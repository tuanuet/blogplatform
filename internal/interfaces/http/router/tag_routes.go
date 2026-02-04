package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	tags := v1.Group("/tags")
	{
		tags.GET("", p.TagHandler.List)
		tags.GET("/popular", p.RecommendationHandler.GetPopularTags) // Popular tags
		tags.GET("/:id", p.TagHandler.GetByID)
		tags.POST("", sessionAuth, auth.RequireCreate("tags"), p.TagHandler.Create)       // Requires CREATE permission
		tags.PUT("/:id", sessionAuth, auth.RequireUpdate("tags"), p.TagHandler.Update)    // Requires UPDATE permission
		tags.DELETE("/:id", sessionAuth, auth.RequireDelete("tags"), p.TagHandler.Delete) // Requires DELETE permission
	}
}
