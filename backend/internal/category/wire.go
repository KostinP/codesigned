package category

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/category/repository"
	http "github.com/kostinp/edu-platform-backend/internal/category/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/category/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

var CategorySet = wire.NewSet(
	db.ConnectPostgres,

	// Репозитории
	repository.NewPostgresCategoryRepository,
	wire.Bind(new(repository.CategoryRepository), new(*repository.PostgresCategoryRepository)),

	repository.NewPostgresCategoryAssignmentRepository,
	wire.Bind(new(repository.CategoryAssignmentRepository), new(*repository.PostgresCategoryAssignmentRepository)),

	repository.NewPostgresCategoryNavigationRepository,
	wire.Bind(new(repository.CategoryNavigationRepository), new(*repository.PostgresCategoryNavigationRepository)),

	// Usecases
	usecase.NewCategoryUsecase,
	usecase.NewCategoryAssignmentUsecase,
	usecase.NewCategoryNavigationUsecase,

	// Handlers
	http.NewCategoryHandler,
	http.NewCategoryNavigationHandler,
)
