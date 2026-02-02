package recommendation

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type recommendationHandler struct {
	recUseCase usecase.RecommendationUseCase
}

func NewRecommendationHandler(recUseCase usecase.RecommendationUseCase) RecommendationHandler {
	return &recommendationHandler{
		recUseCase: recUseCase,
	}
}

// GetPopularTags godoc
// @Summary Get popular tags
// @Description Get top 10 most used tags
// @Tags Recommendations
// @Accept json
// @Produce json
// @Success 200 {array} dto.TagResponse
// @Failure 500 {object} response.Response
// @Router /api/v1/tags/popular [get]
func (h *recommendationHandler) GetPopularTags(c *gin.Context) {
	tags, err := h.recUseCase.GetPopularTags(c.Request.Context(), 10)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, tags)
}

// UpdateInterests godoc
// @Summary Update user interests
// @Description Update the logged-in user's interested tags
// @Tags Recommendations
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserInterestsRequest true "Tag IDs"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /api/v1/users/me/interests [post]
func (h *recommendationHandler) UpdateInterests(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req dto.UpdateUserInterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.recUseCase.UpdateInterests(c.Request.Context(), userID.(uuid.UUID), req.TagIDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "interests updated successfully"})
}

// GetPersonalizedFeed godoc
// @Summary Get personalized feed
// @Description Get blogs based on user interests (fallback to recent if no interests)
// @Tags Recommendations
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/feed [get]
func (h *recommendationHandler) GetPersonalizedFeed(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		// If not logged in, we can't personalize.
		// Requirement says: "Fallback: Show trending/recent posts if no matches or no interests set."
		// If not logged in, we should probably redirect to standard list or return error?
		// Spec says: "Personalized 'For You' Feed". Usually implies auth.
		// "Fallback: Show trending/recent posts if no matches or no interests set."
		// I'll return generic feed if not logged in (empty userID won't match anything, or I handle it in UseCase).
		// UseCase expects userID. I'll pass nil or handle here.
		// UseCase signature: GetPersonalizedFeed(ctx, userID, ...)
		// If I pass Nil UUID, it will fetch empty interests and return generic feed.
		// But passing Nil UUID is tricky if type is uuid.UUID (value type).
		// I'll assume this endpoint requires Auth. If they want public feed, they use /blogs.
		response.Unauthorized(c, "authentication required for personalized feed")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.recUseCase.GetPersonalizedFeed(c.Request.Context(), userID.(uuid.UUID), page, pageSize)
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

// GetRelatedBlogs godoc
// @Summary Get related blogs
// @Description Get related blogs based on tags
// @Tags Recommendations
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {array} dto.BlogListResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/blogs/{id}/related [get]
func (h *recommendationHandler) GetRelatedBlogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	blogs, err := h.recUseCase.GetRelatedBlogs(c.Request.Context(), id, 3)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, blogs)
}
