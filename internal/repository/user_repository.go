package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	Create(ctx context.Context, tx pgx.Tx, user *models.UserDTO) (uuid.UUID, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `
		SELECT id, email, name, google_id, picture_id, is_active, created_at, updated_at
		FROM users
		WHERE google_id = $1 AND is_active = true
	`
	var user models.User
	err := r.pool.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.GoogleID,
		&user.PictureID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, tx pgx.Tx, user *models.UserDTO) (uuid.UUID, error) {
	query := `
		INSERT INTO users (email, name, google_id, picture_id, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var userID uuid.UUID
	err := tx.QueryRow(ctx, query, user.Email, user.Name, user.GoogleID, user.PictureID, user.IsActive).Scan(&userID)

	if err != nil {
		return uuid.Nil, errors.New("Error creating user")
	}

	return userID, nil
}
