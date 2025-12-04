package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/course/entity"
	"github.com/kostinp/edu-platform-backend/internal/course/repository"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type CourseUsecase interface {
	Create(ctx context.Context, course *entity.Course, authorID uuid.UUID) error
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Course, int, error)
}

type courseUsecase struct {
	repo repository.CourseRepository
}

func NewCourseUsecase(repo repository.CourseRepository) CourseUsecase {
	return &courseUsecase{repo: repo}
}

func (u *courseUsecase) Create(ctx context.Context, course *entity.Course, authorID uuid.UUID) error {
	course.Init(authorID)
	return u.repo.Create(ctx, course)
}

func (u *courseUsecase) Update(ctx context.Context, course *entity.Course) error {
	course.Touch()
	return u.repo.Update(ctx, course)
}

func (u *courseUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *courseUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *courseUsecase) List(ctx context.Context, pag pagination.Params) ([]*entity.Course, int, error) {
	return u.repo.List(ctx, pag)
}
