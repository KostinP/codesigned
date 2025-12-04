package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/entity"
)

type Lesson struct {
	entity.Base
	ModuleID  uuid.UUID  `json:"module_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Duration  int        `json:"duration"`
	Ordinal   int        `json:"ordinal"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
