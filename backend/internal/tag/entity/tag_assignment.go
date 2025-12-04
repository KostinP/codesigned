package entity

import (
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/entity"
)

type TagAssignment struct {
	entity.Base

	TagID      uuid.UUID `json:"tag_id" db:"tag_id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id" db:"entity_id"`
}
