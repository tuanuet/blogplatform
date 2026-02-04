package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	// Admin Dashboard
	v1.GET("/admin/dashboard/stats", sessionAuth, auth.RequireAdmin("analytics"), p.AdminHandler.GetDashboardStats)
}
