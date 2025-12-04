package middleware

import (
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

func LinkVisitorWithUser(userUC *usecase.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			visitorIDRaw := c.Get(VisitorIDKey)
			userIDRaw := c.Get(UserIDKey)

			if visitorIDRaw == nil || userIDRaw == nil {
				// Логика: либо гость, либо неавторизованный пользователь — пропускаем
				return next(c)
			}

			visitorID, err := uuid.Parse(visitorIDRaw.(string))
			if err != nil {
				// Логировать ошибку можно, но не прерывать
				return next(c)
			}

			userID, ok := userIDRaw.(uuid.UUID)
			if !ok {
				// Логировать ошибку можно, но не прерывать
				return next(c)
			}

			// Пробуем связать visitor и user
			if err := userUC.LinkVisitorToUser(c.Request().Context(), userID, visitorID); err != nil {
				// Логгируем ошибку, но не блокируем запрос
				// Например: logger.Warnf("LinkVisitorWithUser: %v", err)
			}

			return next(c)
		}
	}
}
