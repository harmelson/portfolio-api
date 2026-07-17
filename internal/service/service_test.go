package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/harmelson/tocouaboa-portfolio/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type userRepositoryStub struct {
	getByGoogleID func(context.Context, string) (*models.User, error)
}

func (s userRepositoryStub) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	return s.getByGoogleID(ctx, googleID)
}

func (s userRepositoryStub) Create(context.Context, pgx.Tx, *models.UserDTO) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

type planRepositoryStub struct {
	getByName func(context.Context, string) (*models.Plan, error)
	getByID   func(context.Context, uuid.UUID) (*models.Plan, error)
}

func (s planRepositoryStub) GetByName(ctx context.Context, name string) (*models.Plan, error) {
	return s.getByName(ctx, name)
}

func (s planRepositoryStub) GetByID(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	return s.getByID(ctx, id)
}

type subscriptionRepositoryStub struct {
	getByUserID func(context.Context, uuid.UUID) (*models.Subscription, error)
}

func (s subscriptionRepositoryStub) Create(context.Context, pgx.Tx, *models.Subscription) error {
	return errors.New("not implemented")
}

func (s subscriptionRepositoryStub) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	return s.getByUserID(ctx, userID)
}

type ServiceSuite struct {
	suite.Suite
	ctx context.Context
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ServiceSuite) TestUserService() {
	s.Run("gets user by Google ID", func() {
		expected := &models.User{GoogleID: "google-id"}
		service := NewUserService(nil, userRepositoryStub{
			getByGoogleID: func(_ context.Context, googleID string) (*models.User, error) {
				s.Equal(expected.GoogleID, googleID)
				return expected, nil
			},
		}, nil, nil)

		user, err := service.GetByGoogleID(s.ctx, expected.GoogleID)
		s.Require().NoError(err)
		s.Same(expected, user)
	})

	s.Run("rejects an existing user", func() {
		service := NewUserService(nil, userRepositoryStub{
			getByGoogleID: func(context.Context, string) (*models.User, error) {
				return &models.User{}, nil
			},
		}, nil, nil)

		s.EqualError(service.Create(s.ctx, &models.UserDTO{GoogleID: "existing"}), "user already exists")
	})

	s.Run("returns lookup failure", func() {
		service := NewUserService(nil, userRepositoryStub{
			getByGoogleID: func(context.Context, string) (*models.User, error) {
				return nil, errors.New("database unavailable")
			},
		}, nil, nil)

		s.EqualError(service.Create(s.ctx, &models.UserDTO{GoogleID: "new"}), "failed to check existing user")
	})

	s.Run("returns missing free plan", func() {
		service := NewUserService(nil, userRepositoryStub{
			getByGoogleID: func(context.Context, string) (*models.User, error) {
				return nil, pgx.ErrNoRows
			},
		}, planRepositoryStub{
			getByName: func(context.Context, string) (*models.Plan, error) {
				return nil, errors.New("not found")
			},
			getByID: func(context.Context, uuid.UUID) (*models.Plan, error) {
				return nil, errors.New("not implemented")
			},
		}, nil)

		s.EqualError(service.Create(s.ctx, &models.UserDTO{GoogleID: "new"}), "free plan not found, contact support")
	})
}

func (s *ServiceSuite) TestSubscriptionService() {
	s.Run("returns subscription with plan", func() {
		userID := uuid.New()
		planID := uuid.New()
		service := NewSubscriptionService(subscriptionRepositoryStub{
			getByUserID: func(_ context.Context, id uuid.UUID) (*models.Subscription, error) {
				s.Equal(userID, id)
				return &models.Subscription{UserID: userID, PlanID: planID, Status: string(models.StatusFree)}, nil
			},
		}, planRepositoryStub{
			getByName: func(context.Context, string) (*models.Plan, error) {
				return nil, errors.New("not implemented")
			},
			getByID: func(_ context.Context, id uuid.UUID) (*models.Plan, error) {
				s.Equal(planID, id)
				return &models.Plan{ID: planID, Name: "free", MaxApps: 1, Price: 0}, nil
			},
		})

		subscription, err := service.GetByUserIDWithPlan(s.ctx, userID)
		s.Require().NoError(err)
		s.Equal("free", subscription.PlanName)
		s.Equal(1, subscription.MaxApps)
		s.Equal(planID, subscription.PlanID)
	})

	s.Run("returns error when plan is missing", func() {
		service := NewSubscriptionService(subscriptionRepositoryStub{
			getByUserID: func(context.Context, uuid.UUID) (*models.Subscription, error) {
				return &models.Subscription{PlanID: uuid.New()}, nil
			},
		}, planRepositoryStub{
			getByName: func(context.Context, string) (*models.Plan, error) {
				return nil, errors.New("not implemented")
			},
			getByID: func(context.Context, uuid.UUID) (*models.Plan, error) {
				return nil, errors.New("not found")
			},
		})

		_, err := service.GetByUserIDWithPlan(s.ctx, uuid.New())
		s.EqualError(err, "plan not found for subscription")
	})
}

var _ repository.UserRepository = userRepositoryStub{}
var _ repository.PlanRepository = planRepositoryStub{}
var _ repository.SubscriptionRepository = subscriptionRepositoryStub{}
