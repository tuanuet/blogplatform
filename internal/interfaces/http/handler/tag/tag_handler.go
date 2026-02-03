package tag

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	tagUsecase "github.com/aiagent/internal/application/usecase/tag"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tagHandler struct {
	tagUseCase tagUsecase.TagUseCase
}

func NewTagHandler(tagUseCase tagUsecase.TagUseCase) TagHandler {
	return &tagHandler{
		tagUseCase: tagUseCase,
	}
}

// Create godoc
// @Summary Create a tag
// @Description Create a new blog tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param request body dto.CreateTagRequest true "Tag data"
// @Success 201 {object} dto.TagResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Security Bearer
// @Router /api/v1/tags [post]
func (h *tagHandler) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.tagUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		if err == tagUsecase.ErrTagSlugExists {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, tag)
}

// GetByID godoc
// @Summary Get tag by ID
// @Description Get a tag by its ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} dto.TagResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/tags/{id} [get]
func (h *tagHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid tag ID")
		return
	}

	tag, err := h.tagUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == tagUsecase.ErrTagNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, tag)
}

// List godoc
// @Summary List tags
// @Description Get all tags
// @Tags Tags
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(50)
// @Success 200 {object} response.Response
// @Router /api/v1/tags [get]
func (h *tagHandler) List(c *gin.Context) {
	page := 1
	pageSize := 50

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	result, err := h.tagUseCase.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, result.Data, &response.Meta{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}

// Update godoc
// @Summary Update a tag
// @Description Update an existing tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param request body dto.UpdateTagRequest true "Tag data"
// @Success 200 {object} dto.TagResponse
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Security Bearer
// @Router /api/v1/tags/{id} [put]
func (h *tagHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid tag ID")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag, err := h.tagUseCase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case tagUsecase.ErrTagNotFound:
			response.NotFound(c, err.Error())
		case tagUsecase.ErrTagSlugExists:
			response.Conflict(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, tag)
}

// Delete godoc
// @Summary Delete a tag
// @Description Delete a tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Success 204
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/tags/{id} [delete]
func (h *tagHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid tag ID")
		return
	}

	if err := h.tagUseCase.Delete(c.Request.Context(), id); err != nil {
		if err == tagUsecase.ErrTagNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
