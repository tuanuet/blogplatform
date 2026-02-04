package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterSubscriptionRoutes(v1 *gin.RouterGroup, p Params, sessionAuth gin.HandlerFunc) {
	// Authors & Subscriptions
	authors := v1.Group("/authors")
	{
		authors.GET("/:authorId/subscribers", p.SubscriptionHandler.GetSubscribers)
		authors.GET("/:authorId/subscribers/count", p.SubscriptionHandler.CountSubscribers)
		authors.POST("/:authorId/subscribe", sessionAuth, p.SubscriptionHandler.Subscribe)
		authors.POST("/:authorId/unsubscribe", sessionAuth, p.SubscriptionHandler.Unsubscribe)
	}

	// My subscriptions
	v1.GET("/subscriptions", sessionAuth, p.SubscriptionHandler.GetMySubscriptions)

	// Unified Subscription/Follow API (users can follow/subscribe to each other)
	v1.GET("/users/:userId/followers", p.SubscriptionHandler.GetSubscribers)
	v1.GET("/users/:userId/following", p.SubscriptionHandler.GetUserSubscriptions)
	v1.GET("/users/:userId/follow-counts", p.SubscriptionHandler.GetSubscriptionCounts)
	v1.POST("/users/:userId/follow", sessionAuth, p.SubscriptionHandler.Subscribe)
	v1.DELETE("/users/:userId/follow", sessionAuth, p.SubscriptionHandler.Unsubscribe)
}
