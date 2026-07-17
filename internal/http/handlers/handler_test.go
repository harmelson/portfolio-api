package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/harmelson/tocouaboa-portfolio/internal/contextutil"
	"github.com/harmelson/tocouaboa-portfolio/internal/models"
	"github.com/stretchr/testify/suite"
)

type subscriptionServiceStub struct {
	getByUserIDWithPlan func(context.Context, uuid.UUID) (*models.SubscriptionWithPlan, error)
}

func (s subscriptionServiceStub) GetByUserIDWithPlan(ctx context.Context, userID uuid.UUID) (*models.SubscriptionWithPlan, error) {
	return s.getByUserIDWithPlan(ctx, userID)
}

type userServiceStub struct{}

func (userServiceStub) GetByGoogleID(context.Context, string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (userServiceStub) Create(context.Context, *models.UserDTO) error {
	return errors.New("not implemented")
}

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) TestSubscriptionHandlerGetMy() {
	s.Run("returns unauthorized without an authenticated user", func() {
		handler := NewSubscriptionHandler(subscriptionServiceStub{
			getByUserIDWithPlan: func(context.Context, uuid.UUID) (*models.SubscriptionWithPlan, error) {
				s.Fail("service must not be called")
				return nil, nil
			},
		})
		request := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/me", nil)
		response := httptest.NewRecorder()

		handler.GetMy(response, request)

		s.Equal(http.StatusUnauthorized, response.Code)
		s.Contains(response.Body.String(), "user not authenticated")
	})

	s.Run("returns subscription for authenticated user", func() {
		userID := uuid.New()
		handler := NewSubscriptionHandler(subscriptionServiceStub{
			getByUserIDWithPlan: func(_ context.Context, id uuid.UUID) (*models.SubscriptionWithPlan, error) {
				s.Equal(userID, id)
				return &models.SubscriptionWithPlan{PlanName: "free", MaxApps: 1}, nil
			},
		})
		request := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/me", nil)
		request = request.WithContext(contextutil.WithUserID(request.Context(), userID))
		response := httptest.NewRecorder()

		handler.GetMy(response, request)

		s.Equal(http.StatusOK, response.Code)
		s.Contains(response.Body.String(), `"plan_name":"free"`)
	})

	s.Run("returns not found when subscription is missing", func() {
		handler := NewSubscriptionHandler(subscriptionServiceStub{
			getByUserIDWithPlan: func(context.Context, uuid.UUID) (*models.SubscriptionWithPlan, error) {
				return nil, errors.New("subscription not found")
			},
		})
		request := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/me", nil)
		request = request.WithContext(contextutil.WithUserID(request.Context(), uuid.New()))
		response := httptest.NewRecorder()

		handler.GetMy(response, request)

		s.Equal(http.StatusNotFound, response.Code)
		s.Contains(response.Body.String(), "subscription not found")
	})
}

func (s *HandlerSuite) TestUserHandler() {
	s.Run("get user rejects missing token", func() {
		handler := NewUserHandler(userServiceStub{})
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/get", nil)
		response := httptest.NewRecorder()

		handler.GetByGoogleID(response, request)

		s.Equal(http.StatusUnauthorized, response.Code)
		s.Contains(response.Body.String(), "auth token not provided")
	})

	s.Run("create user rejects missing token", func() {
		handler := NewUserHandler(userServiceStub{})
		request := httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
		response := httptest.NewRecorder()

		handler.Create(response, request)

		s.Equal(http.StatusUnauthorized, response.Code)
		s.Contains(response.Body.String(), "invalid authentication token")
	})
}
