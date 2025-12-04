package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/lesson/entity"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *entity.Lesson) error
	Update(ctx context.Context, lesson *entity.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Lesson, int, error)
}

type PostgresLessonRepository struct {
	db *pgxpool.Pool
}

func NewPostgresLessonRepository(db *pgxpool.Pool) *PostgresLessonRepository {
	return &PostgresLessonRepository{db: db}
}

func (r *PostgresLessonRepository) Create(ctx context.Context, lesson *entity.Lesson) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO lessons (id, module_id, title, content, duration, ordinal, author_id, created_at, updated_at) 		
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, lesson.ID, lesson.ModuleID, lesson.Title, lesson.Content, lesson.Duration, lesson.Ordinal, lesson.AuthorID, lesson.CreatedAt, lesson.UpdatedAt)
	return err
}

func (r *PostgresLessonRepository) Update(ctx context.Context, lesson *entity.Lesson) error {
	_, err := r.db.Exec(ctx, `
			UPDATE lessons
			SET module_id = $1, title = $2, content = $3, duration = $4, ordinal = $5, updated_at = $6, deleted_at = $7 		
			WHERE id = $8
		`, lesson.ModuleID, lesson.Title, lesson.Content, lesson.Duration, lesson.Ordinal, lesson.UpdatedAt, lesson.DeletedAt, lesson.ID)
	return err
}

func (r *PostgresLessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE lessons SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *PostgresLessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, module_id, title, content, duration, ordinal, author_id, created_at, updated_at, deleted_at 		
		FROM lessons WHERE id = $1 AND deleted_at IS NULL
		`, id)
	lesson := &entity.Lesson{}

	err := row.Scan(&lesson.ID, &lesson.ModuleID, &lesson.Title, &lesson.Content, &lesson.Duration, &lesson.Ordinal, &lesson.AuthorID, &lesson.CreatedAt, &lesson.UpdatedAt, &lesson.DeletedAt)
	if err != nil {
		return nil, err
	}

	return lesson, nil
}
func (r *PostgresLessonRepository) List(ctx context.Context, pag pagination.Params) ([]*entity.Lesson, int, error) {
	baseQuery := `
		SELECT id, module_id, title, content, duration, ordinal, author_id, created_at, updated_at, deleted_at 		
		FROM lessons WHERE deleted_at IS NULL`
	query, args := pagination.SQLWithPagination(baseQuery, pag, map[string]string{"created_at": "created_at", "title": "title"})

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var lessons []*entity.Lesson
	for rows.Next() {
		lesson := &entity.Lesson{}
		err := rows.Scan(&lesson.ID, &lesson.ModuleID, &lesson.Title, &lesson.Content, &lesson.Duration, &lesson.Ordinal, &lesson.AuthorID, &lesson.CreatedAt, &lesson.UpdatedAt, &lesson.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		lessons = append(lessons, lesson)
	}

	countQuery := `SELECT COUNT(*) FROM lessons WHERE deleted_at IS NULL`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}
