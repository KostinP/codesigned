package transport

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/middleware"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
	"github.com/labstack/echo/v4"
)

type VisitorEventHandler struct {
	usecase usecase.VisitorEventUsecase
}

func NewVisitorEventHandler(uc usecase.VisitorEventUsecase) *VisitorEventHandler {
	return &VisitorEventHandler{usecase: uc}
}

type LogEventRequest struct {
	EventType string                 `json:"event_type" validate:"required"`
	EventData map[string]interface{} `json:"event_data"`
	VisitorID string                 `json:"visitor_id,omitempty"`
}

// @Summary Логировать событие посетителя
// @Tags Visitors
// @Accept json
// @Produce json
// @Param event body LogEventRequest true "Событие"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /visitor/events [post]
func (h *VisitorEventHandler) LogEvent(c echo.Context) error {
	var visitorID uuid.UUID

	// Сначала пробуем получить visitor_id из тела запроса
	req := new(LogEventRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный формат запроса"})
	}

	log.Printf("Received event: Type=%s, VisitorID from body=%s", req.EventType, req.VisitorID)

	// Валидация
	if err := c.Validate(req); err != nil {
		log.Printf("Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Получаем visitor_id из разных источников (приоритет: тело запроса > middleware > новый)
	if req.VisitorID != "" {
		// Пробуем использовать visitor_id из тела запроса
		parsedID, err := uuid.Parse(req.VisitorID)
		if err == nil {
			visitorID = parsedID
		} else {
			// Если невалидный UUID из тела, создаем новый
			visitorID = uuid.New()
		}
	} else {
		// Пробуем получить из middleware
		visitorIDRaw := c.Get(middleware.VisitorIDKey)
		if visitorIDRaw != nil {
			visitorIDStr, ok := visitorIDRaw.(string)
			if ok && visitorIDStr != "" {
				parsedID, err := uuid.Parse(visitorIDStr)
				if err == nil {
					visitorID = parsedID
				} else {
					visitorID = uuid.New()
				}
			} else {
				visitorID = uuid.New()
			}
		} else {
			// Если нет нигде, создаем новый
			visitorID = uuid.New()
		}
	}

	log.Printf("Processing event for visitor: %s, type: %s", visitorID.String(), req.EventType)

	// Если event_data nil, инициализируем пустую map
	if req.EventData == nil {
		req.EventData = make(map[string]interface{})
	}

	err := h.usecase.LogEvent(c.Request().Context(), visitorID, req.EventType, req.EventData)
	if err != nil {
		log.Printf("Failed to save event: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось сохранить событие"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":     "событие сохранено",
		"visitor_id": visitorID.String(),
	})
}
