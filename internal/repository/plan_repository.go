package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlanRepository interface {
	GetByName(ctx context.Context, planName string) (*models.Plan, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error)
}

type planRepository struct {
	pool *pgxpool.Pool
}

func NewPlanRepository(pool *pgxpool.Pool) PlanRepository {
	return &planRepository{pool: pool}
}

func (r *planRepository) GetByName(ctx context.Context, planName string) (*models.Plan, error) {
	query := `
		SELECT id, name, max_apps, price, created_at, updated_at FROM plans WHERE name = $1
	`
	var plan models.Plan
	err := r.pool.QueryRow(ctx, query, planName).Scan(
		&plan.ID,
		&plan.Name,
		&plan.MaxApps,
		&plan.Price,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &plan, nil
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	query := `
		SELECT id, name, max_apps, price, created_at, updated_at FROM plans WHERE id = $1
	`

	var plan models.Plan
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.MaxApps,
		&plan.Price,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &plan, nil
}
