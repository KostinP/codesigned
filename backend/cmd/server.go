// cmd/server.go
package main

import (
	category_http "github.com/kostinp/edu-platform-backend/internal/category/transport/http"
	category_navigation_http "github.com/kostinp/edu-platform-backend/internal/category/transport/http"
	course_http "github.com/kostinp/edu-platform-backend/internal/course/transport/http"
	lesson_http "github.com/kostinp/edu-platform-backend/internal/lesson/transport/http"
	module_http "github.com/kostinp/edu-platform-backend/internal/module/transport/http"
	search_http "github.com/kostinp/edu-platform-backend/internal/search/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/shared/abac"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/shared/middleware"
	customMiddleware "github.com/kostinp/edu-platform-backend/internal/shared/middleware"
	"github.com/kostinp/edu-platform-backend/internal/shared/validation"
	tag_http "github.com/kostinp/edu-platform-backend/internal/tag/transport/http"
	transport "github.com/kostinp/edu-platform-backend/internal/user/transport/http"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewServer создает и конфигурирует Echo-сервер
func newEchoServer(
	cfg *config.Config,
	userHandler *transport.UserHandler,
	visitorEventHandler *transport.VisitorEventHandler,
	telegramAuthHandler *transport.TelegramAuthHandler,
	sessionHandler *transport.SessionHandler,
	analyticsHandler *transport.AnalyticsHandler,
	sessionUsecase usecase.SessionUsecase,
	userService *usecase.UserService,
	abacEngine *abac.Engine,
	courseHandler *course_http.CourseHandler,
	moduleHandler *module_http.ModuleHandler,
	lessonHandler *lesson_http.LessonHandler,
	categoryHandler *category_http.CategoryHandler,
	tagHandler *tag_http.TagHandler,
	categoryNavigationHandler *category_navigation_http.CategoryNavigationHandler,
	searchHandler *search_http.SearchHandler,
) (*echo.Echo, error) {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()

	// CORS middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3001",
			"http://localhost:3000",
			"http://localhost:3002", // добавить на всякий случай
			"https://codesigned.ru",
			"chrome-extension://*",
		},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-User-ID",
			"X-Visitor-ID",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposeHeaders: []string{
			"X-Refresh-Token", // если используете обновление токенов
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Swagger - публичный, без JWT
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health - публичный, без JWT
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Prometheus метрики
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Ваш JWT middleware
	jwtMiddleware := customMiddleware.JWTMiddleware([]byte(cfg.JWT.Secret), sessionUsecase)

	// Другие middleware, которые нужны для всех запросов
	e.Use(customMiddleware.VisitorMiddleware)
	e.Use(customMiddleware.SetUserIDMiddleware)
	e.Use(customMiddleware.LinkVisitorWithUser(userService))

	// Публичные маршруты без JWT
	e.POST("/api/users/:user_id/link-visitor", userHandler.LinkVisitorToUser)
	e.GET("/api/visitor", transport.GetVisitorIDHandler)
	e.POST("/api/visitor/events", visitorEventHandler.LogEvent)
	e.POST("/api/telegram/auth", telegramAuthHandler.Auth)

	// Создаем группу для маршрутов, защищённых JWT
	apiProtected := e.Group("/api")
	apiProtected.Use(jwtMiddleware)

	// ========== КАТЕГОРИИ (навигация) ==========
	// Публичные роуты - доступны всем
	e.GET("/api/categories/tree", categoryNavigationHandler.GetCategoryTree)
	e.GET("/api/categories/:id/breadcrumbs", categoryNavigationHandler.GetBreadcrumbs)
	e.GET("/api/categories/:id/content", categoryNavigationHandler.GetCategoryContent)

	// Поиск - доступен всем
	e.GET("/api/search", searchHandler.Search)
	e.GET("/api/search/autocomplete", searchHandler.Autocomplete)

	// Защищенные роуты категорий
	apiProtected.POST("/categories", middleware.ABACMiddleware(abacEngine, "category", "create")(categoryHandler.Create))
	apiProtected.GET("/categories", middleware.ABACMiddleware(abacEngine, "category", "read")(categoryHandler.List))
	apiProtected.GET("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "read")(categoryHandler.Get))
	apiProtected.PUT("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "update")(categoryHandler.Update))
	apiProtected.DELETE("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "delete")(categoryHandler.Delete))

	// Убираем RequireRole — заменяем на ABAC
	apiProtected.GET("/me/sessions",
		middleware.ABACMiddleware(abacEngine, "user_sessions", "read")(sessionHandler.ListSessions))
	apiProtected.DELETE("/me/sessions/:id",
		middleware.ABACMiddleware(abacEngine, "user_sessions", "delete")(sessionHandler.DeleteSession))
	apiProtected.GET("/analytics/page-views",
		middleware.ABACMiddleware(abacEngine, "analytics", "read")(analyticsHandler.GetPageViews))

	// Для курсов
	apiProtected.POST("/courses", middleware.ABACMiddleware(abacEngine, "course", "create")(courseHandler.Create))
	apiProtected.GET("/courses", middleware.ABACMiddleware(abacEngine, "course", "read")(courseHandler.List))
	apiProtected.GET("/courses/:id", middleware.ABACMiddleware(abacEngine, "course", "read")(courseHandler.Get))
	apiProtected.PUT("/courses/:id", middleware.ABACMiddleware(abacEngine, "course", "update")(courseHandler.Update))
	apiProtected.DELETE("/courses/:id", middleware.ABACMiddleware(abacEngine, "course", "delete")(courseHandler.Delete))

	// Для модулей
	apiProtected.POST("/modules", middleware.ABACMiddleware(abacEngine, "module", "create")(moduleHandler.Create))
	apiProtected.GET("/modules", middleware.ABACMiddleware(abacEngine, "module", "read")(moduleHandler.List))
	apiProtected.GET("/modules/:id", middleware.ABACMiddleware(abacEngine, "module", "read")(moduleHandler.Get))
	apiProtected.PUT("/modules/:id", middleware.ABACMiddleware(abacEngine, "module", "update")(moduleHandler.Update))
	apiProtected.DELETE("/modules/:id", middleware.ABACMiddleware(abacEngine, "module", "delete")(moduleHandler.Delete))

	// Для уроков
	apiProtected.POST("/lessons", middleware.ABACMiddleware(abacEngine, "lesson", "create")(lessonHandler.Create))
	apiProtected.GET("/lessons", middleware.ABACMiddleware(abacEngine, "lesson", "read")(lessonHandler.List))
	apiProtected.GET("/lessons/:id", middleware.ABACMiddleware(abacEngine, "lesson", "read")(lessonHandler.Get))
	apiProtected.PUT("/lessons/:id", middleware.ABACMiddleware(abacEngine, "lesson", "update")(lessonHandler.Update))
	apiProtected.DELETE("/lessons/:id", middleware.ABACMiddleware(abacEngine, "lesson", "delete")(lessonHandler.Delete))

	// Для категорий
	apiProtected.POST("/categories", middleware.ABACMiddleware(abacEngine, "category", "create")(categoryHandler.Create))
	apiProtected.GET("/categories", middleware.ABACMiddleware(abacEngine, "category", "read")(categoryHandler.List))
	apiProtected.GET("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "read")(categoryHandler.Get))
	apiProtected.PUT("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "update")(categoryHandler.Update))
	apiProtected.DELETE("/categories/:id", middleware.ABACMiddleware(abacEngine, "category", "delete")(categoryHandler.Delete))
	apiProtected.POST("/category-assignments", middleware.ABACMiddleware(abacEngine, "category_assignment", "create")(categoryHandler.Assign))
	apiProtected.DELETE("/category-assignments/:id", middleware.ABACMiddleware(abacEngine, "category_assignment", "delete")(categoryHandler.Remove))
	apiProtected.GET("/category-assignments/by-entity", middleware.ABACMiddleware(abacEngine, "category_assignment", "read")(categoryHandler.ListByEntity))
	apiProtected.GET("/category-assignments/by-category", middleware.ABACMiddleware(abacEngine, "category_assignment", "read")(categoryHandler.ListByCategory))

	// Для тегов
	apiProtected.POST("/tags", middleware.ABACMiddleware(abacEngine, "tag", "create")(tagHandler.CreateTag))
	apiProtected.GET("/tags", middleware.ABACMiddleware(abacEngine, "tag", "read")(tagHandler.ListTags))
	apiProtected.GET("/tags/:id", middleware.ABACMiddleware(abacEngine, "tag", "read")(tagHandler.GetTagByID))

	// apiProtected.PUT("/tags/:id", middleware.ABACMiddleware(abacEngine, "tag", "update")(tagHandler.UpdateTag))
	apiProtected.DELETE("/tags/:id", middleware.ABACMiddleware(abacEngine, "tag", "delete")(tagHandler.DeleteTag))
	apiProtected.POST("/tag-assignments", middleware.ABACMiddleware(abacEngine, "tag_assignment", "create")(tagHandler.AssignTag))
	apiProtected.DELETE("/tag-assignments/:id", middleware.ABACMiddleware(abacEngine, "tag_assignment", "delete")(tagHandler.RemoveAssignment))
	apiProtected.GET("/tag-assignments/by-entity", middleware.ABACMiddleware(abacEngine, "tag_assignment", "read")(tagHandler.ListAssignmentsByEntity))
	apiProtected.GET("/tag-assignments/by-tag", middleware.ABACMiddleware(abacEngine, "tag_assignment", "read")(tagHandler.ListAssignmentsByTag))

	// Добавьте сюда все защищённые роуты, например:
	apiProtected.GET("/me/sessions", sessionHandler.ListSessions)
	apiProtected.DELETE("/me/sessions/:id", sessionHandler.DeleteSession)
	apiProtected.POST("/me/inactivity-timeout", sessionHandler.SetInactivityTimeout)
	apiProtected.GET("/me/inactivity-timeout", sessionHandler.GetInactivityTimeout)

	// Аналитика
	apiProtected.GET("/analytics/page-views", analyticsHandler.GetPageViews)
	apiProtected.GET("/analytics/utm-stats", analyticsHandler.GetUTMStats)
	apiProtected.GET("/analytics/conversion-rate", analyticsHandler.GetConversionRate)
	return e, nil
}
