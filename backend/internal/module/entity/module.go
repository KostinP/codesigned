package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/entity"
)

type Module struct {
	entity.Base

	CourseID    uuid.UUID  `json:"course_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Ordinal     int        `json:"ordinal"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
