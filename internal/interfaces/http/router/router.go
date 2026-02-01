package router

import (
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler"
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

	HealthHandler       *handler.HealthHandler
	BlogHandler         *handler.BlogHandler
	CategoryHandler     *handler.CategoryHandler
	TagHandler          *handler.TagHandler
	CommentHandler      *handler.CommentHandler
	SubscriptionHandler *handler.SubscriptionHandler
	ProfileHandler      *handler.ProfileHandler
	RoleHandler         *handler.RoleHandler
	RoleUseCase         usecase.RoleUseCase // For authorization middleware
	Config              *config.Config
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
			blogs.GET("/:id", p.BlogHandler.GetByID)
			blogs.POST("", auth.RequireCreate("blogs"), p.BlogHandler.Create)                  // Requires CREATE permission
			blogs.PUT("/:id", auth.RequireUpdate("blogs"), p.BlogHandler.Update)               // Requires UPDATE permission
			blogs.DELETE("/:id", auth.RequireDelete("blogs"), p.BlogHandler.Delete)            // Requires DELETE permission
			blogs.POST("/:id/publish", auth.RequireUpdate("blogs"), p.BlogHandler.Publish)     // Requires UPDATE permission
			blogs.POST("/:id/unpublish", auth.RequireUpdate("blogs"), p.BlogHandler.Unpublish) // Requires UPDATE permission
			blogs.POST("/:id/reaction", p.BlogHandler.React)                                   // Authenticated users

			// Blog comments
			blogs.GET("/:id/comments", p.CommentHandler.GetByBlogID)
			blogs.POST("/:id/comments", auth.RequireCreate("comments"), p.CommentHandler.Create)
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
	}

	return engine
}
