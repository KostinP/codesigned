package search

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/search/repository"
	http "github.com/kostinp/edu-platform-backend/internal/search/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/search/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

var SearchSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresSearchRepository,
	wire.Bind(new(repository.SearchRepository), new(*repository.PostgresSearchRepository)),
	usecase.NewSearchUsecase,
	http.NewSearchHandler,
)
