package modules

import (
	"github.com/aiagent/internal/application/usecase/payment"
	"github.com/aiagent/internal/infrastructure/config"
	"github.com/aiagent/internal/interfaces/http/handler/admin"
	"github.com/aiagent/internal/interfaces/http/handler/auth"
	"github.com/aiagent/internal/interfaces/http/handler/blog"
	"github.com/aiagent/internal/interfaces/http/handler/bookmark"
	"github.com/aiagent/internal/interfaces/http/handler/category"
	"github.com/aiagent/internal/interfaces/http/handler/comment"
	"github.com/aiagent/internal/interfaces/http/handler/fraud"
	"github.com/aiagent/internal/interfaces/http/handler/health"
	"github.com/aiagent/internal/interfaces/http/handler/notification"
	paymentH "github.com/aiagent/internal/interfaces/http/handler/payment"
	"github.com/aiagent/internal/interfaces/http/handler/profile"
	"github.com/aiagent/internal/interfaces/http/handler/ranking"
	"github.com/aiagent/internal/interfaces/http/handler/reading_history"
	"github.com/aiagent/internal/interfaces/http/handler/recommendation"
	"github.com/aiagent/internal/interfaces/http/handler/role"
	"github.com/aiagent/internal/interfaces/http/handler/series"
	"github.com/aiagent/internal/interfaces/http/handler/subscription"
	"github.com/aiagent/internal/interfaces/http/handler/tag"
	"go.uber.org/fx"
)

// HandlerModule provides HTTP handler dependencies
// Uses constructors directly - no wrapper functions needed
var HandlerModule = fx.Module("handler",
	fx.Provide(
		health.NewHealthHandler,
		admin.NewAdminHandler,
		blog.NewBlogHandler,
		bookmark.NewBookmarkHandler,
		category.NewCategoryHandler,
		tag.NewTagHandler,
		comment.NewCommentHandler,
		subscription.NewSubscriptionHandler,
		profile.NewProfileHandler,
		role.NewRoleHandler,
		series.NewSeriesHandler,
		ranking.NewRankingHandler,
		reading_history.NewReadingHistoryHandler,
		fraud.NewFraudHandler,
		recommendation.NewRecommendationHandler,
		auth.NewAuthHandler,
		notification.NewNotificationHandler,
		paymentH.NewPaymentHandler,
		func(cfg *config.Config, uc payment.ProcessWebhookUseCase) paymentH.WebhookHandler {
			return paymentH.NewWebhookHandler(uc, cfg.SePay.APIKey)
		},
	),
)
