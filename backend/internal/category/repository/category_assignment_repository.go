package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
)

type CategoryAssignmentRepository interface {
	Assign(ctx context.Context, assignment *entity.CategoryAssignment) error
	Unassign(ctx context.Context, id uuid.UUID) error
	ListByEntity(ctx context.Context, targetID uuid.UUID, targetType string) ([]*entity.CategoryAssignment, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryAssignment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CategoryAssignment, error)
}

type PostgresCategoryAssignmentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryAssignmentRepository(db *pgxpool.Pool) *PostgresCategoryAssignmentRepository {
	return &PostgresCategoryAssignmentRepository{db: db}
}

func (r *PostgresCategoryAssignmentRepository) Assign(ctx context.Context, ca *entity.CategoryAssignment) error {
	_, err := r.db.Exec(ctx, `
			INSERT INTO category_assignments (id, category_id, target_type, target_id, author_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, ca.ID, ca.CategoryID, ca.TargetType, ca.TargetID, ca.AuthorID, ca.CreatedAt, ca.UpdatedAt)
	return err
}

func (r *PostgresCategoryAssignmentRepository) Unassign(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM category_assignments WHERE id = $1`, id)
	return err
}

func (r *PostgresCategoryAssignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CategoryAssignment, error) {
	row := r.db.QueryRow(ctx, `
			SELECT id, category_id, target_type, target_id, author_id, created_at, updated_at
			FROM category_assignments WHERE id = $1
		`, id)
	ca := &entity.CategoryAssignment{}

	err := row.Scan(&ca.ID, &ca.CategoryID, &ca.TargetType, &ca.TargetID, &ca.AuthorID, &ca.CreatedAt, &ca.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (r *PostgresCategoryAssignmentRepository) ListByEntity(ctx context.Context, targetID uuid.UUID, targetType string) ([]*entity.CategoryAssignment, error) {
	rows, err := r.db.Query(ctx, `
			SELECT id, category_id, target_type, target_id, author_id, created_at, updated_at
			FROM category_assignments
			WHERE target_id = $1 AND target_type = $2
		`, targetID, targetType)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*entity.CategoryAssignment{}
	for rows.Next() {
		ca := &entity.CategoryAssignment{}
		err := rows.Scan(&ca.ID, &ca.CategoryID, &ca.TargetType, &ca.TargetID, &ca.AuthorID, &ca.CreatedAt, &ca.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, ca)
	}

	return assignments, nil
}

func (r *PostgresCategoryAssignmentRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryAssignment, error) {
	rows, err := r.db.Query(ctx, `
			SELECT id, category_id, target_type, target_id, author_id, created_at, updated_at
			FROM category_assignments
			WHERE category_id = $1
		`, categoryID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*entity.CategoryAssignment{}
	for rows.Next() {
		ca := &entity.CategoryAssignment{}
		err := rows.Scan(&ca.ID, &ca.CategoryID, &ca.TargetType, &ca.TargetID, &ca.AuthorID, &ca.CreatedAt, &ca.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, ca)
	}

	return assignments, nil
}
