package category

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type categoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) CategoryHandler {
	return &categoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// Create godoc
// @Summary Create a category
// @Description Create a new blog category
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Security Bearer
// @Router /api/v1/categories [post]
func (h *categoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		if err == usecase.ErrCategorySlugExists {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, category)
}

// GetByID godoc
// @Summary Get category by ID
// @Description Get a category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/categories/{id} [get]
func (h *categoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid category ID")
		return
	}

	category, err := h.categoryUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == usecase.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, category)
}

// List godoc
// @Summary List categories
// @Description Get all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/categories [get]
func (h *categoryHandler) List(c *gin.Context) {
	page := 1
	pageSize := 20

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

	result, err := h.categoryUseCase.List(c.Request.Context(), page, pageSize)
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
// @Summary Update a category
// @Description Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} dto.CategoryResponse
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Security Bearer
// @Router /api/v1/categories/{id} [put]
func (h *categoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid category ID")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryUseCase.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case usecase.ErrCategoryNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrCategorySlugExists:
			response.Conflict(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, category)
}

// Delete godoc
// @Summary Delete a category
// @Description Delete a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 204
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/categories/{id} [delete]
func (h *categoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid category ID")
		return
	}

	if err := h.categoryUseCase.Delete(c.Request.Context(), id); err != nil {
		if err == usecase.ErrCategoryNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
