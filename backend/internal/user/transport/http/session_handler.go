package transport

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	SessionUsecase usecase.SessionUsecase
}

func NewSessionHandler(sessionUC usecase.SessionUsecase) *SessionHandler {
	return &SessionHandler{SessionUsecase: sessionUC}
}

// SetInactivityTimeoutRequest — структура запроса для установки таймаута бездействия.
// @Description Таймаут в секундах.
type SetInactivityTimeoutRequest struct {
	// Таймаут в секундах
	TimeoutSeconds int64 `json:"timeout_seconds" example:"3600"`
}

// ListSessions возвращает все сессии пользователя
// @Summary Получить список сессий пользователя
// @Tags Session
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.UserSession
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /me/sessions [get]
func (h *SessionHandler) ListSessions(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	rawSessions, err := h.SessionUsecase.ListSessions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	sessions := make([]entity.UserSession, 0, len(rawSessions))
	for _, s := range rawSessions {
		sessions = append(sessions, *s) // разыменовываем указатель
	}

	return c.JSON(http.StatusOK, sessions)
}

// DeleteSession удаляет указанную сессию пользователя
// @Summary Удалить сессию по ID
// @Tags Session
// @Security BearerAuth
// @Param id path string true "ID сессии"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /me/sessions/{id} [delete]
func (h *SessionHandler) DeleteSession(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session id"})
	}

	err = h.SessionUsecase.DeleteSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete session"})
	}

	return c.NoContent(http.StatusNoContent)
}

// SetInactivityTimeout устанавливает таймаут неактивности для пользователя
// @Summary Установить таймаут неактивности
// @Tags Session
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body SetInactivityTimeoutRequest true "Таймаут в секундах"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /me/inactivity-timeout [post]
func (h *SessionHandler) SetInactivityTimeout(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	type Request struct {
		TimeoutSeconds int64 `json:"timeout_seconds"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if timeout < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "timeout must be positive"})
	}

	err = h.SessionUsecase.SetInactivityTimeout(c.Request().Context(), userID, timeout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to set inactivity timeout"})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetInactivityTimeout возвращает текущий таймаут неактивности пользователя
// @Summary Получить таймаут неактивности
// @Tags Session
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]int64
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /me/inactivity-timeout [get]
func (h *SessionHandler) GetInactivityTimeout(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	timeout, err := h.SessionUsecase.GetInactivityTimeout(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get inactivity timeout"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"timeout_seconds": int64(timeout.Seconds()),
	})
}
