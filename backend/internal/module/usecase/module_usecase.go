package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/module/entity"
	"github.com/kostinp/edu-platform-backend/internal/module/repository"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type ModuleUsecase interface {
	Create(ctx context.Context, module *entity.Module, authorID uuid.UUID) error
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Module, int, error)
}

type moduleUsecase struct {
	repo repository.ModuleRepository
}

func NewModuleUsecase(repo repository.ModuleRepository) ModuleUsecase {
	return &moduleUsecase{repo: repo}
}

func (u *moduleUsecase) Create(ctx context.Context, module *entity.Module, authorID uuid.UUID) error {
	module.Init(authorID)
	return u.repo.Create(ctx, module)
}

func (u *moduleUsecase) Update(ctx context.Context, module *entity.Module) error {
	module.Touch()
	return u.repo.Update(ctx, module)
}

func (u *moduleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *moduleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *moduleUsecase) List(ctx context.Context, pag pagination.Params) ([]*entity.Module, int, error) {
	return u.repo.List(ctx, pag)
}
