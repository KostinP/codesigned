// internal/shared/middleware/abac.go
package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/abac"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
	"github.com/labstack/echo/v4"
)

func ABACMiddleware(engine *abac.Engine, resourceType, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr, ok := c.Get(UserIDKey).(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "auth required")
			}

			userID, _ := uuid.Parse(userIDStr)
			user := &entity.User{ID: userID} // потом можно подгружать полные данные

			// Подгружаем автора ресурса, если нужно
			resourceAuthorID := c.Get("resource_author_id")
			targetAuthorID := c.Get("target_author_id")

			ctx := abac.Context{
				User: user,
				Resource: map[string]interface{}{
					"type":             resourceType,
					"id":               c.Param("id"),
					"author_id":        resourceAuthorID,
					"target_author_id": targetAuthorID,
					"user_id":          userIDStr,
				},
				Action: action,
				Environment: map[string]interface{}{
					"time": map[string]int{
						"hour":    time.Now().Hour(),
						"weekday": int(time.Now().Weekday()),
					},
					"ip": c.RealIP(),
				},
			}

			allowed, err := engine.Evaluate(ctx)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "auth error")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			return next(c)
		}
	}
}
