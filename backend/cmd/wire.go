//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kostinp/edu-platform-backend/internal/category"
	"github.com/kostinp/edu-platform-backend/internal/course"
	"github.com/kostinp/edu-platform-backend/internal/lesson"
	"github.com/kostinp/edu-platform-backend/internal/module"
	"github.com/kostinp/edu-platform-backend/internal/search"
	"github.com/kostinp/edu-platform-backend/internal/shared/abac"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/user"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	echo "github.com/labstack/echo/v4"
)

// Главный wire-компонент
func InitializeServer(cfg *config.Config, abacEngine *abac.Engine) (*echo.Echo, error) {
	wire.Build(
		user.UserSet,
		course.CourseSet,
		module.ModuleSet,
		lesson.LessonSet,
		category.CategorySet,
		search.SearchSet,
		newEchoServer,
	)
	return nil, nil
}

// Для middleware и background job'ов
func InitializeSessionUsecase(cfg *config.Config) (usecase.SessionUsecase, error) {
	wire.Build(user.SessionUsecaseSet)
	return nil, nil
}
