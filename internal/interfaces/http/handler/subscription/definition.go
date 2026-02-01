package subscription

import "github.com/gin-gonic/gin"

type SubscriptionHandler interface {
	Subscribe(c *gin.Context)
	Unsubscribe(c *gin.Context)
	GetMySubscriptions(c *gin.Context)
	GetSubscribers(c *gin.Context)
	CountSubscribers(c *gin.Context)
	GetSubscriptionCounts(c *gin.Context)
	GetUserSubscriptions(c *gin.Context)
	CheckSubscriptionStatus(c *gin.Context)
}
