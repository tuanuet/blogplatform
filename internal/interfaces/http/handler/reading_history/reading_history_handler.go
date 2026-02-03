package reading_history

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"net/http"
	"strconv"

	readingHistoryUsecase "github.com/aiagent/internal/application/usecase/reading_history"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReadingHistoryHandler interface {
	MarkAsRead(c *gin.Context)
	GetHistory(c *gin.Context)
}

type readingHistoryHandler struct {
	historyUseCase readingHistoryUsecase.ReadingHistoryUseCase
}

func NewReadingHistoryHandler(historyUseCase readingHistoryUsecase.ReadingHistoryUseCase) ReadingHistoryHandler {
	return &readingHistoryHandler{
		historyUseCase: historyUseCase,
	}
}

// MarkAsRead godoc
// @Summary Record a blog view
// @Description Record that a user viewed a blog post (Upsert)
// @Tags ReadingHistory
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Success 204
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Security Bearer
// @Router /api/v1/blogs/{id}/read [post]
func (h *readingHistoryHandler) MarkAsRead(c *gin.Context) {
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

	err = h.historyUseCase.MarkAsRead(c.Request.Context(), userID.(uuid.UUID), blogID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHistory godoc
// @Summary Get reading history
// @Description Get the list of recently viewed blogs
// @Tags ReadingHistory
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of records (default 20)"
// @Success 200 {object} dto.ReadingHistoryListResponse
// @Failure 401 {object} response.Response
// @Security Bearer
// @Router /api/v1/me/history [get]
func (h *readingHistoryHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	history, err := h.historyUseCase.GetHistory(c.Request.Context(), userID.(uuid.UUID), limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, history)
}
