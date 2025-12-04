// internal/module/wire.go
package module

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/module/repository"
	http "github.com/kostinp/edu-platform-backend/internal/module/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/module/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

var ModuleSet = wire.NewSet(
	db.ConnectPostgres,
	repository.NewPostgresModuleRepository,
	wire.Bind(new(repository.ModuleRepository), new(*repository.PostgresModuleRepository)),
	usecase.NewModuleUsecase,
	http.NewModuleHandler,
)
