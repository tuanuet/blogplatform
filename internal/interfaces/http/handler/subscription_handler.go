package handler

import (
	"net/http"
	"strconv"

	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionHandler handles subscription-related HTTP requests
type SubscriptionHandler struct {
	subscriptionUseCase usecase.SubscriptionUseCase
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionUseCase usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionUseCase: subscriptionUseCase,
	}
}

// Subscribe godoc
// @Summary Subscribe to an author
// @Description Subscribe to receive access to subscriber-only content
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param authorId path string true "Author ID"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Security Bearer
// @Router /api/v1/authors/{authorId}/subscribe [post]
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	subscriberID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		response.BadRequest(c, "invalid author ID")
		return
	}

	subscription, err := h.subscriptionUseCase.Subscribe(c.Request.Context(), subscriberID.(uuid.UUID), authorID)
	if err != nil {
		switch err {
		case usecase.ErrCannotSubscribeToSelf:
			response.BadRequest(c, err.Error())
		case usecase.ErrAlreadySubscribed:
			response.Conflict(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, subscription)
}

// Unsubscribe godoc
// @Summary Unsubscribe from an author
// @Description Remove subscription from an author
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param authorId path string true "Author ID"
// @Success 204
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/authors/{authorId}/unsubscribe [post]
func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	subscriberID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		response.BadRequest(c, "invalid author ID")
		return
	}

	if err := h.subscriptionUseCase.Unsubscribe(c.Request.Context(), subscriberID.(uuid.UUID), authorID); err != nil {
		if err == usecase.ErrSubscriptionNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMySubscriptions godoc
// @Summary Get my subscriptions
// @Description Get list of authors the current user is subscribed to
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Security Bearer
// @Router /api/v1/subscriptions [get]
func (h *SubscriptionHandler) GetMySubscriptions(c *gin.Context) {
	subscriberID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
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

	result, err := h.subscriptionUseCase.GetSubscriptions(c.Request.Context(), subscriberID.(uuid.UUID), page, pageSize)
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

// GetSubscribers godoc
// @Summary Get author's subscribers
// @Description Get list of subscribers for an author
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param authorId path string true "Author ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/authors/{authorId}/subscribers [get]
func (h *SubscriptionHandler) GetSubscribers(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		response.BadRequest(c, "invalid author ID")
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

	result, err := h.subscriptionUseCase.GetSubscribers(c.Request.Context(), authorID, page, pageSize)
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

// CountSubscribers godoc
// @Summary Get subscriber count
// @Description Get the number of subscribers for an author
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param authorId path string true "Author ID"
// @Success 200 {object} dto.SubscriptionCountResponse
// @Router /api/v1/authors/{authorId}/subscribers/count [get]
func (h *SubscriptionHandler) CountSubscribers(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		response.BadRequest(c, "invalid author ID")
		return
	}

	count, err := h.subscriptionUseCase.CountSubscribers(c.Request.Context(), authorID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, count)
}
