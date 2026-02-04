package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFraudRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	// User risk score and badge
	v1.GET("/users/:id/risk-score", p.FraudHandler.GetUserRiskScore)
	v1.GET("/users/:id/badge", p.FraudHandler.GetUserBadgeStatus)
	v1.GET("/users/:id/bot-notifications", p.FraudHandler.GetUserBotNotifications)

	// Admin fraud dashboard
	admin := v1.Group("/admin", sessionAuth, auth.RequireAdmin("fraud"))
	{
		admin.GET("/fraud-dashboard", p.FraudHandler.GetFraudDashboard)
		admin.POST("/users/:id/review", p.FraudHandler.ReviewUser)
		admin.POST("/users/:id/ban", p.FraudHandler.BanUser)
	}

	// Analytics
	v1.GET("/analytics/fraud-trends", sessionAuth, auth.RequireAdmin("analytics"), p.FraudHandler.GetFraudTrends)

	// Batch operations
	v1.POST("/followers/batch-analyze", sessionAuth, auth.RequireAdmin("fraud"), p.FraudHandler.TriggerBatchAnalysis)

	// Notifications
	v1.POST("/notifications/:id/read", sessionAuth, p.FraudHandler.MarkNotificationAsRead)
}
