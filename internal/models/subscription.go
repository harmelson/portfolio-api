package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	PlanID        uuid.UUID `json:"plan_id"`
	PurchaseToken string    `json:"purchase_token"`
	OrderID       string    `json:"order_id"`
	Status        string    `json:"status"`
	AutoRenew     bool      `json:"auto_renew"`
	StartedAt     time.Time `json:"started_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SubscriptionWithPlan struct {
	Subscription
	PlanName string  `json:"plan_name"`
	MaxApps  int     `json:"max_apps"`
	Price    float64 `json:"price"`
}

type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusCanceled SubscriptionStatus = "canceled"
	StatusExpired  SubscriptionStatus = "expired"
	StatusPending  SubscriptionStatus = "pending"
	StatusFree     SubscriptionStatus = "free"
)
