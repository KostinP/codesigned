package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

type SessionRepository interface {
	Save(ctx context.Context, session *entity.UserSession) error
	UpdateLastActive(ctx context.Context, sessionID uuid.UUID) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserSession, error)
	Delete(ctx context.Context, userID, sessionID uuid.UUID) error
	SaveInactivityTimeout(ctx context.Context, userID uuid.UUID, timeout time.Duration) error
	GetInactivityTimeout(ctx context.Context, userID uuid.UUID) (time.Duration, error)
	FindByID(ctx context.Context, sessionID uuid.UUID) (*entity.UserSession, error)
	DeleteExpiredSessions(ctx context.Context) error
}

type PostgresSessionRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

// Сохраняем сессию
func (r *PostgresSessionRepository) Save(ctx context.Context, s *entity.UserSession) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_sessions (
			id, user_id, token, user_agent, ip_address, country, city, created_at, last_active_at, expires_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`,
		s.ID, s.UserID, s.Token, s.UserAgent, s.IPAddress, s.Country, s.City,
		s.CreatedAt, s.LastActiveAt, s.ExpiresAt,
	)
	return err
}

// Обновляем время последней активности
func (r *PostgresSessionRepository) UpdateLastActive(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE user_sessions
		SET last_active_at = NOW()
		WHERE id = $1
	`, sessionID)
	return err
}

// Получаем список сессий пользователя
func (r *PostgresSessionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserSession, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, token, user_agent, ip_address, country, city, created_at, last_active_at, expires_at
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entity.UserSession
	for rows.Next() {
		var s entity.UserSession
		err := rows.Scan(
			&s.ID, &s.Token, &s.UserAgent, &s.IPAddress, &s.Country, &s.City,
			&s.CreatedAt, &s.LastActiveAt, &s.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		s.UserID = userID
		sessions = append(sessions, &s)
	}

	return sessions, nil
}

// Удаление конкретной сессии
func (r *PostgresSessionRepository) Delete(ctx context.Context, userID, sessionID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM user_sessions
		WHERE id = $1 AND user_id = $2
	`, sessionID, userID)
	return err
}

// Храним таймаут неактивности (в секундах)
func (r *PostgresSessionRepository) SaveInactivityTimeout(ctx context.Context, userID uuid.UUID, timeout time.Duration) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_inactivity_timeout (user_id, timeout_seconds)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET timeout_seconds = EXCLUDED.timeout_seconds
	`, userID, int64(timeout.Seconds()))
	return err
}

func (r *PostgresSessionRepository) GetInactivityTimeout(ctx context.Context, userID uuid.UUID) (time.Duration, error) {
	var seconds int64
	err := r.db.QueryRow(ctx, `
		SELECT timeout_seconds FROM user_inactivity_timeout WHERE user_id = $1
	`, userID).Scan(&seconds)
	if err != nil {
		if err == pgx.ErrNoRows {
			// По умолчанию 6 месяцев
			return 4380 * time.Hour, nil
		}
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

func (r *PostgresSessionRepository) FindByID(ctx context.Context, sessionID uuid.UUID) (*entity.UserSession, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, token, user_agent, ip_address, country, city, created_at, last_active_at, expires_at
		FROM user_sessions
		WHERE id = $1
	`, sessionID)

	var s entity.UserSession
	err := row.Scan(
		&s.ID, &s.UserID, &s.Token, &s.UserAgent, &s.IPAddress, &s.Country, &s.City,
		&s.CreatedAt, &s.LastActiveAt, &s.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}
