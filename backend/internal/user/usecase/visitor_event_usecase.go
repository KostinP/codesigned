package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
	"github.com/kostinp/edu-platform-backend/internal/user/repository"
)

type VisitorEventUsecase interface {
	LogEvent(ctx context.Context, visitorID uuid.UUID, eventType string, eventData map[string]any) error
	GetEventsByVisitor(ctx context.Context, visitorID string) ([]*entity.VisitorEvent, error)
}

type visitorEventUsecase struct {
	repo repository.VisitorEventRepository
}

func NewVisitorEventUsecase(repo repository.VisitorEventRepository) VisitorEventUsecase {
	return &visitorEventUsecase{repo: repo}
}

func (u *visitorEventUsecase) LogEvent(ctx context.Context, visitorID uuid.UUID, eventType string, eventData map[string]any) error {
	event := &entity.VisitorEvent{
		ID:        uuid.New(),
		VisitorID: visitorID,
		EventType: eventType,
		EventData: eventData,
		CreatedAt: time.Now(),
	}
	return u.repo.Create(ctx, event)
}

func (u *visitorEventUsecase) GetEventsByVisitor(ctx context.Context, visitorID string) ([]*entity.VisitorEvent, error) {
	return u.repo.GetEventsByVisitor(ctx, visitorID)
}
