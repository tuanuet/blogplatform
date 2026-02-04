package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	comments := v1.Group("/comments", sessionAuth)
	{
		comments.PUT("/:id", auth.RequireUpdate("comments"), p.CommentHandler.Update)
		comments.DELETE("/:id", auth.RequireDelete("comments"), p.CommentHandler.Delete)
	}
}
