package router

import (
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/blog"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/bookmark"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/category"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/comment"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/fraud"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/health"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/profile"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/ranking"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/reading_history"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/recommendation"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/role"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/series"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/subscription"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/tag"
	"github.com/aiagent/boilerplate/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"

	// Import generated swagger docs
	_ "github.com/aiagent/boilerplate/docs"
)

// Params groups all handler dependencies for the router (FX dependency injection)
type Params struct {
	fx.In

	HealthHandler         health.HealthHandler
	BlogHandler           blog.BlogHandler
	BookmarkHandler       bookmark.BookmarkHandler
	CategoryHandler       category.CategoryHandler
	TagHandler            tag.TagHandler
	CommentHandler        comment.CommentHandler
	SubscriptionHandler   subscription.SubscriptionHandler
	ProfileHandler        profile.ProfileHandler
	RoleHandler           role.RoleHandler
	SeriesHandler         series.SeriesHandler
	RankingHandler        ranking.RankingHandler
	ReadingHistoryHandler reading_history.ReadingHistoryHandler
	RecommendationHandler recommendation.RecommendationHandler
	FraudHandler          fraud.FraudHandler
	RoleUseCase           usecase.RoleUseCase // For authorization middleware
	Config                *config.Config
}

// New creates and configures a new Gin engine with all routes
func New(p Params) *gin.Engine {
	gin.SetMode(p.Config.Server.Mode)
	engine := gin.New()

	// Global middleware
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logging())

	// OpenTelemetry Middleware
	if p.Config.Telemetry.Enabled {
		engine.Use(otelgin.Middleware(p.Config.Telemetry.ServiceName))
	}

	engine.Use(middleware.CORS())

	// Authorization middleware (for protected routes)
	auth := middleware.NewAuthorization(p.RoleUseCase)

	// Serve static files for avatar uploads
	engine.Static("/uploads", "./uploads")

	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check routes
	engine.GET("/ping", p.HealthHandler.Ping)

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Health
		v1.GET("/health", p.HealthHandler.Check)

		// Profile (authenticated user)
		v1.GET("/profile", p.ProfileHandler.GetMyProfile)
		v1.PUT("/profile", p.ProfileHandler.UpdateMyProfile)
		v1.POST("/profile/avatar", p.ProfileHandler.UploadAvatar)

		// Permissions
		v1.GET("/permissions", p.RoleHandler.GetMyPermission)

		// Users
		users := v1.Group("/users")
		{
			users.POST("/me/interests", p.RecommendationHandler.UpdateInterests) // Update interests
			users.GET("/:id/profile", p.ProfileHandler.GetPublicProfile)
			users.GET("/:id/roles", p.RoleHandler.GetUserRoles)
			users.POST("/:id/roles", auth.RequireAdmin("users"), p.RoleHandler.AssignRole)           // Admin only
			users.DELETE("/:id/roles/:roleId", auth.RequireAdmin("users"), p.RoleHandler.RemoveRole) // Admin only
		}

		// Roles (admin only for write operations)
		roles := v1.Group("/roles")
		{
			roles.GET("", p.RoleHandler.List)
			roles.GET("/:id", p.RoleHandler.GetByID)
			roles.POST("", auth.RequireAdmin("roles"), p.RoleHandler.Create)                        // Admin only
			roles.PUT("/:id", auth.RequireAdmin("roles"), p.RoleHandler.Update)                     // Admin only
			roles.DELETE("/:id", auth.RequireAdmin("roles"), p.RoleHandler.Delete)                  // Admin only
			roles.POST("/:id/permissions", auth.RequireAdmin("roles"), p.RoleHandler.SetPermission) // Admin only
		}

		// Blogs
		blogs := v1.Group("/blogs")
		{
			blogs.GET("", p.BlogHandler.List)
			blogs.GET("/feed", p.RecommendationHandler.GetPersonalizedFeed) // Personalized feed
			blogs.GET("/:id", p.BlogHandler.GetByID)
			blogs.GET("/:id/related", p.RecommendationHandler.GetRelatedBlogs)                 // Related blogs
			blogs.POST("", auth.RequireCreate("blogs"), p.BlogHandler.Create)                  // Requires CREATE permission
			blogs.PUT("/:id", auth.RequireUpdate("blogs"), p.BlogHandler.Update)               // Requires UPDATE permission
			blogs.DELETE("/:id", auth.RequireDelete("blogs"), p.BlogHandler.Delete)            // Requires DELETE permission
			blogs.POST("/:id/publish", auth.RequireUpdate("blogs"), p.BlogHandler.Publish)     // Requires UPDATE permission
			blogs.POST("/:id/unpublish", auth.RequireUpdate("blogs"), p.BlogHandler.Unpublish) // Requires UPDATE permission
			blogs.POST("/:id/reaction", p.BlogHandler.React)                                   // Authenticated users
			blogs.POST("/:id/read", p.ReadingHistoryHandler.MarkAsRead)                        // Authenticated users
			blogs.POST("/:id/bookmark", p.BookmarkHandler.Bookmark)
			blogs.DELETE("/:id/bookmark", p.BookmarkHandler.Unbookmark)

			// Blog comments
			blogs.GET("/:id/comments", p.CommentHandler.GetByBlogID)
			blogs.POST("/:id/comments", auth.RequireCreate("comments"), p.CommentHandler.Create)
		}

		// Series
		seriesGroup := v1.Group("/series")
		{
			seriesGroup.GET("", p.SeriesHandler.List)
			seriesGroup.GET("/:id", p.SeriesHandler.GetByID)
			seriesGroup.GET("/slug/:slug", p.SeriesHandler.GetBySlug)
			seriesGroup.POST("", auth.RequireCreate("series"), p.SeriesHandler.Create)
			seriesGroup.PUT("/:id", auth.RequireUpdate("series"), p.SeriesHandler.Update)
			seriesGroup.DELETE("/:id", auth.RequireDelete("series"), p.SeriesHandler.Delete)
			seriesGroup.POST("/:id/blogs", auth.RequireUpdate("series"), p.SeriesHandler.AddBlog)
			seriesGroup.DELETE("/:id/blogs/:blogId", auth.RequireUpdate("series"), p.SeriesHandler.RemoveBlog)
		}

		// Comments (for update/delete)
		comments := v1.Group("/comments")
		{
			comments.PUT("/:id", auth.RequireUpdate("comments"), p.CommentHandler.Update)
			comments.DELETE("/:id", auth.RequireDelete("comments"), p.CommentHandler.Delete)
		}

		// Categories
		categories := v1.Group("/categories")
		{
			categories.GET("", p.CategoryHandler.List)
			categories.GET("/:id", p.CategoryHandler.GetByID)
			categories.POST("", auth.RequireCreate("categories"), p.CategoryHandler.Create)       // Requires CREATE permission
			categories.PUT("/:id", auth.RequireUpdate("categories"), p.CategoryHandler.Update)    // Requires UPDATE permission
			categories.DELETE("/:id", auth.RequireDelete("categories"), p.CategoryHandler.Delete) // Requires DELETE permission
		}

		// Tags
		tags := v1.Group("/tags")
		{
			tags.GET("", p.TagHandler.List)
			tags.GET("/popular", p.RecommendationHandler.GetPopularTags) // Popular tags
			tags.GET("/:id", p.TagHandler.GetByID)
			tags.POST("", auth.RequireCreate("tags"), p.TagHandler.Create)       // Requires CREATE permission
			tags.PUT("/:id", auth.RequireUpdate("tags"), p.TagHandler.Update)    // Requires UPDATE permission
			tags.DELETE("/:id", auth.RequireDelete("tags"), p.TagHandler.Delete) // Requires DELETE permission
		}

		// Authors & Subscriptions
		authors := v1.Group("/authors")
		{
			authors.GET("/:authorId/subscribers", p.SubscriptionHandler.GetSubscribers)
			authors.GET("/:authorId/subscribers/count", p.SubscriptionHandler.CountSubscribers)
			authors.POST("/:authorId/subscribe", p.SubscriptionHandler.Subscribe)
			authors.POST("/:authorId/unsubscribe", p.SubscriptionHandler.Unsubscribe)
		}

		// My subscriptions
		v1.GET("/subscriptions", p.SubscriptionHandler.GetMySubscriptions)

		// My bookmarks
		v1.GET("/bookmarks", p.BookmarkHandler.List)

		// My reading history
		v1.GET("/me/history", p.ReadingHistoryHandler.GetHistory)

		// Unified Subscription/Follow API (users can follow/subscribe to each other)
		v1.GET("/users/:userId/followers", p.SubscriptionHandler.GetSubscribers)
		v1.GET("/users/:userId/following", p.SubscriptionHandler.GetUserSubscriptions)
		v1.GET("/users/:userId/follow-counts", p.SubscriptionHandler.GetSubscriptionCounts)
		v1.POST("/users/:userId/follow", p.SubscriptionHandler.Subscribe)
		v1.DELETE("/users/:userId/follow", p.SubscriptionHandler.Unsubscribe)

		// Rankings
		rankings := v1.Group("/rankings")
		{
			rankings.GET("/trending", p.RankingHandler.GetTrending)
			rankings.GET("/top", p.RankingHandler.GetTop)
			rankings.GET("/users/:userId", p.RankingHandler.GetUserRanking)
			rankings.POST("/recalculate", auth.RequireAdmin("rankings"), p.RankingHandler.RecalculateScores)
		}

		// Fraud Detection & Risk Management
		// User risk score and badge
		v1.GET("/users/:id/risk-score", p.FraudHandler.GetUserRiskScore)
		v1.GET("/users/:id/badge", p.FraudHandler.GetUserBadgeStatus)
		v1.GET("/users/:id/bot-notifications", p.FraudHandler.GetUserBotNotifications)

		// Admin fraud dashboard
		admin := v1.Group("/admin", auth.RequireAdmin("fraud"))
		{
			admin.GET("/fraud-dashboard", p.FraudHandler.GetFraudDashboard)
			admin.POST("/users/:id/review", p.FraudHandler.ReviewUser)
			admin.POST("/users/:id/ban", p.FraudHandler.BanUser)
		}

		// Analytics
		v1.GET("/analytics/fraud-trends", auth.RequireAdmin("analytics"), p.FraudHandler.GetFraudTrends)

		// Batch operations
		v1.POST("/followers/batch-analyze", auth.RequireAdmin("fraud"), p.FraudHandler.TriggerBatchAnalysis)

		// Notifications
		v1.POST("/notifications/:id/read", p.FraudHandler.MarkNotificationAsRead)
	}

	return engine
}
