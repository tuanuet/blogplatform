package modules

import (
	"github.com/aiagent/internal/application/usecase/admin"
	"github.com/aiagent/internal/application/usecase/auth"
	"github.com/aiagent/internal/application/usecase/blog"
	"github.com/aiagent/internal/application/usecase/bookmark"
	"github.com/aiagent/internal/application/usecase/category"
	"github.com/aiagent/internal/application/usecase/comment"
	"github.com/aiagent/internal/application/usecase/health"
	"github.com/aiagent/internal/application/usecase/notification"
	"github.com/aiagent/internal/application/usecase/permission"
	"github.com/aiagent/internal/application/usecase/profile"
	"github.com/aiagent/internal/application/usecase/ranking"
	"github.com/aiagent/internal/application/usecase/reading_history"
	"github.com/aiagent/internal/application/usecase/recommendation"
	"github.com/aiagent/internal/application/usecase/role"
	"github.com/aiagent/internal/application/usecase/series"
	"github.com/aiagent/internal/application/usecase/subscription"
	"github.com/aiagent/internal/application/usecase/tag"
	"go.uber.org/fx"
)

// UseCaseModule provides application use case dependencies
var UseCaseModule = fx.Module("usecase",
	fx.Provide(
		admin.NewAdminUseCase,
		auth.NewAuthUseCase,
		blog.NewBlogUseCase,
		bookmark.NewBookmarkUseCase,
		category.NewCategoryUseCase,
		comment.NewCommentUseCase,
		health.NewHealthUseCase,
		notification.NewNotificationUseCase,
		permission.NewPermissionUseCase,
		profile.NewProfileUseCase,
		ranking.NewRankingUseCase,
		reading_history.NewReadingHistoryUseCase,
		recommendation.NewRecommendationUseCase,
		role.NewRoleUseCase,
		series.NewSeriesUseCase,
		subscription.NewSubscriptionUseCase,
		tag.NewTagUseCase,
	),
)
