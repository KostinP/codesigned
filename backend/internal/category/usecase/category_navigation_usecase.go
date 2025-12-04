// internal/category/usecase/category_navigation_usecase.go
package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/category/repository"
	searchEntity "github.com/kostinp/edu-platform-backend/internal/search/entity"
)

type CategoryNavigationUsecase interface {
	GetTree(ctx context.Context) ([]*entity.TreeCategory, error)
	GetBreadcrumbs(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error)
	GetContentByCategory(ctx context.Context, categoryID uuid.UUID, contentType string) ([]searchEntity.SearchEntity, error)
}

type categoryNavigationUsecase struct {
	repo repository.CategoryNavigationRepository
}

func NewCategoryNavigationUsecase(repo repository.CategoryNavigationRepository) CategoryNavigationUsecase {
	return &categoryNavigationUsecase{repo: repo}
}

func (u *categoryNavigationUsecase) GetTree(ctx context.Context) ([]*entity.TreeCategory, error) {
	return u.repo.GetTree(ctx)
}

func (u *categoryNavigationUsecase) GetBreadcrumbs(ctx context.Context, categoryID uuid.UUID) ([]*entity.Category, error) {
	return u.repo.GetBreadcrumbs(ctx, categoryID)
}

func (u *categoryNavigationUsecase) GetContentByCategory(ctx context.Context, categoryID uuid.UUID, contentType string) ([]searchEntity.SearchEntity, error) {
	return u.repo.GetContentByCategory(ctx, categoryID, contentType)
}
