// internal/tag/wire.go
package tag

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
	"github.com/kostinp/edu-platform-backend/internal/tag/repository"
	http "github.com/kostinp/edu-platform-backend/internal/tag/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/tag/usecase"
)

var TagSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresTagRepository,
	wire.Bind(new(repository.TagRepository), new(*repository.PostgresTagRepository)),
	usecase.NewTagUsecase,
	repository.NewPostgresTagAssignmentRepository,
	usecase.NewTagAssignmentUsecase,
	http.NewTagHandler,
)
