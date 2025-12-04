package entity

import (
	"time"

	"github.com/google/uuid"
)

type VisitorEvent struct {
	ID        uuid.UUID      `json:"id"`
	VisitorID uuid.UUID      `json:"visitor_id"`
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}
