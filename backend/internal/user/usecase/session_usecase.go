package usecase

import (
	"context"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/user/repository"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

// Интерфейс для SessionUsecase, который будет использоваться в wire.Bind
type SessionUsecase interface {
	CreateSession(ctx context.Context, userID uuid.UUID, token, userAgent, ip, country, city string, expiresIn time.Duration) (*entity.UserSession, error)
	UpdateLastActive(ctx context.Context, sessionID uuid.UUID) error
	ListSessions(ctx context.Context, userID uuid.UUID) ([]*entity.UserSession, error)
	DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error
	SetInactivityTimeout(ctx context.Context, userID uuid.UUID, timeout time.Duration) error
	GetInactivityTimeout(ctx context.Context, userID uuid.UUID) (time.Duration, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entity.UserSession, error)
	DeleteExpiredSessions(ctx context.Context) error
}

// Структура-реализация
type SessionUsecaseImpl struct {
	repo repository.SessionRepository
}

// Конструктор
func NewSessionUsecase(repo repository.SessionRepository) *SessionUsecaseImpl {
	return &SessionUsecaseImpl{repo: repo}
}

// Реализация методов

func (s *SessionUsecaseImpl) CreateSession(
	ctx context.Context,
	userID uuid.UUID,
	token, userAgent, ip, country, city string,
	expiresIn time.Duration,
) (*entity.UserSession, error) {
	now := time.Now()
	session := &entity.UserSession{
		ID:           uuid.New(),
		UserID:       userID,
		Token:        token,
		UserAgent:    userAgent,
		IPAddress:    ip,
		Country:      country,
		City:         city,
		CreatedAt:    now,
		LastActiveAt: now,
		ExpiresAt:    now.Add(expiresIn),
	}
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionUsecaseImpl) UpdateLastActive(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.UpdateLastActive(ctx, sessionID)
}

func (s *SessionUsecaseImpl) ListSessions(ctx context.Context, userID uuid.UUID) ([]*entity.UserSession, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *SessionUsecaseImpl) DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	return s.repo.Delete(ctx, userID, sessionID)
}

func (s *SessionUsecaseImpl) SetInactivityTimeout(ctx context.Context, userID uuid.UUID, timeout time.Duration) error {
	return s.repo.SaveInactivityTimeout(ctx, userID, timeout)
}

func (s *SessionUsecaseImpl) GetInactivityTimeout(ctx context.Context, userID uuid.UUID) (time.Duration, error) {
	return s.repo.GetInactivityTimeout(ctx, userID)
}

func (s *SessionUsecaseImpl) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entity.UserSession, error) {
	return s.repo.FindByID(ctx, sessionID)
}

func (s *SessionUsecaseImpl) DeleteExpiredSessions(ctx context.Context) error {
	return s.repo.DeleteExpiredSessions(ctx)
}
