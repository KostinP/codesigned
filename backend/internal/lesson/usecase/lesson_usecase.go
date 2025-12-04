package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/lesson/entity"
	"github.com/kostinp/edu-platform-backend/internal/lesson/repository"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
)

type LessonUsecase interface {
	Create(ctx context.Context, lesson *entity.Lesson, authorID uuid.UUID) error
	Update(ctx context.Context, lesson *entity.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error)
	List(ctx context.Context, pag pagination.Params) ([]*entity.Lesson, int, error)
}

type lessonUsecase struct {
	repo repository.LessonRepository
}

func NewLessonUsecase(repo repository.LessonRepository) LessonUsecase {
	return &lessonUsecase{repo: repo}
}

func (u *lessonUsecase) Create(ctx context.Context, lesson *entity.Lesson, authorID uuid.UUID) error {
	lesson.Init(authorID)
	return u.repo.Create(ctx, lesson)
}

func (u *lessonUsecase) Update(ctx context.Context, lesson *entity.Lesson) error {
	lesson.Touch()
	return u.repo.Update(ctx, lesson)
}

func (u *lessonUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *lessonUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *lessonUsecase) List(ctx context.Context, pag pagination.Params) ([]*entity.Lesson, int, error) {
	return u.repo.List(ctx, pag)
}
