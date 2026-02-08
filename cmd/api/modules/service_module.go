package modules

import (
	"time"

	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/adapter"
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
		service.NewRecommendationService,
		service.NewSocialAuthService,
		service.NewPaymentService,
		service.NewPlanManagementService,
		service.NewTagTierService,
		service.NewContentAccessService,
		service.NewVersionService,
		service.NewNotificationAggregator,
		service.NewNotificationDispatcher,
		// Task Runner for async tasks
		func() service.TaskRunner {
			return service.NewTaskRunner(30 * time.Second)
		},
		// Email Service
		func(userRepo repository.UserRepository, provider adapter.EmailProvider, taskRunner service.TaskRunner) service.EmailService {
			return service.NewEmailServiceImpl(userRepo, provider, taskRunner, "internal/infrastructure/email/templates")
		},
	),
)
