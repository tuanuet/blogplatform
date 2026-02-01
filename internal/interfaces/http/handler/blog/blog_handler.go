package blog

import (
	"net/http"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type blogHandler struct {
	blogUseCase usecase.BlogUseCase
}

func NewBlogHandler(blogUseCase usecase.BlogUseCase) BlogHandler {
	return &blogHandler{
		blogUseCase: blogUseCase,
	}
}

// Create godoc
// @Summary Create a new blog
// @Description Create a new blog post (draft by default)
// @Tags Blogs
// @Accept json
// @Produce json
// @Param request body dto.CreateBlogRequest true "Blog data"
// @Success 201 {object} dto.BlogResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs [post]
func (h *blogHandler) Create(c *gin.Context) {
	// Get author ID from context (set by auth middleware)
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req dto.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	blog, err := h.blogUseCase.Create(c.Request.Context(), authorID.(uuid.UUID), &req)
	if err != nil {
		if err == usecase.ErrSlugAlreadyExists {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, blog)
}

// GetByID godoc
// @Summary Get blog by ID
// @Description Get a blog post by its ID
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} dto.BlogResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/blogs/{id} [get]
func (h *blogHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	// Get viewer ID if authenticated
	var viewerID *uuid.UUID
	if userID, exists := c.Get("userID"); exists {
		uid := userID.(uuid.UUID)
		viewerID = &uid
	}

	blog, err := h.blogUseCase.GetByID(c.Request.Context(), id, viewerID)
	if err != nil {
		if err == usecase.ErrBlogNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == usecase.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, blog)
}

// List godoc
// @Summary List blogs
// @Description List blogs with optional filters
// @Tags Blogs
// @Accept json
// @Produce json
// @Param authorId query string false "Filter by author ID"
// @Param categoryId query string false "Filter by category ID"
// @Param status query string false "Filter by status (draft, published)"
// @Param visibility query string false "Filter by visibility (public, subscribers_only)"
// @Param search query string false "Search in title or content"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/blogs [get]
func (h *blogHandler) List(c *gin.Context) {
	var params dto.BlogFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	// Get viewer ID if authenticated
	var viewerID *uuid.UUID
	if userID, exists := c.Get("userID"); exists {
		uid := userID.(uuid.UUID)
		viewerID = &uid
	}

	result, err := h.blogUseCase.List(c.Request.Context(), &params, viewerID)
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
// @Summary Update a blog
// @Description Update an existing blog (author only)
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param request body dto.UpdateBlogRequest true "Blog data"
// @Success 200 {object} dto.BlogResponse
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id} [put]
func (h *blogHandler) Update(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	var req dto.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	blog, err := h.blogUseCase.Update(c.Request.Context(), id, authorID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case usecase.ErrBlogNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrBlogAccessDenied:
			response.Forbidden(c, err.Error())
		case usecase.ErrSlugAlreadyExists:
			response.Conflict(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, blog)
}

// Delete godoc
// @Summary Delete a blog
// @Description Soft delete a blog (author only)
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 204
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id} [delete]
func (h *blogHandler) Delete(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	if err := h.blogUseCase.Delete(c.Request.Context(), id, authorID.(uuid.UUID)); err != nil {
		switch err {
		case usecase.ErrBlogNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrBlogAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// Publish godoc
// @Summary Publish a blog
// @Description Publish a draft blog with visibility setting
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param request body dto.PublishBlogRequest true "Visibility setting"
// @Success 200 {object} dto.BlogResponse
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/publish [post]
func (h *blogHandler) Publish(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	var req dto.PublishBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	blog, err := h.blogUseCase.Publish(c.Request.Context(), id, authorID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case usecase.ErrBlogNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrBlogAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, blog)
}

// Unpublish godoc
// @Summary Unpublish a blog
// @Description Revert a published blog back to draft
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} dto.BlogResponse
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/unpublish [post]
func (h *blogHandler) Unpublish(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	blog, err := h.blogUseCase.Unpublish(c.Request.Context(), id, authorID.(uuid.UUID))
	if err != nil {
		switch err {
		case usecase.ErrBlogNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrBlogAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, blog)
}

// React godoc
// @Summary React to a blog
// @Description Upvote, downvote, or remove reaction from a blog
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param request body dto.ReactionRequest true "Reaction"
// @Success 200 {object} dto.ReactionResponse
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/reaction [post]
func (h *blogHandler) React(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	var req dto.ReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.blogUseCase.React(c.Request.Context(), id, userID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case usecase.ErrBlogNotFound:
			response.NotFound(c, err.Error())
		case usecase.ErrBlogAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, res)
}
