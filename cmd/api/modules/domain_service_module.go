package modules

import (
	"github.com/aiagent/boilerplate/internal/domain/service"
	"go.uber.org/fx"
)

// DomainServiceModule provides domain service dependencies
var DomainServiceModule = fx.Module("domain_service",
	fx.Provide(
		service.NewRoleService,
		service.NewPermissionService,
		service.NewUserService,
		service.NewSystemService,
		service.NewCategoryService,
		service.NewTagService,
		service.NewSubscriptionService,
		service.NewCommentService,
		service.NewBlogService,
		service.NewRankingService,
		service.NewFraudDetectionService,
		service.NewNotificationService,
		service.NewBotDetectionAlgorithm,
		service.NewBatchJobService,
	),
)
