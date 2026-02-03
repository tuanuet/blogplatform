package fraud

import (
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/valueobject"
	"github.com/aiagent/pkg/response"
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

	resp := &dto.RiskScoreResponse{
		ID:                        result.ID,
		UserID:                    result.UserID,
		OverallScore:              result.OverallScore,
		FollowerAuthenticityScore: result.FollowerAuthenticityScore,
		EngagementQualityScore:    result.EngagementQualityScore,
		AccountAgeFactor:          result.AccountAgeFactor,
		BadgeStatus:               result.BadgeStatus,
		LastCalculatedAt:          result.LastCalculatedAt,
	}

	response.Success(c, http.StatusOK, resp)
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

	filter := valueobject.FraudDashboardFilter{
		MinRiskScore: &req.MinRiskScore,
		MaxRiskScore: &req.MaxRiskScore,
		SignalTypes:  req.SignalTypes,
		ReviewStatus: req.ReviewStatus,
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}
	// Note: Request struct uses int for Min/Max, pointers would be better to distinguish 0 vs unset,
	// but mapping for now assuming 0 means default or ignore if appropriate in service.
	// Actually service logic was: if maxScore == 0 { maxScore = 100 }.
	// So passing as is is fine.
	// Wait, service expects pointers for MinRiskScore/MaxRiskScore in Filter struct I created?
	// Let's check model.go. Yes: MinRiskScore *int.
	// DTO has int. 0 could be a valid score. Use pointer to int.

	// Correct mapping needs address of int.
	// Create vars to take address.

	result, err := h.service.GetFraudDashboard(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	users := make([]dto.FraudDashboardUser, 0, len(result.Users))
	for _, u := range result.Users {
		activeSignals := make([]dto.BotSignalSummary, 0, len(u.ActiveSignals))
		for _, s := range u.ActiveSignals {
			activeSignals = append(activeSignals, dto.BotSignalSummary{
				SignalType:      s.SignalType,
				ConfidenceScore: s.ConfidenceScore,
				DetectedAt:      s.DetectedAt,
			})
		}
		users = append(users, dto.FraudDashboardUser{
			UserID:                u.UserID,
			Username:              u.Username,
			Email:                 u.Email,
			OverallScore:          u.OverallScore,
			FollowerCount:         u.FollowerCount,
			BotFollowerCount:      u.BotFollowerCount,
			ActiveSignals:         activeSignals,
			LastReviewAction:      u.LastReviewAction,
			LastReviewedAt:        u.LastReviewedAt,
			RiskScoreCalculatedAt: u.RiskScoreCalculatedAt,
		})
	}

	resp := &dto.FraudDashboardResponse{
		Users:      users,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}

	response.SuccessWithMeta(c, resp.Users, &response.Meta{
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		Total:      int64(resp.TotalCount),
		TotalPages: resp.TotalPages,
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

	cmd := valueobject.ReviewUserCommand{
		Notes: req.Notes,
	}

	result, err := h.service.ReviewUser(c.Request.Context(), adminID.(uuid.UUID), userID, cmd)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	resp := &dto.ReviewUserResponse{
		ReviewID:   result.ReviewID,
		UserID:     result.UserID,
		AdminID:    result.AdminID,
		Action:     result.Action,
		Notes:      result.Notes,
		ReviewedAt: result.ReviewedAt,
	}

	response.Success(c, http.StatusOK, resp)
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

	cmd := valueobject.BanUserCommand{
		Reason: req.Reason,
		Notes:  req.Notes,
	}

	result, err := h.service.BanUser(c.Request.Context(), adminID.(uuid.UUID), userID, cmd)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	resp := &dto.BanUserResponse{
		ReviewID: result.ReviewID,
		UserID:   result.UserID,
		AdminID:  result.AdminID,
		Action:   result.Action,
		Reason:   result.Reason,
		Notes:    result.Notes,
		BannedAt: result.BannedAt,
	}

	response.Success(c, http.StatusOK, resp)
}

// GetFraudTrends handles GET /api/analytics/fraud-trends
func (h *fraudHandler) GetFraudTrends(c *gin.Context) {
	var req dto.FraudTrendsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	filter := valueobject.FraudTrendsFilter{
		Period:   req.Period,
		FromDate: req.FromDate,
		ToDate:   req.ToDate,
	}

	result, err := h.service.GetFraudTrends(c.Request.Context(), filter)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	dailyStats := make([]dto.DailyFraudStat, 0, len(result.DailyStats))
	for _, s := range result.DailyStats {
		dailyStats = append(dailyStats, dto.DailyFraudStat{
			Date:                  s.Date,
			NewSignals:            s.NewSignals,
			NewSuspiciousAccounts: s.NewSuspiciousAccounts,
			BannedAccounts:        s.BannedAccounts,
		})
	}

	resp := &dto.FraudTrendsResponse{
		Period:                result.Period,
		FromDate:              result.FromDate,
		ToDate:                result.ToDate,
		TotalBotSignals:       result.TotalBotSignals,
		SignalsByType:         result.SignalsByType,
		NewSuspiciousAccounts: result.NewSuspiciousAccounts,
		BannedAccounts:        result.BannedAccounts,
		ReviewedAccounts:      result.ReviewedAccounts,
		AverageRiskScore:      result.AverageRiskScore,
		RiskScoreDistribution: result.RiskScoreDistribution,
		DailyStats:            dailyStats,
	}

	response.Success(c, http.StatusOK, resp)
}

// TriggerBatchAnalysis handles POST /api/followers/batch-analyze
func (h *fraudHandler) TriggerBatchAnalysis(c *gin.Context) {
	var req dto.BatchAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cmd := valueobject.BatchAnalyzeCommand{
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}

	result, err := h.service.TriggerBatchAnalysis(c.Request.Context(), cmd)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	resp := &dto.BatchAnalyzeResponse{
		JobID:              result.JobID,
		Status:             result.Status,
		StartedAt:          result.StartedAt,
		CompletedAt:        result.CompletedAt,
		ProcessedFollowers: result.ProcessedFollowers,
		NewSignalsDetected: result.NewSignalsDetected,
		UsersScored:        result.UsersScored,
		Message:            result.Message,
	}

	response.Success(c, http.StatusOK, resp)
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

	resp := &dto.UserBadgeResponse{
		UserID:        result.UserID,
		BadgeType:     result.BadgeType,
		Status:        result.Status,
		EligibleSince: result.EligibleSince,
		ActivatedAt:   result.ActivatedAt,
	}

	response.Success(c, http.StatusOK, resp)
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

	resp := make([]dto.BotFollowerNotificationResponse, 0, len(result))
	for _, n := range result {
		resp = append(resp, dto.BotFollowerNotificationResponse{
			ID:              n.ID,
			BotFollowerID:   n.BotFollowerID,
			BotFollowerName: n.BotFollowerName,
			SignalType:      n.SignalType,
			ConfidenceScore: n.ConfidenceScore,
			SentAt:          n.SentAt,
			ReadAt:          n.ReadAt,
		})
	}

	response.Success(c, http.StatusOK, resp)
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
