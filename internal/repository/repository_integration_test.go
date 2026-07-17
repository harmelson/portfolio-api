package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	return pool
}

type RepositorySuite struct {
	suite.Suite
	pool *pgxpool.Pool
	ctx  context.Context
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) SetupTest() {
	s.pool = testPool(s.T())
	s.ctx = context.Background()
}

func (s *RepositorySuite) TestReadSeededData() {
	var userID uuid.UUID
	var planID uuid.UUID

	s.Run("gets seeded user", func() {
		user, err := NewUserRepository(s.pool).GetByGoogleID(s.ctx, "dev-test-user")
		s.Require().NoError(err)
		s.Equal("dev@example.com", user.Email)
		s.True(user.IsActive)
		userID = user.ID
	})

	s.Run("gets free plan by name and ID", func() {
		plan, err := NewPlanRepository(s.pool).GetByName(s.ctx, "free")
		s.Require().NoError(err)
		s.Equal(1, plan.MaxApps)
		s.Equal(float64(0), plan.Price)
		planID = plan.ID

		planByID, err := NewPlanRepository(s.pool).GetByID(s.ctx, planID)
		s.Require().NoError(err)
		s.Equal("free", planByID.Name)
	})

	s.Run("gets seeded subscription", func() {
		s.Require().NotEqual(uuid.Nil, userID)
		s.Require().NotEqual(uuid.Nil, planID)

		subscription, err := NewSubscriptionRepository(s.pool).GetByUserID(s.ctx, userID)
		s.Require().NoError(err)
		s.Equal(planID, subscription.PlanID)
		s.Equal(string(models.StatusFree), subscription.Status)
	})
}

func (s *RepositorySuite) TestCreateUserAndSubscriptionInTransaction() {
	tx, err := s.pool.Begin(s.ctx)
	s.Require().NoError(err)
	defer tx.Rollback(s.ctx)

	plan, err := NewPlanRepository(s.pool).GetByName(s.ctx, "free")
	s.Require().NoError(err)

	googleID := "test-" + uuid.NewString()
	userID, err := NewUserRepository(s.pool).Create(s.ctx, tx, &models.UserDTO{
		Email:     googleID + "@example.com",
		Name:      "Repository Test User",
		GoogleID:  googleID,
		PictureID: "https://example.com/test.png",
		IsActive:  true,
	})
	s.Require().NoError(err)

	s.Run("creates subscription for the new user", func() {
		subscription := &models.Subscription{
			UserID:        userID,
			PlanID:        plan.ID,
			PurchaseToken: "free_" + userID.String(),
			Status:        string(models.StatusFree),
			AutoRenew:     false,
			StartedAt:     time.Now(),
			ExpiresAt:     time.Now().AddDate(100, 0, 0),
		}
		s.Require().NoError(NewSubscriptionRepository(s.pool).Create(s.ctx, tx, subscription))

		var subscriptionCount int
		s.Require().NoError(tx.QueryRow(s.ctx, "SELECT COUNT(*) FROM subscriptions WHERE user_id = $1", userID).Scan(&subscriptionCount))
		s.Equal(1, subscriptionCount)
	})
}
