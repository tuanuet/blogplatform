package modules

import (
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/blog"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/category"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/comment"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/fraud"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/health"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/profile"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/ranking"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/role"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/subscription"
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler/tag"
	"go.uber.org/fx"
)

// HandlerModule provides HTTP handler dependencies
// Uses constructors directly - no wrapper functions needed
var HandlerModule = fx.Module("handler",
	fx.Provide(
		health.NewHealthHandler,
		blog.NewBlogHandler,
		category.NewCategoryHandler,
		tag.NewTagHandler,
		comment.NewCommentHandler,
		subscription.NewSubscriptionHandler,
		profile.NewProfileHandler,
		role.NewRoleHandler,
		ranking.NewRankingHandler,
		fraud.NewFraudHandler,
	),
)
