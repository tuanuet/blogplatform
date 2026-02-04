package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	categories := v1.Group("/categories")
	{
		categories.GET("", p.CategoryHandler.List)
		categories.GET("/:id", p.CategoryHandler.GetByID)
		categories.POST("", sessionAuth, auth.RequireCreate("categories"), p.CategoryHandler.Create)       // Requires CREATE permission
		categories.PUT("/:id", sessionAuth, auth.RequireUpdate("categories"), p.CategoryHandler.Update)    // Requires UPDATE permission
		categories.DELETE("/:id", sessionAuth, auth.RequireDelete("categories"), p.CategoryHandler.Delete) // Requires DELETE permission
	}
}
