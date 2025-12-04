// internal/module/transport/http/module_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/module/entity"
	"github.com/kostinp/edu-platform-backend/internal/module/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/dto"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/labstack/echo/v4"
)

type ModuleHandler struct {
	usecase usecase.ModuleUsecase
}

func NewModuleHandler(uc usecase.ModuleUsecase) *ModuleHandler {
	return &ModuleHandler{usecase: uc}
}

// CreateModule godoc
// @Summary Create a new module
// @Tags modules
// @Accept json
// @Produce json
// @Param module body entity.Module true "Module object"
// @Success 201 {object} entity.Module
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /modules [post]
func (h *ModuleHandler) Create(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	module := new(entity.Module)
	if err := c.Bind(module); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.usecase.Create(c.Request().Context(), module, authorID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, module)
}

// GetModule godoc
// @Summary Get module by ID
// @Tags modules
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} entity.Module
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /modules/{id} [get]
func (h *ModuleHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	module, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "module not found"})
	}
	return c.JSON(http.StatusOK, module)
}

// UpdateModule godoc
// @Summary Update module
// @Tags modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param module body entity.Module true "Updated module object"
// @Success 200 {object} entity.Module
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /modules/{id} [put]
func (h *ModuleHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	module := new(entity.Module)
	if err := c.Bind(module); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	module.ID = id
	if err := h.usecase.Update(c.Request().Context(), module); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, module)
}

// DeleteModule godoc
// @Summary Delete module
// @Tags modules
// @Param id path string true "Module ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /modules/{id} [delete]
func (h *ModuleHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListModules godoc
// @Summary List modules
// @Tags modules
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sort_by query string false "Sort by"
// @Param order query string false "Order"
// @Success 200 {object} dto.PaginatedResponse[*entity.Module]
// @Failure 500 {object} map[string]string
// @Router /modules [get]
func (h *ModuleHandler) List(c echo.Context) error {
	pagQuery := pagination.ParsePaginationParams(c)
	pag := pagQuery.ToDomainParams()
	modules, total, err := h.usecase.List(c.Request().Context(), pag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Module]{
		Items:  modules,
		Total:  total,
		Limit:  pag.Limit,
		Offset: pag.Offset,
	})
}
