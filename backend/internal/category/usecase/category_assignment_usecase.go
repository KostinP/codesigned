// internal/category/usecase/category_assignment_usecase.go
package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/category/repository"
)

type CategoryAssignmentUsecase interface {
	Assign(ctx context.Context, ca *entity.CategoryAssignment, authorID uuid.UUID) error
	Remove(ctx context.Context, id uuid.UUID) error
	ListByEntity(ctx context.Context, targetID uuid.UUID, targetType string) ([]*entity.CategoryAssignment, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryAssignment, error)
}

type categoryAssignmentUsecase struct {
	repo repository.CategoryAssignmentRepository
}

func NewCategoryAssignmentUsecase(repo repository.CategoryAssignmentRepository) CategoryAssignmentUsecase {
	return &categoryAssignmentUsecase{repo: repo}
}

func (u *categoryAssignmentUsecase) Assign(ctx context.Context, ca *entity.CategoryAssignment, authorID uuid.UUID) error {
	ca.Init(authorID)
	return u.repo.Assign(ctx, ca)
}

func (u *categoryAssignmentUsecase) Remove(ctx context.Context, id uuid.UUID) error {
	return u.repo.Unassign(ctx, id)
}

func (u *categoryAssignmentUsecase) ListByEntity(ctx context.Context, targetID uuid.UUID, targetType string) ([]*entity.CategoryAssignment, error) {
	return u.repo.ListByEntity(ctx, targetID, targetType)
}

func (u *categoryAssignmentUsecase) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.CategoryAssignment, error) {
	return u.repo.ListByCategory(ctx, categoryID)
}
