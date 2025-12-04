// internal/course/transport/http/course_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/course/entity"
	"github.com/kostinp/edu-platform-backend/internal/course/usecase"
	"github.com/kostinp/edu-platform-backend/internal/shared/dto"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/labstack/echo/v4"
)

type CourseHandler struct {
	usecase usecase.CourseUsecase
}

func NewCourseHandler(uc usecase.CourseUsecase) *CourseHandler {
	return &CourseHandler{usecase: uc}
}

// CreateCourse godoc
// @Summary Create a new course
// @Tags courses
// @Accept json
// @Produce json
// @Param course body entity.Course true "Course object"
// @Success 201 {object} entity.Course
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /courses [post]
func (h *CourseHandler) Create(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	course := new(entity.Course)
	if err := c.Bind(course); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.usecase.Create(c.Request().Context(), course, authorID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, course)
}

// GetCourse godoc
// @Summary Get course by ID
// @Tags courses
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} entity.Course
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /courses/{id} [get]
func (h *CourseHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	course, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "course not found"})
	}
	return c.JSON(http.StatusOK, course)
}

// UpdateCourse godoc
// @Summary Update course
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param course body entity.Course true "Updated course object"
// @Success 200 {object} entity.Course
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /courses/{id} [put]
func (h *CourseHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	course := new(entity.Course)
	if err := c.Bind(course); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	course.ID = id
	if err := h.usecase.Update(c.Request().Context(), course); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, course)
}

// DeleteCourse godoc
// @Summary Delete course
// @Tags courses
// @Param id path string true "Course ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /courses/{id} [delete]
func (h *CourseHandler) Delete(c echo.Context) error {
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

// ListCourses godoc
// @Summary List courses
// @Tags courses
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sort_by query string false "Sort by"
// @Param order query string false "Order"
// @Success 200 {object} dto.PaginatedResponse[*entity.Course]
// @Failure 500 {object} map[string]string
// @Router /courses [get]
func (h *CourseHandler) List(c echo.Context) error {
	pagQuery := pagination.ParsePaginationParams(c)
	pag := pagQuery.ToDomainParams()
	courses, total, err := h.usecase.List(c.Request().Context(), pag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Course]{
		Items:  courses,
		Total:  total,
		Limit:  pag.Limit,
		Offset: pag.Offset,
	})
}
