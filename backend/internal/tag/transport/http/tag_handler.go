// internal/tag/transport/http/tag_handler.go
package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/shared/dto"
	"github.com/kostinp/edu-platform-backend/internal/shared/pagination"
	"github.com/kostinp/edu-platform-backend/internal/tag/entity"
	"github.com/kostinp/edu-platform-backend/internal/tag/usecase"
	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	TagUsecase           usecase.TagUsecase
	TagAssignmentUsecase usecase.TagAssignmentUsecase
}

func NewTagHandler(tu usecase.TagUsecase, tau usecase.TagAssignmentUsecase) *TagHandler {
	return &TagHandler{
		TagUsecase:           tu,
		TagAssignmentUsecase: tau,
	}
}

// CreateTag
// @Summary Создать тег
// @Tags Tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param tag body entity.Tag true "Тег"
// @Success 201 {object} entity.Tag
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags [post]
func (h *TagHandler) CreateTag(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}
	tag := new(entity.Tag)
	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err = h.TagUsecase.CreateTag(c.Request().Context(), tag, authorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create tag"})
	}
	return c.JSON(http.StatusCreated, tag)
}

// ListTags
// @Summary Получить список тегов
// @Tags Tag
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Количество" default(50)
// @Param offset query int false "Смещение" default(0)
// @Param sort_by query string false "Поле сортировки (name, created_at)" default(created_at)
// @Param order query string false "Порядок сортировки (asc, desc)" default(desc)
// @Success 200 {object} dto.PaginatedResponse[*entity.Tag]
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags [get]
func (h *TagHandler) ListTags(c echo.Context) error {
	pagQuery := pagination.ParsePaginationParams(c)
	pag := pagQuery.ToDomainParams()
	tags, total, err := h.TagUsecase.ListTags(c.Request().Context(), pag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list tags",
		})
	}
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Tag]{
		Items:  tags,
		Total:  total,
		Limit:  pag.Limit,
		Offset: pag.Offset,
	})
}

// GetTagByID
// @Summary Получить тег по ID
// @Tags Tag
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID тега"
// @Success 200 {object} entity.Tag
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags/{id} [get]
func (h *TagHandler) GetTagByID(c echo.Context) error {
	idStr := c.Param("id")
	tagID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag id"})
	}
	tag, err := h.TagUsecase.GetTagByID(c.Request().Context(), tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get tag"})
	}
	if tag == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "tag not found"})
	}
	return c.JSON(http.StatusOK, tag)
}

// DeleteTag
// @Summary Удалить тег по ID
// @Tags Tag
// @Security BearerAuth
// @Param id path string true "ID тега"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c echo.Context) error {
	idStr := c.Param("id")
	tagID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag id"})
	}
	err = h.TagUsecase.DeleteTag(c.Request().Context(), tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete tag"})
	}
	return c.NoContent(http.StatusNoContent)
}

// AssignTag
// @Summary Назначить тег сущности
// @Tags TagAssignment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param assignment body entity.TagAssignment true "Назначение тега"
// @Success 201 {object} entity.TagAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tag-assignments [post]
func (h *TagHandler) AssignTag(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	authorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}
	ta := new(entity.TagAssignment)
	if err := c.Bind(ta); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	err = h.TagAssignmentUsecase.AssignTag(c.Request().Context(), ta, authorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to assign tag"})
	}
	return c.JSON(http.StatusCreated, ta)
}

// RemoveAssignment
// @Summary Удалить назначение тега по ID
// @Tags TagAssignment
// @Security BearerAuth
// @Param id path string true "ID назначения тега"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tag-assignments/{id} [delete]
func (h *TagHandler) RemoveAssignment(c echo.Context) error {
	idStr := c.Param("id")
	assignmentID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid assignment id"})
	}
	err = h.TagAssignmentUsecase.RemoveAssignment(c.Request().Context(), assignmentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to remove assignment"})
	}
	return c.NoContent(http.StatusNoContent)
}

// ListAssignmentsByEntity
// @Summary Получить назначения тегов для сущности
// @Tags TagAssignment
// @Security BearerAuth
// @Produce json
// @Param entity_id query string true "ID сущности"
// @Param entity_type query string true "Тип сущности"
// @Success 200 {array} entity.TagAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tag-assignments/by-entity [get]
func (h *TagHandler) ListAssignmentsByEntity(c echo.Context) error {
	entityIDStr := c.QueryParam("entity_id")
	entityType := c.QueryParam("entity_type")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid entity id"})
	}
	if entityType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "entity_type required"})
	}
	assignments, err := h.TagAssignmentUsecase.ListAssignmentsByEntity(c.Request().Context(), entityID, entityType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list assignments"})
	}
	return c.JSON(http.StatusOK, assignments)
}

// ListAssignmentsByTag
// @Summary Получить назначения по тегу
// @Tags TagAssignment
// @Security BearerAuth
// @Produce json
// @Param tag_id query string true "ID тега"
// @Success 200 {array} entity.TagAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tag-assignments/by-tag [get]
func (h *TagHandler) ListAssignmentsByTag(c echo.Context) error {
	tagIDStr := c.QueryParam("tag_id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag id"})
	}
	assignments, err := h.TagAssignmentUsecase.ListAssignmentsByTag(c.Request().Context(), tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list assignments"})
	}
	return c.JSON(http.StatusOK, assignments)
}
