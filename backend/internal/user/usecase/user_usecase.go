package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/telegram"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

// UserRepository описывает методы репозитория, которые нужны сервису
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

// UserService реализует бизнес-логику для пользователей
type UserService struct {
	sessionUC SessionUsecase
	repo      UserRepository
}

// NewUserService создаёт UserService с переданным репозиторием
func NewUserService(repo UserRepository, sessionUC SessionUsecase) *UserService {
	return &UserService{
		repo:      repo,
		sessionUC: sessionUC,
	}
}

// CreateGuestUser создаёт нового гостевого пользователя
func (s *UserService) CreateGuestUser(ctx context.Context) (*entity.User, error) {
	user := &entity.User{
		Role:      entity.RoleGuest,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpgradeToTelegramUser обновляет гостя до пользователя с TelegramID и ролью student по умолчанию
func (s *UserService) UpgradeToTelegramUser(ctx context.Context, userID uuid.UUID, telegramID int64, username, fullName string) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	user.TelegramID = &telegramID
	user.Username = &username
	user.FullName = &fullName
	user.Role = entity.RoleStudent
	user.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// LinkVisitorToUser связывает посетителя (visitor) с пользователем (user)
func (s *UserService) LinkVisitorToUser(ctx context.Context, userID, visitorID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("не удалось найти пользователя: %w", err)
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	user.VisitorID = &visitorID
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("не удалось обновить пользователя: %w", err)
	}
	return nil
}

// GetByTelegramID возвращает пользователя по Telegram ID
func (s *UserService) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	return s.repo.GetByTelegramID(ctx, telegramID)
}

// CreateFromTelegramAuth создает нового пользователя из Telegram-авторизации
func (s *UserService) CreateFromTelegramAuth(ctx context.Context, data telegram.AuthData) (*entity.User, error) {
	user := &entity.User{
		TelegramID: &data.ID,
		Username:   &data.Username,
		FullName:   &data.FirstName,
		Role:       entity.RoleStudent,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := s.repo.Create(ctx, user)
	return user, err
}

func (s *UserService) CreateUserSession(
	ctx context.Context,
	userID uuid.UUID,
	ip, userAgent, country, city string,
) (*entity.UserSession, error) {
	// Срок жизни — 6 месяцев
	expiresIn := time.Hour * 24 * 30 * 6

	token := uuid.New().String() // можно использовать как внутренний ID
	return s.sessionUC.CreateSession(ctx, userID, token, userAgent, ip, country, city, expiresIn)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}
