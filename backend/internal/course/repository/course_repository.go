// internal/course/repository/course_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/course/entity"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Course, int, error)
}

type PostgresCourseRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCourseRepository(db *pgxpool.Pool) *PostgresCourseRepository {
	return &PostgresCourseRepository{db: db}
}

func (r *PostgresCourseRepository) Create(ctx context.Context, course *entity.Course) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO courses (id, slug, title, description, price, image_url, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, course.ID, course.Slug, course.Title, course.Description, course.Price, course.ImageURL, course.AuthorID, course.CreatedAt, course.UpdatedAt)
	return err
}

func (r *PostgresCourseRepository) Update(ctx context.Context, course *entity.Course) error {
	_, err := r.db.Exec(ctx, `
		UPDATE courses
		SET slug = $1, title = $2, description = $3, price = $4, image_url = $5, updated_at = $6, deleted_at = $7
		WHERE id = $8
	`, course.Slug, course.Title, course.Description, course.Price, course.ImageURL, course.UpdatedAt, course.DeletedAt, course.ID)
	return err
}

func (r *PostgresCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE courses SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *PostgresCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, slug, title, description, price, image_url, author_id, created_at, updated_at, deleted_at
		FROM courses WHERE id = $1 AND deleted_at IS NULL
	`, id)
	course := &entity.Course{}
	err := row.Scan(&course.ID, &course.Slug, &course.Title, &course.Description, &course.Price, &course.ImageURL, &course.AuthorID, &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (r *PostgresCourseRepository) List(ctx context.Context, pag pagination.Params) ([]*entity.Course, int, error) {
	baseQuery := `
		SELECT id, slug, title, description, price, image_url, author_id, created_at, updated_at, deleted_at
		FROM courses WHERE deleted_at IS NULL
	`
	query, args := pagination.SQLWithPagination(baseQuery, pag, map[string]string{"created_at": "created_at", "title": "title"})
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var courses []*entity.Course
	for rows.Next() {
		course := &entity.Course{}
		err := rows.Scan(&course.ID, &course.Slug, &course.Title, &course.Description, &course.Price, &course.ImageURL, &course.AuthorID, &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		courses = append(courses, course)
	}
	countQuery := `SELECT COUNT(*) FROM courses WHERE deleted_at IS NULL`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return courses, total, nil
}
