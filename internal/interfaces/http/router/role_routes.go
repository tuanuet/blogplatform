package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	roles := v1.Group("/roles", sessionAuth)
	{
		roles.GET("", p.RoleHandler.List)
		roles.GET("/:id", p.RoleHandler.GetByID)
		roles.POST("", auth.RequireAdmin("roles"), p.RoleHandler.Create)                        // Admin only
		roles.PUT("/:id", auth.RequireAdmin("roles"), p.RoleHandler.Update)                     // Admin only
		roles.DELETE("/:id", auth.RequireAdmin("roles"), p.RoleHandler.Delete)                  // Admin only
		roles.POST("/:id/permissions", auth.RequireAdmin("roles"), p.RoleHandler.SetPermission) // Admin only
	}
}
