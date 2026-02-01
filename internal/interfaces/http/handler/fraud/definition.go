package fraud

import "github.com/gin-gonic/gin"

type FraudHandler interface {
	GetUserRiskScore(c *gin.Context)
	GetFraudDashboard(c *gin.Context)
	ReviewUser(c *gin.Context)
	BanUser(c *gin.Context)
	GetFraudTrends(c *gin.Context)
	TriggerBatchAnalysis(c *gin.Context)
	GetUserBadgeStatus(c *gin.Context)
	GetUserBotNotifications(c *gin.Context)
	MarkNotificationAsRead(c *gin.Context)
}
