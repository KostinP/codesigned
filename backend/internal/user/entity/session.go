package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserSession модель сессии пользователя
type UserSession struct {
	ID           uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID       uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Token        string    `json:"token" example:"jwt.token.here"`
	UserAgent    string    `json:"user_agent" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"`
	IPAddress    string    `json:"ip_address" example:"192.168.1.1"`
	Country      string    `json:"country" example:"Russia"`
	City         string    `json:"city" example:"Moscow"`
	CreatedAt    time.Time `json:"created_at" example:"2025-07-24T18:25:43.511Z"`
	LastActiveAt time.Time `json:"last_active_at" example:"2025-07-25T10:15:00.000Z"`
	ExpiresAt    time.Time `json:"expires_at" example:"2025-08-24T18:25:43.511Z"`
}
