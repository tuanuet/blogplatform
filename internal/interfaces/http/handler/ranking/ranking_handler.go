package ranking

import (
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/ranking"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type rankingHandler struct {
	rankingUseCase ranking.RankingUseCase
}

func NewRankingHandler(rankingUseCase ranking.RankingUseCase) RankingHandler {
	return &rankingHandler{
		rankingUseCase: rankingUseCase,
	}
}

// GetTrending godoc
// @Summary Get trending users
// @Description Get users ranked by velocity score (follower growth + blog activity)
// @Tags Rankings
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/rankings/trending [get]
func (h *rankingHandler) GetTrending(c *gin.Context) {
	var params dto.RankingFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	result, err := h.rankingUseCase.GetTrendingUsers(c.Request.Context(), &params)
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

// GetTop godoc
// @Summary Get top users
// @Description Get top users by total followers
// @Tags Rankings
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/rankings/top [get]
func (h *rankingHandler) GetTop(c *gin.Context) {
	var params dto.RankingFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	result, err := h.rankingUseCase.GetTopUsers(c.Request.Context(), &params)
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

// GetUserRanking godoc
// @Summary Get user ranking
// @Description Get detailed ranking information for a specific user
// @Tags Rankings
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dto.UserRankingDetailResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/rankings/users/{userId} [get]
func (h *rankingHandler) GetUserRanking(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	ranking, err := h.rankingUseCase.GetUserRanking(c.Request.Context(), userID)
	if err != nil {
		if err == domainService.ErrUserNotFound {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, ranking)
}

// RecalculateScores godoc
// @Summary Recalculate all scores
// @Description Trigger a recalculation of all velocity scores (admin only)
// @Tags Rankings
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Security Bearer
// @Router /api/v1/rankings/recalculate [post]
func (h *rankingHandler) RecalculateScores(c *gin.Context) {
	// TODO: Add admin authorization check

	if err := h.rankingUseCase.RecalculateAllScores(c.Request.Context()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "rankings recalculated successfully"})
}
