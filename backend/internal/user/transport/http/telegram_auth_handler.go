package transport

import (
	"net"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/shared/geo"
	"github.com/kostinp/edu-platform-backend/internal/shared/telegram"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

type TelegramAuthHandler struct {
	userService *usecase.UserService
	botToken    string
	jwtSecret   []byte
}

func NewTelegramAuthHandler(userService *usecase.UserService, botToken config.BotToken, jwtSecret config.JwtSecret) *TelegramAuthHandler {
	return &TelegramAuthHandler{
		userService: userService,
		botToken:    string(botToken),
		jwtSecret:   []byte(jwtSecret),
	}
}

// @Summary Авторизация через Telegram
// @Tags Telegram
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /telegram/auth [post]
func (h *TelegramAuthHandler) Auth(c echo.Context) error {
	authData := telegram.ParseTelegramAuth(c.Request())

	if !telegram.VerifyTelegramAuth(authData, h.botToken) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid telegram auth"})
	}

	ctx := c.Request().Context()

	user, err := h.userService.GetByTelegramID(ctx, authData.ID)
	if err != nil {
		user, err = h.userService.CreateFromTelegramAuth(ctx, authData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать пользователя"})
		}
	}

	// Получаем данные сессии
	ip := getIP(c)
	userAgent := c.Request().UserAgent()
	country, city := geo.Lookup(ip)

	// Создаем сессию пользователя в БД
	session, err := h.userService.CreateUserSession(ctx, user.ID, ip, userAgent, country, city)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать сессию"})
	}

	// По умолчанию — 6 месяцев
	expiration := time.Now().AddDate(0, 6, 0)

	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID.String(),
		"session_id": session.ID.String(),
		"exp":        expiration.Unix(),
		"iat":        time.Now().Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось создать токен"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token":      tokenString,
		"user_id":    user.ID.String(),
		"session_id": session.ID.String(),
		"username":   user.Username,
		"full_name":  user.FullName,
		"role":       user.Role,
	})
}

// Получение IP-адреса из запроса
func getIP(c echo.Context) string {
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(c.Request().RemoteAddr)
	}
	return ip
}
