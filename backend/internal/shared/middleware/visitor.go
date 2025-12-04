package middleware

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const VisitorIDKey = "visitor_id"

func VisitorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Пробуем получить visitor_id из cookie
		cookie, err := c.Cookie(VisitorIDKey)
		var visitorID string

		if err != nil || cookie == nil || cookie.Value == "" {
			// Создаем новый visitor_id
			visitorID = uuid.New().String()
			cookie = &http.Cookie{
				Name:     VisitorIDKey,
				Value:    visitorID,
				Path:     "/",
				HttpOnly: false,
				MaxAge:   365 * 24 * 60 * 60,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			}
			c.SetCookie(cookie)
			log.Printf("Created new visitor ID: %s", visitorID)
		} else {
			visitorID = cookie.Value
			log.Printf("Using existing visitor ID: %s", visitorID)
		}

		c.Set(VisitorIDKey, visitorID)
		return next(c)
	}
}
