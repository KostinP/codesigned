package entity

import "github.com/kostinp/edu-platform-backend/internal/shared/entity"

type Tag struct {
	entity.Base

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}
