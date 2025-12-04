// internal/search/transport/http/search_handler.go
package http

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/search/entity"
	"github.com/kostinp/edu-platform-backend/internal/search/usecase"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	usecase usecase.SearchUsecase
}

func NewSearchHandler(uc usecase.SearchUsecase) *SearchHandler {
	return &SearchHandler{usecase: uc}
}

// Search godoc
// @Summary Поиск контента
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string false "Поисковый запрос"
// @Param category_id query string false "ID категории"
// @Param tags query string false "Теги через запятую"
// @Param content_type query string false "Типы контента через запятую (course,lesson,module)"
// @Param price_min query int false "Минимальная цена"
// @Param price_max query int false "Максимальная цена"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} entity.SearchResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
func (h *SearchHandler) Search(c echo.Context) error {
	var filters entity.SearchFilters
	// Парсим параметры запроса
	filters.Query = c.QueryParam("q")
	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID) // Исправлено: Parse для UUID
		if err == nil {
			filters.CategoryID = &id // Используем id
		}
	}
	if tags := c.QueryParam("tags"); tags != "" {
		// Парсим теги через запятую
		// filters.Tags = strings.Split(tags, ",")
	}
	if contentType := c.QueryParam("content_type"); contentType != "" {
		// filters.ContentType = strings.Split(contentType, ",")
	}
	if priceMin := c.QueryParam("price_min"); priceMin != "" {
		if val, err := strconv.Atoi(priceMin); err == nil {
			filters.PriceMin = &val
		}
	}
	if priceMax := c.QueryParam("price_max"); priceMax != "" {
		if val, err := strconv.Atoi(priceMax); err == nil {
			filters.PriceMax = &val
		}
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	filters.Limit = limit
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}
	filters.Offset = offset
	result, err := h.usecase.Search(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка поиска",
		})
	}
	return c.JSON(http.StatusOK, result)
}

// SearchAdvanced godoc
// @Summary Расширенный поиск
// @Tags Search
// @Accept json
// @Produce json
// @Param filters body entity.SearchFilters true "Фильтры поиска"
// @Success 200 {object} entity.SearchResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search/advanced [post]
func (h *SearchHandler) SearchAdvanced(c echo.Context) error {
	var filters entity.SearchFilters
	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректные параметры поиска",
		})
	}
	result, err := h.usecase.SearchAdvanced(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка расширенного поиска",
		})
	}
	return c.JSON(http.StatusOK, result)
}

// Autocomplete godoc
// @Summary Автодополнение поиска
// @Tags Search
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Param limit query int false "Лимит" default(10)
// @Success 200 {array} string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search/autocomplete [get]
func (h *SearchHandler) Autocomplete(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Параметр q обязателен",
		})
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	suggestions, err := h.usecase.GetAutocomplete(c.Request().Context(), query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка автодополнения",
		})
	}
	return c.JSON(http.StatusOK, suggestions)
}
