package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, sub *models.Subscription) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error)
}

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, tx pgx.Tx, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, purchase_token, order_id, status, auto_renew, started_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := tx.Exec(ctx, query,
		sub.UserID, sub.PlanID, sub.PurchaseToken, sub.OrderID,
		sub.Status, sub.AutoRenew, sub.StartedAt, sub.ExpiresAt,
	)
	if err != nil {
		return errors.New("error creating subscription")
	}
	return nil
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, purchase_token, order_id, status, auto_renew, started_at, expires_at, created_at, updated_at
		FROM subscriptions WHERE user_id = $1
		ORDER BY created_at DESC LIMIT 1
	`
	var sub models.Subscription
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.PurchaseToken, &sub.OrderID,
		&sub.Status, &sub.AutoRenew, &sub.StartedAt, &sub.ExpiresAt,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("subscription not found")
	}
	return &sub, nil
}
