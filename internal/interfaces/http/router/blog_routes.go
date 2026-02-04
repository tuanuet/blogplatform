package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	blogs := v1.Group("/blogs")
	{
		blogs.GET("", p.BlogHandler.List)
		blogs.GET("/feed", sessionAuth, p.RecommendationHandler.GetPersonalizedFeed) // Personalized feed
		blogs.GET("/:id", p.BlogHandler.GetByID)
		blogs.GET("/:id/related", p.RecommendationHandler.GetRelatedBlogs)                              // Related blogs
		blogs.POST("", sessionAuth, auth.RequireCreate("blogs"), p.BlogHandler.Create)                  // Requires CREATE permission
		blogs.PUT("/:id", sessionAuth, auth.RequireUpdate("blogs"), p.BlogHandler.Update)               // Requires UPDATE permission
		blogs.DELETE("/:id", sessionAuth, auth.RequireDelete("blogs"), p.BlogHandler.Delete)            // Requires DELETE permission
		blogs.POST("/:id/publish", sessionAuth, auth.RequireUpdate("blogs"), p.BlogHandler.Publish)     // Requires UPDATE permission
		blogs.POST("/:id/unpublish", sessionAuth, auth.RequireUpdate("blogs"), p.BlogHandler.Unpublish) // Requires UPDATE permission
		blogs.POST("/:id/reaction", sessionAuth, p.BlogHandler.React)                                   // Authenticated users
		blogs.POST("/:id/read", sessionAuth, p.ReadingHistoryHandler.MarkAsRead)                        // Authenticated users
		blogs.POST("/:id/bookmark", sessionAuth, p.BookmarkHandler.Bookmark)
		blogs.DELETE("/:id/bookmark", sessionAuth, p.BookmarkHandler.Unbookmark)

		// Blog comments
		blogs.GET("/:id/comments", p.CommentHandler.GetByBlogID)
		blogs.POST("/:id/comments", sessionAuth, auth.RequireCreate("comments"), p.CommentHandler.Create)
	}
}
