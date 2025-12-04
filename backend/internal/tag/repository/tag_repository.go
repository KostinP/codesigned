// internal/tag/repository/tag_repository.go
package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/kostinp/edu-platform-backend/internal/tag/entity"
)

var ErrTagNotFound = errors.New("tag not found")

// TagRepository описывает интерфейс работы с тегами
type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	Update(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Tag, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresTagRepository — реализация TagRepository для PostgreSQL
type PostgresTagRepository struct {
	db *pgxpool.Pool
}

// allowedSortFields: API-поля => БД-поля
var allowedSortFields = map[string]string{
	"created_at": "created_at",
	"name":       "name",
}

// NewPostgresTagRepository создает новый экземпляр PostgresTagRepository
func NewPostgresTagRepository(db *pgxpool.Pool) *PostgresTagRepository {
	return &PostgresTagRepository{db: db}
}

// Create добавляет новый тег в базу
func (r *PostgresTagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tags (id, name, description, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tag.ID, tag.Name, tag.Description, tag.AuthorID, tag.CreatedAt, tag.UpdatedAt)
	return err
}

// Update обновляет существующий тег
func (r *PostgresTagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	existingTag, err := r.GetByID(ctx, tag.ID)
	if err != nil {
		return err
	}
	if existingTag == nil {
		return ErrTagNotFound
	}

	_, err = r.db.Exec(ctx, `
		UPDATE tags
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`, tag.Name, tag.Description, tag.UpdatedAt, tag.ID)
	return err
}

// Delete удаляет тег по ID
func (r *PostgresTagRepository) Delete(ctx context.Context, tagID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM tags WHERE id = $1
	`, tagID)
	return err
}

// GetByID возвращает тег по ID
func (r *PostgresTagRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, description, author_id, created_at, updated_at
		FROM tags WHERE id = $1
	`, id)

	tag := &entity.Tag{}
	err := row.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.AuthorID, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// List возвращает список тегов с пагинацией
func (r *PostgresTagRepository) List(ctx context.Context, pag pagination.Params) ([]*entity.Tag, int, error) {
	baseQuery := `
		SELECT id, name, description, author_id, created_at, updated_at
		FROM tags
	`
	query, args := pagination.SQLWithPagination(baseQuery, pag, allowedSortFields)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var tags []*entity.Tag
	for rows.Next() {
		tag := &entity.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.AuthorID, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, tag)
	}
	countQuery := `SELECT COUNT(*) FROM tags`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return tags, total, nil
}
