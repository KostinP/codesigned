package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Category, int, error)
}

type PostgresCategoryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryRepository(db *pgxpool.Pool) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	_, err := r.db.Exec(ctx, `
			INSERT INTO categories (id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, category.ID, category.Name, category.Slug, category.Description, category.ParentID, category.AuthorID, category.SortOrder, category.IsVisible, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	_, err := r.db.Exec(ctx, `
			UPDATE categories
			SET name = $1, slug = $2, description = $3, parent_id = $4, sort_order = $5, is_visible = $6, updated_at = $7
			WHERE id = $8
		`, category.Name, category.Slug, category.Description, category.ParentID, category.SortOrder, category.IsVisible, category.UpdatedAt, category.ID)
	return err
}

func (r *PostgresCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
func (r *PostgresCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	row := r.db.QueryRow(ctx, `
			SELECT id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at
			FROM categories WHERE id = $1
		`, id)
	category := &entity.Category{}

	err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.ParentID, &category.AuthorID, &category.SortOrder, &category.IsVisible, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}
func (r *PostgresCategoryRepository) List(ctx context.Context, pag pagination.Params) ([]*entity.Category, int, error) {
	baseQuery := `
		SELECT id, name, slug, description, parent_id, author_id, sort_order, is_visible, created_at, updated_at
		FROM categories`
	query, args := pagination.SQLWithPagination(baseQuery, pag, map[string]string{"created_at": "created_at", "name": "name"})

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*entity.Category
	for rows.Next() {
		category := &entity.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.ParentID, &category.AuthorID, &category.SortOrder, &category.IsVisible, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	countQuery := `SELECT COUNT(*) FROM categories`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
