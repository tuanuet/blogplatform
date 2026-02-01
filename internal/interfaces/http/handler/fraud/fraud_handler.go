package fraud

import (
	"net/http"

	"github.com/aiagent/boilerplate/internal/domain/service"
	"github.com/aiagent/boilerplate/internal/interfaces/http/dto"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// fraudHandler implements the FraudHandler interface
type fraudHandler struct {
	service service.FraudDetectionService
}

// NewFraudHandler creates a new fraud handler instance
func NewFraudHandler(service service.FraudDetectionService) FraudHandler {
	return &fraudHandler{
		service: service,
	}
}

// GetUserRiskScore handles GET /api/users/:id/risk-score
func (h *fraudHandler) GetUserRiskScore(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.service.GetUserRiskScore(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	if result == nil {
		response.NotFound(c, "Risk score not found for user")
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetFraudDashboard handles GET /api/admin/fraud-dashboard
func (h *fraudHandler) GetFraudDashboard(c *gin.Context) {
	var req dto.FraudDashboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	result, err := h.service.GetFraudDashboard(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, result.Users, &response.Meta{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      int64(result.TotalCount),
		TotalPages: result.TotalPages,
	})
}

// ReviewUser handles POST /api/admin/users/:id/review
func (h *fraudHandler) ReviewUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.ReviewUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get admin ID from context (set by auth middleware)
	adminID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Admin authentication required")
		return
	}

	result, err := h.service.ReviewUser(c.Request.Context(), adminID.(uuid.UUID), userID, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// BanUser handles POST /api/admin/users/:id/ban
func (h *fraudHandler) BanUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get admin ID from context (set by auth middleware)
	adminID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Admin authentication required")
		return
	}

	result, err := h.service.BanUser(c.Request.Context(), adminID.(uuid.UUID), userID, req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetFraudTrends handles GET /api/analytics/fraud-trends
func (h *fraudHandler) GetFraudTrends(c *gin.Context) {
	var req dto.FraudTrendsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.GetFraudTrends(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// TriggerBatchAnalysis handles POST /api/followers/batch-analyze
func (h *fraudHandler) TriggerBatchAnalysis(c *gin.Context) {
	var req dto.BatchAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.TriggerBatchAnalysis(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetUserBadgeStatus handles GET /api/users/:id/badge
func (h *fraudHandler) GetUserBadgeStatus(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	result, err := h.service.GetUserBadgeStatus(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	if result == nil {
		response.NotFound(c, "Badge status not found for user")
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetUserBotNotifications handles GET /api/users/:id/bot-notifications
func (h *fraudHandler) GetUserBotNotifications(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	unreadOnly := c.Query("unread_only") == "true"

	result, err := h.service.GetUserBotNotifications(c.Request.Context(), userID, unreadOnly)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

// MarkNotificationAsRead handles POST /api/notifications/:id/read
func (h *fraudHandler) MarkNotificationAsRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid notification ID")
		return
	}

	// This would typically call a service method to mark as read
	// For now, return success
	response.Success(c, http.StatusOK, gin.H{
		"notification_id": notificationID,
		"status":          "marked_as_read",
	})
}
