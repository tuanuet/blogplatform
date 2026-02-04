package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterSeriesRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	seriesGroup := v1.Group("/series")
	{
		seriesGroup.GET("", p.SeriesHandler.List)
		seriesGroup.GET("/:id", p.SeriesHandler.GetByID)
		seriesGroup.GET("/slug/:slug", p.SeriesHandler.GetBySlug)
		seriesGroup.POST("", sessionAuth, auth.RequireCreate("series"), p.SeriesHandler.Create)
		seriesGroup.PUT("/:id", sessionAuth, auth.RequireUpdate("series"), p.SeriesHandler.Update)
		seriesGroup.DELETE("/:id", sessionAuth, auth.RequireDelete("series"), p.SeriesHandler.Delete)
		seriesGroup.POST("/:id/blogs", sessionAuth, auth.RequireUpdate("series"), p.SeriesHandler.AddBlog)
		seriesGroup.DELETE("/:id/blogs/:blogId", sessionAuth, auth.RequireUpdate("series"), p.SeriesHandler.RemoveBlog)
	}
}
