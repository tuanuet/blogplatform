package router

import (
	"time"

	roleUseCase "github.com/aiagent/internal/application/usecase/role"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/config"
	"github.com/aiagent/internal/interfaces/http/handler/admin"
	"github.com/aiagent/internal/interfaces/http/handler/auth"
	"github.com/aiagent/internal/interfaces/http/handler/blog"
	"github.com/aiagent/internal/interfaces/http/handler/bookmark"
	"github.com/aiagent/internal/interfaces/http/handler/category"
	"github.com/aiagent/internal/interfaces/http/handler/comment"
	"github.com/aiagent/internal/interfaces/http/handler/fraud"
	"github.com/aiagent/internal/interfaces/http/handler/health"
	"github.com/aiagent/internal/interfaces/http/handler/payment"
	"github.com/aiagent/internal/interfaces/http/handler/plan"
	"github.com/aiagent/internal/interfaces/http/handler/profile"
	"github.com/aiagent/internal/interfaces/http/handler/ranking"
	"github.com/aiagent/internal/interfaces/http/handler/reading_history"
	"github.com/aiagent/internal/interfaces/http/handler/recommendation"
	"github.com/aiagent/internal/interfaces/http/handler/role"
	"github.com/aiagent/internal/interfaces/http/handler/series"
	"github.com/aiagent/internal/interfaces/http/handler/subscription"
	"github.com/aiagent/internal/interfaces/http/handler/tag"
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"

	// Import generated swagger docs
	_ "github.com/aiagent/docs"
)

// Params groups all handler dependencies for the router (FX dependency injection)
type Params struct {
	fx.In

	HealthHandler         health.HealthHandler
	AdminHandler          admin.AdminHandler
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
	PaymentHandler        payment.PaymentHandler
	WebhookHandler        payment.WebhookHandler
	PlanHandler           plan.PlanHandler
	AuthHandler           auth.AuthHandler
	SessionRepository     repository.SessionRepository
	RedisClient           *redis.Client
	RoleUseCase           roleUseCase.RoleUseCase // For authorization middleware
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
	sessionAuth := middleware.SessionAuth(p.SessionRepository)
	rateLimit := middleware.RateLimit(p.RedisClient, 100, time.Minute)

	// Serve static files for avatar uploads
	engine.Static("/uploads", "./uploads")

	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		RegisterHealthRoutes(engine, v1, p)
		RegisterAuthRoutes(v1, p, rateLimit, sessionAuth)
		RegisterProfileRoutes(v1, p, sessionAuth)
		RegisterUserRoutes(v1, p, auth, sessionAuth)
		RegisterRoleRoutes(v1, p, auth, sessionAuth)
		RegisterBlogRoutes(v1, p, auth, sessionAuth)
		RegisterSeriesRoutes(v1, p, auth, sessionAuth)
		RegisterCommentRoutes(v1, p, auth, sessionAuth)
		RegisterCategoryRoutes(v1, p, auth, sessionAuth)
		RegisterTagRoutes(v1, p, auth, sessionAuth)
		RegisterSubscriptionRoutes(v1, p, sessionAuth)
		RegisterBookmarkRoutes(v1, p, sessionAuth)
		RegisterReadingHistoryRoutes(v1, p, sessionAuth)
		RegisterRankingRoutes(v1, p, auth, sessionAuth)
		RegisterAdminRoutes(v1, p, auth, sessionAuth)
		RegisterFraudRoutes(v1, p, auth, sessionAuth)

		// Payment & Webhooks
		RegisterPaymentRoutes(v1, p.PaymentHandler, p.WebhookHandler, sessionAuth)

		// Plan routes (multi-tier subscription)
		RegisterPlanRoutes(v1, p.PlanHandler, sessionAuth)
	}

	return engine
}
