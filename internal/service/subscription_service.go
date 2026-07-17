package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/harmelson/tocouaboa-portfolio/internal/repository"
)

type SubscriptionService interface {
	GetByUserIDWithPlan(ctx context.Context, userID uuid.UUID) (*models.SubscriptionWithPlan, error)
}

type subscriptionService struct {
	subRepo  repository.SubscriptionRepository
	planRepo repository.PlanRepository
}

func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
) SubscriptionService {
	return &subscriptionService{subRepo: subRepo, planRepo: planRepo}
}

func (s *subscriptionService) GetByUserIDWithPlan(ctx context.Context, userID uuid.UUID) (*models.SubscriptionWithPlan, error) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.GetByID(ctx, sub.PlanID)
	if err != nil {
		return nil, errors.New("plan not found for subscription")
	}

	return &models.SubscriptionWithPlan{
		Subscription: *sub,
		PlanName:     plan.Name,
		MaxApps:      plan.MaxApps,
		Price:        plan.Price,
	}, nil
}
