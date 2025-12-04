// internal/category/transport/http/category_navigation_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/category/usecase"
	"github.com/labstack/echo/v4"
)

type CategoryNavigationHandler struct {
	navigationUC usecase.CategoryNavigationUsecase
}

func NewCategoryNavigationHandler(uc usecase.CategoryNavigationUsecase) *CategoryNavigationHandler {
	return &CategoryNavigationHandler{navigationUC: uc}
}

// GetCategoryTree godoc
// @Summary Получить дерево категорий
// @Tags Categories
// @Produce json
// @Success 200 {object} CategoryTreeResponse
// @Failure 500 {object} map[string]string
// @Router /categories/tree [get]
func (h *CategoryNavigationHandler) GetCategoryTree(c echo.Context) error {
	tree, err := h.navigationUC.GetTree(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить дерево категорий",
		})
	}
	return c.JSON(http.StatusOK, CategoryTreeResponse{
		Categories: tree,
	})
}

// GetBreadcrumbs godoc
// @Summary Получить хлебные крошки для категории
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {array} entity.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/breadcrumbs [get]
func (h *CategoryNavigationHandler) GetBreadcrumbs(c echo.Context) error {
	idStr := c.Param("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректный ID категории",
		})
	}
	breadcrumbs, err := h.navigationUC.GetBreadcrumbs(c.Request().Context(), categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить хлебные крошки",
		})
	}
	return c.JSON(http.StatusOK, breadcrumbs)
}

// GetCategoryContent godoc
// @Summary Получить контент категории
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Param type query string false "Тип контента (course, lesson)" default(course)
// @Success 200 {array} searchEntity.SearchEntity
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/content [get]
func (h *CategoryNavigationHandler) GetCategoryContent(c echo.Context) error {
	idStr := c.Param("id")
	categoryID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректный ID категории",
		})
	}
	contentType := c.QueryParam("type")
	if contentType == "" {
		contentType = "course"
	}
	content, err := h.navigationUC.GetContentByCategory(c.Request().Context(), categoryID, contentType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить контент категории",
		})
	}
	return c.JSON(http.StatusOK, content)
}

// DTO для ответов
type CategoryTreeResponse struct {
	Categories []*entity.TreeCategory `json:"categories"`
}
