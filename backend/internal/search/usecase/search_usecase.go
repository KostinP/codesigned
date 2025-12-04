// internal/search/usecase/search_usecase.go
package usecase

import (
	"context"

	"github.com/kostinp/edu-platform-backend/internal/search/entity"
	"github.com/kostinp/edu-platform-backend/internal/search/repository"
)

type SearchUsecase interface {
	Search(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error)
	SearchAdvanced(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error)
	GetAutocomplete(ctx context.Context, query string, limit int) ([]string, error)
	GetPopularSearches(ctx context.Context, limit int) ([]string, error)
	LogSearch(ctx context.Context, query string, filters entity.SearchFilters, userID *string) error
}

type searchUsecase struct {
	repo repository.SearchRepository
}

func NewSearchUsecase(repo repository.SearchRepository) SearchUsecase {
	return &searchUsecase{repo: repo}
}

func (u *searchUsecase) Search(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error) {
	// Валидация параметров
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}
	return u.repo.Search(ctx, filters)
}

func (u *searchUsecase) SearchAdvanced(ctx context.Context, filters entity.SearchFilters) (*entity.SearchResult, error) {
	// Валидация параметров
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}
	return u.repo.SearchAdvanced(ctx, filters)
}

func (u *searchUsecase) GetAutocomplete(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return u.repo.GetAutocomplete(ctx, query, limit)
}

func (u *searchUsecase) GetPopularSearches(ctx context.Context, limit int) ([]string, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return u.repo.GetPopularSearches(ctx, limit)
}

func (u *searchUsecase) LogSearch(ctx context.Context, query string, filters entity.SearchFilters, userID *string) error {
	// TODO: Логировать поисковые запросы для аналитики
	return nil
}
