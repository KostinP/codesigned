// internal/lesson/transport/http/lesson_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/lesson/entity"
	"github.com/kostinp/edu-platform-backend/internal/lesson/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/dto"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/labstack/echo/v4"
)

type LessonHandler struct {
	usecase usecase.LessonUsecase
}

func NewLessonHandler(uc usecase.LessonUsecase) *LessonHandler {
	return &LessonHandler{usecase: uc}
}

// CreateLesson godoc
// @Summary Create a new lesson
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body entity.Lesson true "Lesson object"
// @Success 201 {object} entity.Lesson
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lessons [post]
func (h *LessonHandler) Create(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	lesson := new(entity.Lesson)
	if err := c.Bind(lesson); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.usecase.Create(c.Request().Context(), lesson, authorID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, lesson)
}

// GetLesson godoc
// @Summary Get lesson by ID
// @Tags lessons
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} entity.Lesson
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lessons/{id} [get]
func (h *LessonHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	lesson, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "lesson not found"})
	}
	return c.JSON(http.StatusOK, lesson)
}

// UpdateLesson godoc
// @Summary Update lesson
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path string true "Lesson ID"
// @Param lesson body entity.Lesson true "Updated lesson object"
// @Success 200 {object} entity.Lesson
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lessons/{id} [put]
func (h *LessonHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	lesson := new(entity.Lesson)
	if err := c.Bind(lesson); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	lesson.ID = id
	if err := h.usecase.Update(c.Request().Context(), lesson); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, lesson)
}

// DeleteLesson godoc
// @Summary Delete lesson
// @Tags lessons
// @Param id path string true "Lesson ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lessons/{id} [delete]
func (h *LessonHandler) Delete(c echo.Context) error {
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

// ListLessons godoc
// @Summary List lessons
// @Tags lessons
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sort_by query string false "Sort by"
// @Param order query string false "Order"
// @Success 200 {object} dto.PaginatedResponse[*entity.Lesson]
// @Failure 500 {object} map[string]string
// @Router /lessons [get]
func (h *LessonHandler) List(c echo.Context) error {
	pagQuery := pagination.ParsePaginationParams(c)
	pag := pagQuery.ToDomainParams()
	lessons, total, err := h.usecase.List(c.Request().Context(), pag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Lesson]{
		Items:  lessons,
		Total:  total,
		Limit:  pag.Limit,
		Offset: pag.Offset,
	})
}
