package models

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	MaxApps   int       `json:"max_apps"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
