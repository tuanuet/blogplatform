package modules

import (
	"github.com/aiagent/boilerplate/internal/interfaces/http/handler"
	"go.uber.org/fx"
)

// HandlerModule provides HTTP handler dependencies
// Uses constructors directly - no wrapper functions needed
var HandlerModule = fx.Module("handler",
	fx.Provide(
		handler.NewHealthHandler,
		handler.NewBlogHandler,
		handler.NewCategoryHandler,
		handler.NewTagHandler,
		handler.NewCommentHandler,
		handler.NewSubscriptionHandler,
		handler.NewProfileHandler,
		handler.NewRoleHandler,
	),
)
