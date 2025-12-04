package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/tag/entity"
)

var ErrTagAssignmentNotFound = errors.New("tag assignment not found")

type TagAssignmentRepository interface {
	Assign(ctx context.Context, assignment *entity.TagAssignment) error
	Unassign(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error
	ListByEntity(ctx context.Context, entityID uuid.UUID, entityType string) ([]*entity.TagAssignment, error)
	ListByTag(ctx context.Context, tagID uuid.UUID) ([]*entity.TagAssignment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TagAssignment, error)
}

type PostgresTagAssignmentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTagAssignmentRepository(db *pgxpool.Pool) *PostgresTagAssignmentRepository {
	return &PostgresTagAssignmentRepository{db: db}
}

func (r *PostgresTagAssignmentRepository) Assign(ctx context.Context, ta *entity.TagAssignment) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tag_assignments (id, tag_id, entity_id, entity_type, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, ta.ID, ta.TagID, ta.EntityID, ta.EntityType, ta.AuthorID, ta.CreatedAt, ta.UpdatedAt)
	return err
}

func (r *PostgresTagAssignmentRepository) Unassign(ctx context.Context, id uuid.UUID, entityID uuid.UUID) error {
	var cmdAssign *entity.TagAssignment
	var err error

	if id != uuid.Nil {
		cmdAssign, err = r.GetByID(ctx, id)
	} else {
		// Если id не передан (uuid.Nil), можно реализовать удаление по entityID - при необходимости
		// В данном случае для удаления по id всегда id передается
		return errors.New("id is required for unassign")
	}
	if err != nil {
		return err
	}
	if cmdAssign == nil {
		return ErrTagAssignmentNotFound
	}

	_, err = r.db.Exec(ctx, `
		DELETE FROM tag_assignments WHERE id = $1
	`, id)
	return err
}

func (r *PostgresTagAssignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TagAssignment, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tag_id, entity_id, entity_type, author_id, created_at, updated_at
		FROM tag_assignments WHERE id = $1
	`, id)

	ta := &entity.TagAssignment{}
	err := row.Scan(&ta.ID, &ta.TagID, &ta.EntityID, &ta.EntityType, &ta.AuthorID, &ta.CreatedAt, &ta.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return ta, nil
}

func (r *PostgresTagAssignmentRepository) ListByEntity(ctx context.Context, entityID uuid.UUID, entityType string) ([]*entity.TagAssignment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tag_id, entity_id, entity_type, author_id, created_at, updated_at
		FROM tag_assignments
		WHERE entity_id = $1 AND entity_type = $2
	`, entityID, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*entity.TagAssignment{}
	for rows.Next() {
		ta := &entity.TagAssignment{}
		err := rows.Scan(&ta.ID, &ta.TagID, &ta.EntityID, &ta.EntityType, &ta.AuthorID, &ta.CreatedAt, &ta.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, ta)
	}
	return assignments, nil
}

func (r *PostgresTagAssignmentRepository) ListByTag(ctx context.Context, tagID uuid.UUID) ([]*entity.TagAssignment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tag_id, entity_id, entity_type, author_id, created_at, updated_at
		FROM tag_assignments
		WHERE tag_id = $1
	`, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*entity.TagAssignment{}
	for rows.Next() {
		ta := &entity.TagAssignment{}
		err := rows.Scan(&ta.ID, &ta.TagID, &ta.EntityID, &ta.EntityType, &ta.AuthorID, &ta.CreatedAt, &ta.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, ta)
	}
	return assignments, nil
}
