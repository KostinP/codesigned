// internal/lesson/wire.go
package lesson

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/lesson/repository"
	http "github.com/kostinp/edu-platform-backend/internal/lesson/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/lesson/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

var LessonSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresLessonRepository,
	wire.Bind(new(repository.LessonRepository), new(*repository.PostgresLessonRepository)),
	usecase.NewLessonUsecase,
	http.NewLessonHandler,
)
