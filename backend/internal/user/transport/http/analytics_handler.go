package transport

import (
	"net/http"
	"strconv"

	"github.com/kostinp/edu-platform-backend/internal/user/repository"
	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	repo *repository.ClickHouseVisitorEventRepo
}

func NewAnalyticsHandler(repo *repository.ClickHouseVisitorEventRepo) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

type AnalyticsResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total,omitempty"`
}

// @Summary Получить статистику просмотров страниц
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param days query int false "Количество дней (по умолчанию 30)" default(30)
// @Success 200 {object} AnalyticsResponse
// @Router /api/analytics/page-views [get]
func (h *AnalyticsHandler) GetPageViews(c echo.Context) error {
	daysStr := c.QueryParam("days")
	days := 30
	if daysStr != "" {
		if h.repo == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Analytics disabled",
			})
		}
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.repo.GetPageViews(c.Request().Context(), days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить статистику просмотров",
		})
	}

	return c.JSON(http.StatusOK, AnalyticsResponse{
		Data:  stats,
		Total: len(stats),
	})
}

// @Summary Получить UTM статистику
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param days query int false "Количество дней (по умолчанию 30)" default(30)
// @Success 200 {object} AnalyticsResponse
// @Router /api/analytics/utm-stats [get]
func (h *AnalyticsHandler) GetUTMStats(c echo.Context) error {
	daysStr := c.QueryParam("days")
	days := 30
	if daysStr != "" {
		if h.repo == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Analytics disabled",
			})
		}
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.repo.GetUTMStats(c.Request().Context(), days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить UTM статистику",
		})
	}

	return c.JSON(http.StatusOK, AnalyticsResponse{
		Data:  stats,
		Total: len(stats),
	})
}

// @Summary Получить конверсию по событию
// @Tags Analytics
// @Security BearerAuth
// @Produce json
// @Param event_type query string true "Тип события (form_submit, button_click, etc.)"
// @Param days query int false "Количество дней (по умолчанию 30)" default(30)
// @Success 200 {object} map[string]float64
// @Router /api/analytics/conversion-rate [get]
func (h *AnalyticsHandler) GetConversionRate(c echo.Context) error {
	eventType := c.QueryParam("event_type")
	if eventType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Параметр event_type обязателен",
		})
	}

	daysStr := c.QueryParam("days")
	days := 30
	if daysStr != "" {
		if h.repo == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Analytics disabled",
			})
		}
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	rate, err := h.repo.GetConversionRate(c.Request().Context(), eventType, days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить конверсию",
		})
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"conversion_rate": rate,
	})
}
