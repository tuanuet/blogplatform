package modules

import (
	"context"

	"github.com/aiagent/internal/application/usecase/notification"
	"github.com/aiagent/internal/infrastructure/adapter"
	"github.com/aiagent/internal/infrastructure/config"
	pgRepo "github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	redisRepo "github.com/aiagent/internal/infrastructure/persistence/redis/repository"
	"go.uber.org/fx"
)

// RepositoryModule provides repository dependencies
// Constructors already return interface types, so no fx.Annotate needed
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		pgRepo.NewBlogRepository,
		pgRepo.NewBlogVersionRepository,
		pgRepo.NewCategoryRepository,
		pgRepo.NewTagRepository,
		pgRepo.NewCommentRepository,
		pgRepo.NewSubscriptionRepository,
		pgRepo.NewSubscriptionPlanRepository,
		pgRepo.NewTagTierMappingRepository,
		pgRepo.NewBookmarkRepository,
		pgRepo.NewTransactionRepository,
		pgRepo.NewUserRepository,
		pgRepo.NewRoleRepository,
		pgRepo.NewUserVelocityScoreRepository,
		pgRepo.NewUserRankingHistoryRepository,
		pgRepo.NewUserFollowerSnapshotRepository,
		pgRepo.NewSeriesRepository,
		pgRepo.NewReadingHistoryRepository,
		pgRepo.NewSocialAccountRepository,
		redisRepo.NewSessionRepository,
		adapter.NewSystemRepository,
		adapter.NewSePayAdapter,
		pgRepo.NewFraudDetectionRepository,
		pgRepo.NewUserSeriesPurchaseRepository,
		// Notification Repositories
		pgRepo.NewNotificationRepository,
		pgRepo.NewDeviceTokenRepository,
		pgRepo.NewNotificationPreferenceRepository,
		// Firebase Adapter Client
		func(cfg *config.Config) (adapter.FirebaseClient, error) {
			return adapter.NewFirebaseClient(context.Background(), cfg.Firebase.ProjectID, cfg.Firebase.ServiceAccountPath, cfg.Firebase.Enabled)
		},
	),
	fx.Provide(
		fx.Annotate(
			adapter.NewFirebaseAdapter,
			fx.As(new(notification.NotificationAdapter)),
		),
	),
)
