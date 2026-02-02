package fraud

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

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
