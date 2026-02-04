package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	// Permissions
	v1.GET("/permissions", sessionAuth, p.RoleHandler.GetMyPermission)

	// Users
	users := v1.Group("/users")
	{
		users.POST("/me/interests", sessionAuth, p.RecommendationHandler.UpdateInterests) // Update interests
		users.GET("/:id/profile", p.ProfileHandler.GetPublicProfile)
		users.GET("/:id/roles", sessionAuth, p.RoleHandler.GetUserRoles)
		users.POST("/:id/roles", sessionAuth, auth.RequireAdmin("users"), p.RoleHandler.AssignRole)           // Admin only
		users.DELETE("/:id/roles/:roleId", sessionAuth, auth.RequireAdmin("users"), p.RoleHandler.RemoveRole) // Admin only
	}
}
