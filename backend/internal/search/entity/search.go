// internal/search/entity/search.go
package entity

import (
	"time"

	"github.com/google/uuid"
	categoryEntity "github.com/kostinp/edu-platform-backend/internal/category/entity"
	tagEntity "github.com/kostinp/edu-platform-backend/internal/tag/entity"
)

type SearchEntity struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"` // "course", "lesson", "module"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	// Связи
	Categories []*categoryEntity.Category `json:"categories,omitempty"`
	Tags       []*tagEntity.Tag           `json:"tags,omitempty"`
	// Поля специфичные для типа
	Price    *int       `json:"price,omitempty"`
	Duration *int       `json:"duration,omitempty"`
	AuthorID *uuid.UUID `json:"author_id,omitempty"`
	// Релевантность
	Relevance float64 `json:"relevance"`
}

type SearchFilters struct {
	Query       string     `json:"query,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	ContentType []string   `json:"content_type,omitempty"`
	PriceMin    *int       `json:"price_min,omitempty"`
	PriceMax    *int       `json:"price_max,omitempty"`
	Level       []string   `json:"level,omitempty"`
	AuthorID    *uuid.UUID `json:"author_id,omitempty"`
	// Пагинация
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

type SearchResult struct {
	Entities []*SearchEntity `json:"entities"`
	Total    int             `json:"total"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
}

// Breadcrumb для навигации
type Breadcrumb struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
	Type string    `json:"type,omitempty"` // "category", "tag", "search"
}
