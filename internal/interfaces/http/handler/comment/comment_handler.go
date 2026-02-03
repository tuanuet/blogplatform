package comment

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	commentUsecase "github.com/aiagent/internal/application/usecase/comment"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type commentHandler struct {
	commentUseCase commentUsecase.CommentUseCase
}

func NewCommentHandler(commentUseCase commentUsecase.CommentUseCase) CommentHandler {
	return &commentHandler{
		commentUseCase: commentUseCase,
	}
}

// Create godoc
// @Summary Create a comment
// @Description Create a comment on a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param blogId path string true "Blog ID"
// @Param request body dto.CreateCommentRequest true "Comment data"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{blogId}/comments [post]
func (h *commentHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	blogID, err := uuid.Parse(c.Param("blogId"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := h.commentUseCase.Create(c.Request.Context(), userID.(uuid.UUID), blogID, &req)
	if err != nil {
		// Note: We might need to handle BlogNotFound if enforced by FK or service check.
		// UseCase currently propagates repo errors.
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, comment)
}

// GetByBlogID godoc
// @Summary Get comments for a blog
// @Description Get all comments for a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param blogId path string true "Blog ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/blogs/{blogId}/comments [get]
func (h *commentHandler) GetByBlogID(c *gin.Context) {
	blogID, err := uuid.Parse(c.Param("blogId"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

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

	result, err := h.commentUseCase.GetByBlogID(c.Request.Context(), blogID, page, pageSize)
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
// @Summary Update a comment
// @Description Update a comment (owner only)
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Comment data"
// @Success 200 {object} dto.CommentResponse
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/comments/{id} [put]
func (h *commentHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid comment ID")
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := h.commentUseCase.Update(c.Request.Context(), id, userID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case commentUsecase.ErrCommentNotFound:
			response.NotFound(c, err.Error())
		case commentUsecase.ErrCommentAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, comment)
}

// Delete godoc
// @Summary Delete a comment
// @Description Delete a comment (owner only)
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 204
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/comments/{id} [delete]
func (h *commentHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid comment ID")
		return
	}

	if err := h.commentUseCase.Delete(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		switch err {
		case commentUsecase.ErrCommentNotFound:
			response.NotFound(c, err.Error())
		case commentUsecase.ErrCommentAccessDenied:
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
