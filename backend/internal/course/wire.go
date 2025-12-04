// internal/course/wire.go
package course

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/course/repository"
	http "github.com/kostinp/edu-platform-backend/internal/course/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/course/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

var CourseSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresCourseRepository,
	wire.Bind(new(repository.CourseRepository), new(*repository.PostgresCourseRepository)),
	usecase.NewCourseUsecase,
	http.NewCourseHandler,
)
