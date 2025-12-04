// internal/category/transport/http/category_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/category/entity"
	"github.com/kostinp/edu-platform-backend/internal/category/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/dto"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryUC   usecase.CategoryUsecase
	assignmentUC usecase.CategoryAssignmentUsecase
}

func NewCategoryHandler(cuc usecase.CategoryUsecase, auc usecase.CategoryAssignmentUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUC: cuc, assignmentUC: auc}
}

// CreateCategory godoc
// @Summary Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body entity.Category true "Category object"
// @Success 201 {object} entity.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) Create(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	category := new(entity.Category)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.categoryUC.Create(c.Request().Context(), category, authorID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	category, err := h.categoryUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
	}
	return c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary Update category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body entity.Category true "Updated category object"
// @Success 200 {object} entity.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	category := new(entity.Category)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	category.ID = id
	if err := h.categoryUC.Update(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete category
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := h.categoryUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListCategories godoc
// @Summary List categories
// @Tags categories
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sort_by query string false "Sort by"
// @Param order query string false "Order"
// @Success 200 {object} dto.PaginatedResponse[*entity.Category]
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) List(c echo.Context) error {
	pagQuery := pagination.ParsePaginationParams(c)
	pag := pagQuery.ToDomainParams()
	categories, total, err := h.categoryUC.List(c.Request().Context(), pag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Category]{
		Items:  categories,
		Total:  total,
		Limit:  pag.Limit,
		Offset: pag.Offset,
	})
}

// AssignCategory godoc
// @Summary Assign category to entity
// @Tags category-assignments
// @Accept json
// @Produce json
// @Param assignment body entity.CategoryAssignment true "Assignment object"
// @Success 201 {object} entity.CategoryAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category-assignments [post]
func (h *CategoryHandler) Assign(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	ca := new(entity.CategoryAssignment)
	if err := c.Bind(ca); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.assignmentUC.Assign(c.Request().Context(), ca, authorID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, ca)
}

// RemoveAssignment godoc
// @Summary Remove category assignment
// @Tags category-assignments
// @Param id path string true "Assignment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category-assignments/{id} [delete]
func (h *CategoryHandler) Remove(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := h.assignmentUC.Remove(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListAssignmentsByEntity godoc
// @Summary List assignments by entity
// @Tags category-assignments
// @Produce json
// @Param target_id query string true "Target ID"
// @Param target_type query string true "Target Type"
// @Success 200 {array} entity.CategoryAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category-assignments/by-entity [get]
func (h *CategoryHandler) ListByEntity(c echo.Context) error {
	targetIDStr := c.QueryParam("target_id")
	targetType := c.QueryParam("target_type")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid target ID"})
	}
	assignments, err := h.assignmentUC.ListByEntity(c.Request().Context(), targetID, targetType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, assignments)
}

// ListAssignmentsByCategory godoc
// @Summary List assignments by category
// @Tags category-assignments
// @Produce json
// @Param category_id query string true "Category ID"
// @Success 200 {array} entity.CategoryAssignment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /category-assignments/by-category [get]
func (h *CategoryHandler) ListByCategory(c echo.Context) error {
	categoryIDStr := c.QueryParam("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category ID"})
	}
	assignments, err := h.assignmentUC.ListByCategory(c.Request().Context(), categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, assignments)
}
