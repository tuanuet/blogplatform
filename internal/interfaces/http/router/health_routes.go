package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(engine *gin.Engine, v1 *gin.RouterGroup, p Params) {
	// Health check routes
	engine.GET("/ping", p.HealthHandler.Ping)

	// API v1 routes
	// Health
	v1.GET("/health", p.HealthHandler.Check)
}
