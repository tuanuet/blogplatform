package bookmark

import (
	"net/http"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bookmarkHandler struct {
	bookmarkUseCase usecase.BookmarkUseCase
}

func NewBookmarkHandler(bookmarkUseCase usecase.BookmarkUseCase) BookmarkHandler {
	return &bookmarkHandler{
		bookmarkUseCase: bookmarkUseCase,
	}
}

// Bookmark godoc
// @Summary Bookmark a blog
// @Description Add a blog to user's bookmarks
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 204
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/bookmark [post]
func (h *bookmarkHandler) Bookmark(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	blogID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	if err := h.bookmarkUseCase.BookmarkBlog(c.Request.Context(), userID.(uuid.UUID), blogID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// Unbookmark godoc
// @Summary Unbookmark a blog
// @Description Remove a blog from user's bookmarks
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 204
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/bookmark [delete]
func (h *bookmarkHandler) Unbookmark(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	blogID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	if err := h.bookmarkUseCase.UnbookmarkBlog(c.Request.Context(), userID.(uuid.UUID), blogID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary Get my bookmarks
// @Description Get list of bookmarked blogs
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Security Bearer
// @Router /api/v1/bookmarks [get]
func (h *bookmarkHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	var params dto.BlogFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	result, err := h.bookmarkUseCase.GetUserBookmarks(c.Request.Context(), userID.(uuid.UUID), &params)
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
