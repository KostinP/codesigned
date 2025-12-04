package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/kostinp/edu-platform-backend/internal/tag/entity"
	"github.com/kostinp/edu-platform-backend/internal/tag/repository"
)

var ErrTagNotFound = errors.New("tag not found")

type TagUsecase interface {
	CreateTag(ctx context.Context, tag *entity.Tag, authorID uuid.UUID) error
	UpdateTag(ctx context.Context, tag *entity.Tag) error
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
	GetTagByID(ctx context.Context, tagID uuid.UUID) (*entity.Tag, error)
	ListTags(ctx context.Context, pag pagination.Params) ([]*entity.Tag, int, error)
}

type tagUsecase struct {
	tagRepo repository.TagRepository
}

func NewTagUsecase(tr repository.TagRepository) TagUsecase {
	return &tagUsecase{tagRepo: tr}
}

func (u *tagUsecase) CreateTag(ctx context.Context, tag *entity.Tag, authorID uuid.UUID) error {
	tag.Init(authorID)
	return u.tagRepo.Create(ctx, tag)
}

func (u *tagUsecase) UpdateTag(ctx context.Context, tag *entity.Tag) error {
	tag.Touch()
	return u.tagRepo.Update(ctx, tag)
}

func (u *tagUsecase) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	return u.tagRepo.Delete(ctx, tagID)
}

func (u *tagUsecase) GetTagByID(ctx context.Context, tagID uuid.UUID) (*entity.Tag, error) {
	return u.tagRepo.GetByID(ctx, tagID)
}

func (u *tagUsecase) ListTags(ctx context.Context, pag pagination.Params) ([]*entity.Tag, int, error) {
	return u.tagRepo.List(ctx, pag)
}
