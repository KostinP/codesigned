package entity

import (
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/entity"
)

type CategoryAssignment struct {
	entity.Base
	CategoryID uuid.UUID `json:"category_id"`
	TargetType string    `json:"target_type"`
	TargetID   uuid.UUID `json:"target_id"`
}
