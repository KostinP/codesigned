// internal/shared/middleware/set_user_id.go
package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-User-ID")
		if userID == "" {
			return next(c)
		}

		id, err := uuid.Parse(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		}

		c.Set(UserIDKey, id)
		return next(c)
	}
}
