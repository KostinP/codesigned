// internal/category/entity/category.go
package entity

import (
	"time"

	"github.com/google/uuid"
	tagEntity "github.com/kostinp/edu-platform-backend/internal/tag/entity"
)

type Category struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	AuthorID    uuid.UUID  `json:"author_id"`
	SortOrder   int        `json:"sort_order"`
	IsVisible   bool       `json:"is_visible"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	// Навигационные поля
	Parent   *Category   `json:"parent,omitempty"`
	Children []*Category `json:"children,omitempty"`
}

// TreeCategory - категория для дерева
type TreeCategory struct {
	*Category
	Children []*TreeCategory  `json:"children"`
	Tags     []*tagEntity.Tag `json:"tags,omitempty"`
}

// CategoryWithStats - категория со статистикой
type CategoryWithStats struct {
	Category
	CoursesCount int `json:"courses_count"`
	LessonsCount int `json:"lessons_count"`
	ModulesCount int `json:"modules_count"`
}
