package entity

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
}

func (b *Base) Init(authorID uuid.UUID) {
	b.ID = uuid.New()
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now
	b.AuthorID = authorID
}

func (b *Base) Touch() {
	b.UpdatedAt = time.Now().UTC()
}
