package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/module/entity"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type ModuleRepository interface {
	Create(ctx context.Context, module *entity.Module) error
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Module, int, error)
}

type PostgresModuleRepository struct {
	db *pgxpool.Pool
}

func NewPostgresModuleRepository(db *pgxpool.Pool) *PostgresModuleRepository {
	return &PostgresModuleRepository{db: db}
}

func (r *PostgresModuleRepository) Create(ctx context.Context, module *entity.Module) error {
	_, err := r.db.Exec(ctx, `
			INSERT INTO modules (id, course_id, title, description, ordinal, author_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, module.ID, module.CourseID, module.Title, module.Description, module.Ordinal, module.AuthorID, module.CreatedAt, module.UpdatedAt)
	return err
}

func (r *PostgresModuleRepository) Update(ctx context.Context, module *entity.Module) error {
	_, err := r.db.Exec(ctx, `
			UPDATE modules
			SET course_id = $1, title = $2, description = $3, ordinal = $4, updated_at = $5, deleted_at = $6
			WHERE id = $7
		`, module.CourseID, module.Title, module.Description, module.Ordinal, module.UpdatedAt, module.DeletedAt, module.ID)
	return err
}

func (r *PostgresModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE modules SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *PostgresModuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error) {
	row := r.db.QueryRow(ctx, `
			SELECT id, course_id, title, description, ordinal, author_id, created_at, updated_at, deleted_at
			FROM modules WHERE id = $1 AND deleted_at IS NULL
		`, id)
	module := &entity.Module{}
	err := row.Scan(&module.ID, &module.CourseID, &module.Title, &module.Description, &module.Ordinal, &module.AuthorID, &module.CreatedAt, &module.UpdatedAt, &module.DeletedAt)
	if err != nil {
		return nil, err
	}
	return module, nil
}

func (r *PostgresModuleRepository) List(ctx context.Context, pag pagination.Params) ([]*entity.Module, int, error) {
	baseQuery := `
		SELECT id, course_id, title, description, ordinal, author_id, created_at, updated_at, deleted_at
		FROM modules WHERE deleted_at IS NULL`
	query, args := pagination.SQLWithPagination(baseQuery, pag, map[string]string{"created_at": "created_at", "title": "title"})

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var modules []*entity.Module

	for rows.Next() {
		module := &entity.Module{}

		err := rows.Scan(&module.ID, &module.CourseID, &module.Title, &module.Description, &module.Ordinal, &module.AuthorID, &module.CreatedAt, &module.UpdatedAt, &module.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		modules = append(modules, module)
	}

	countQuery := `SELECT COUNT(*) FROM modules WHERE deleted_at IS NULL`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return modules, total, nil
}
