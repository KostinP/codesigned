package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtSecret []byte, sessionUC usecase.SessionUsecase) echo.MiddlewareFunc {
	const refreshThreshold = 24 * time.Hour // обновляем если токен истекает менее чем через 1 день

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			sessionIDStr, ok := claims["session_id"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "session_id missing")
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "user_id missing")
			}

			expFloat, ok := claims["exp"].(float64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "exp missing")
			}

			expiration := time.Unix(int64(expFloat), 0)
			now := time.Now()
			if now.After(expiration) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			sessionID, err := uuid.Parse(sessionIDStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid session_id format")
			}

			session, err := sessionUC.GetSessionByID(c.Request().Context(), sessionID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "session not found or expired")
			}

			if now.After(session.ExpiresAt) {
				// Удаляем просроченную сессию из БД
				_ = sessionUC.DeleteSession(c.Request().Context(), session.UserID, session.ID)
				return echo.NewHTTPError(http.StatusUnauthorized, "session expired")
			}

			// Обновляем last_active_at
			if err := sessionUC.UpdateLastActive(c.Request().Context(), sessionID); err != nil {
				c.Logger().Errorf("failed to update session last active: %v", err)
			}

			// Если токен истекает менее чем через refreshThreshold, создаём новый
			timeLeft := expiration.Sub(now)
			if timeLeft < refreshThreshold {
				newExpiration := now.AddDate(0, 6, 0)
				newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id":    userIDStr,
					"session_id": sessionIDStr,
					"exp":        newExpiration.Unix(),
					"iat":        now.Unix(),
				})

				newTokenString, err := newToken.SignedString(jwtSecret)
				if err != nil {
					c.Logger().Errorf("failed to sign new token: %v", err)
				} else {
					c.Response().Header().Set("X-Refresh-Token", newTokenString)
				}
			}

			c.Set("user_id", userIDStr)
			c.Set("session_id", sessionIDStr)

			return next(c)
		}
	}
}
