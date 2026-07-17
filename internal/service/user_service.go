package service

import (
	"context"
	"errors"
	"time"

	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/harmelson/tocouaboa-portfolio/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	Create(ctx context.Context, user *models.UserDTO) error
}

type userService struct {
	pool     *pgxpool.Pool
	userRepo repository.UserRepository
	planRepo repository.PlanRepository
	subRepo  repository.SubscriptionRepository
}

func NewUserService(
	pool *pgxpool.Pool,
	userRepo repository.UserRepository,
	planRepo repository.PlanRepository,
	subRepo repository.SubscriptionRepository,
) UserService {
	return &userService{
		pool:     pool,
		userRepo: userRepo,
		planRepo: planRepo,
		subRepo:  subRepo,
	}
}

func (s *userService) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	return s.userRepo.GetByGoogleID(ctx, googleID)
}

func (s *userService) Create(ctx context.Context, user *models.UserDTO) error {
	existingUser, err := s.userRepo.GetByGoogleID(ctx, user.GoogleID)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return errors.New("failed to check existing user")
	}

	freePlan, err := s.planRepo.GetByName(ctx, "free")
	if err != nil {
		return errors.New("free plan not found, contact support")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	userID, err := s.userRepo.Create(ctx, tx, user)
	if err != nil {
		return err
	}

	now := time.Now()
	const freeSubscriptionYears = 100
	sub := &models.Subscription{
		UserID:        userID,
		PlanID:        freePlan.ID,
		PurchaseToken: "free_" + userID.String(),
		OrderID:       "",
		Status:        string(models.StatusFree),
		AutoRenew:     false,
		StartedAt:     now,
		ExpiresAt:     now.AddDate(freeSubscriptionYears, 0, 0),
	}

	if err := s.subRepo.Create(ctx, tx, sub); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
