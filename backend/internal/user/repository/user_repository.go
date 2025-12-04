package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			id, visitor_id, telegram_id, first_name, last_name, username, photo_url,
			created_at, updated_at, deleted_at, email, subscribe_to_newsletter, role
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	var firstName, lastName string
	if user.FullName != nil {
		names := splitFullName(*user.FullName)
		firstName = names[0]
		if len(names) > 1 {
			lastName = names[1]
		}
	}

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.VisitorID,
		user.TelegramID,
		firstName,
		lastName,
		user.Username,
		user.PhotoURL,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
		user.Email,
		user.SubscribeToNews,
		string(user.Role),
	)
	return err
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users SET
			visitor_id = $1,
			telegram_id = $2,
			first_name = $3,
			last_name = $4,
			username = $5,
			photo_url = $6,
			updated_at = $7,
			deleted_at = $8,
			email = $9,
			subscribe_to_newsletter = $10,
			role = $11
		WHERE id = $12
	`

	var firstName, lastName string
	if user.FullName != nil {
		names := splitFullName(*user.FullName)
		firstName = names[0]
		if len(names) > 1 {
			lastName = names[1]
		}
	}

	cmdTag, err := r.pool.Exec(ctx, query,
		user.VisitorID,
		user.TelegramID,
		firstName,
		lastName,
		user.Username,
		user.PhotoURL,
		user.UpdatedAt,
		user.DeletedAt,
		user.Email,
		user.SubscribeToNews,
		string(user.Role),
		user.ID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("пользователь не найден для обновления")
	}

	return nil
}

func (r *PostgresUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	query := `
		SELECT id, visitor_id, telegram_id, first_name, last_name, username, photo_url,
		       created_at, updated_at, deleted_at, email, subscribe_to_newsletter, role
		FROM users WHERE telegram_id = $1 AND deleted_at IS NULL
	`

	user := &entity.User{}
	var firstName, lastName string

	err := r.pool.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.VisitorID,
		&user.TelegramID,
		&firstName,
		&lastName,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Email,
		&user.SubscribeToNews,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	fullName := combineFullName(firstName, lastName)
	user.FullName = &fullName

	return user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, visitor_id, telegram_id, first_name, last_name, username, photo_url,
		       created_at, updated_at, deleted_at, email, subscribe_to_newsletter, role
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`

	user := &entity.User{}
	var firstName, lastName string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.VisitorID,
		&user.TelegramID,
		&firstName,
		&lastName,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Email,
		&user.SubscribeToNews,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	fullName := combineFullName(firstName, lastName)
	user.FullName = &fullName

	return user, nil
}

func splitFullName(fullName string) []string {
	return strings.Fields(fullName) // лучше split по пробелам с trim
}

func combineFullName(firstName, lastName string) string {
	if lastName == "" {
		return firstName
	}
	return firstName + " " + lastName
}
