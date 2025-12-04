package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUnspecified Role = "unspecified"
	RoleGuest       Role = "guest"
	RoleStudent     Role = "student"
	RoleTeacher     Role = "teacher"
	RoleAdmin       Role = "admin"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	VisitorID       *uuid.UUID `json:"visitor_id,omitempty"`
	TelegramID      *int64     `json:"telegram_id,omitempty"`
	Username        *string    `json:"username,omitempty"`
	FullName        *string    `json:"full_name,omitempty"`
	PhotoURL        *string    `json:"photo_url,omitempty"`
	Email           *string    `json:"email,omitempty"`
	SubscribeToNews bool       `json:"subscribe_to_newsletter"`
	Role            Role       `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
