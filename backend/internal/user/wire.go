//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	http "github.com/kostinp/edu-platform-backend/internal/user/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
)

// PrimaryVisitorEventRepoSet - набор для основного репозитория (ClickHouse)
var PrimaryVisitorEventRepoSet = wire.NewSet(
	repository.NewClickHouseVisitorEventRepo,
	wire.Bind(new(repository.PrimaryVisitorEventRepo), new(*repository.ClickHouseVisitorEventRepo)),
	wire.Bind(new(repository.AnalyticsRepository), new(*repository.ClickHouseVisitorEventRepo)),
)

// VisitorEventRepoSet - набор без fallback
var VisitorEventRepoSet = wire.NewSet(
	PrimaryVisitorEventRepoSet,
	repository.NewVisitorEventRepoWithFallback,
	wire.Bind(new(repository.VisitorEventRepository), new(*repository.VisitorEventRepoWithFallback)),
)

// UserRepoSet - набор для пользовательских репозиториев
var UserRepoSet = wire.NewSet(
	repository.NewPostgresUserRepository,
	wire.Bind(new(usecase.UserRepository), new(*repository.PostgresUserRepository)),
)

// SessionRepoSet - набор для сессий
var SessionRepoSet = wire.NewSet(
	repository.NewPostgresSessionRepository,
	wire.Bind(new(repository.SessionRepository), new(*repository.PostgresSessionRepository)),
)

var UserSet = wire.NewSet(
	// DB
	db.ConnectPostgres,
	ProvideClickHouseConn,
	// --- Repositories ---
	UserRepoSet,
	VisitorEventRepoSet,
	SessionRepoSet,
	// --- Usecases ---
	usecase.NewSessionUsecase,
	wire.Bind(new(usecase.SessionUsecase), new(*usecase.SessionUsecaseImpl)),
	usecase.NewUserService,
	usecase.NewVisitorEventUsecase,
	// --- Handlers ---
	http.NewUserHandler,
	http.NewVisitorEventHandler,
	http.NewSessionHandler,
	http.NewAnalyticsHandler,
	// --- Telegram Auth ---
	ProvideBotToken,
	ProvideJwtSecret,
	http.NewTelegramAuthHandler,
)

// SessionUsecaseSet - отдельный набор для middleware и фоновых задач
var SessionUsecaseSet = wire.NewSet(
	db.ConnectPostgres,
	SessionRepoSet,
	usecase.NewSessionUsecase,
	wire.Bind(new(usecase.SessionUsecase), new(*usecase.SessionUsecaseImpl)),
)
