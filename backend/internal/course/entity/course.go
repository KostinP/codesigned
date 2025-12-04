package entity

import (
	"time"

	"github.com/kostinp/edu-platform-backend/internal/shared/entity"
)

type Course struct {
	entity.Base

	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       int        `json:"price"`
	ImageURL    string     `json:"image_url"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
