// internal/shared/middleware/resource_author.go
package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SetResourceAuthorMiddleware(repo interface {
	GetAuthorID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idStr := c.Param("id")
			id, _ := uuid.Parse(idStr)
			authorID, err := repo.GetAuthorID(c.Request().Context(), id)
			if err == nil {
				c.Set("resource_author_id", authorID.String())
			}
			return next(c)
		}
	}
}
