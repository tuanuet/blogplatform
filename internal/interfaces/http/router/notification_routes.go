package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterNotificationRoutes registers the notification routes
func RegisterNotificationRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	notificationGroup := v1.Group("/notifications")
	notificationGroup.Use(sessionAuth) // User must be logged in
	{
		notificationGroup.GET("", p.NotificationHandler.List)
		notificationGroup.GET("/unread-count", p.NotificationHandler.GetUnreadCount)
		notificationGroup.POST("/:id/read", p.NotificationHandler.MarkAsRead)
		notificationGroup.POST("/read-all", p.NotificationHandler.MarkAllAsRead)
		notificationGroup.GET("/preferences", p.NotificationHandler.GetPreferences)
		notificationGroup.PUT("/preferences", p.NotificationHandler.UpdatePreferences)
		notificationGroup.POST("/device-token", p.NotificationHandler.RegisterDeviceToken)
	}
}
