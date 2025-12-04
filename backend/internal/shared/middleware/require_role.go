package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

func RequireRole(userService *usecase.UserService, allowedRoles ...entity.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDRaw := c.Get(UserIDKey)
			if userIDRaw == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user_id missing in context")
			}

			userID, ok := userIDRaw.(string)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "invalid user_id type")
			}

			parsedID, err := uuid.Parse(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id format")
			}

			user, err := userService.GetUserByID(c.Request().Context(), parsedID)
			if err != nil || user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			for _, role := range allowedRoles {
				if user.Role == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}
