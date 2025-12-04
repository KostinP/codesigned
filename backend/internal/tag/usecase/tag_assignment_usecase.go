package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/tag/entity"
	"github.com/kostinp/edu-platform-backend/internal/tag/repository"
)

type TagAssignmentUsecase interface {
	AssignTag(ctx context.Context, ta *entity.TagAssignment, authorID uuid.UUID) error
	RemoveAssignment(ctx context.Context, id uuid.UUID) error
	ListAssignmentsByEntity(ctx context.Context, entityID uuid.UUID, entityType string) ([]*entity.TagAssignment, error)
	ListAssignmentsByTag(ctx context.Context, tagID uuid.UUID) ([]*entity.TagAssignment, error)
}

type tagAssignmentUsecase struct {
	assignmentRepo repository.TagAssignmentRepository
}

func NewTagAssignmentUsecase(ar repository.TagAssignmentRepository) TagAssignmentUsecase {
	return &tagAssignmentUsecase{assignmentRepo: ar}
}

func (u *tagAssignmentUsecase) AssignTag(ctx context.Context, ta *entity.TagAssignment, authorID uuid.UUID) error {
	ta.Init(authorID)
	return u.assignmentRepo.Assign(ctx, ta) // Вызов Assign, не Create
}

func (u *tagAssignmentUsecase) RemoveAssignment(ctx context.Context, id uuid.UUID) error {
	return u.assignmentRepo.Unassign(ctx, id, uuid.Nil) // Для удаления по id, передаем entityID = uuid.Nil (будет игнорироваться)
}

func (u *tagAssignmentUsecase) ListAssignmentsByEntity(ctx context.Context, entityID uuid.UUID, entityType string) ([]*entity.TagAssignment, error) {
	return u.assignmentRepo.ListByEntity(ctx, entityID, entityType)
}

func (u *tagAssignmentUsecase) ListAssignmentsByTag(ctx context.Context, tagID uuid.UUID) ([]*entity.TagAssignment, error) {
	return u.assignmentRepo.ListByTag(ctx, tagID)
}
