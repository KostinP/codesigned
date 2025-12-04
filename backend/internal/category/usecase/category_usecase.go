package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/category/repository"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type CategoryUsecase interface {
	Create(ctx context.Context, category *entity.Category, authorID uuid.UUID) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Category, int, error)
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (u *categoryUsecase) Create(ctx context.Context, category *entity.Category, authorID uuid.UUID) error {
	category.ID = uuid.New()
	category.AuthorID = authorID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return u.repo.Create(ctx, category)
}

func (u *categoryUsecase) Update(ctx context.Context, category *entity.Category) error {
	category.UpdatedAt = time.Now()
	return u.repo.Update(ctx, category)
}

func (u *categoryUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *categoryUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *categoryUsecase) List(ctx context.Context, pag pagination.Params) ([]*entity.Category, int, error) {
	return u.repo.List(ctx, pag)
}
