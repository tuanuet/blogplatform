package modules

import (
	"github.com/aiagent/internal/infrastructure/adapter"
	pgRepo "github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	redisRepo "github.com/aiagent/internal/infrastructure/persistence/redis/repository"
	"go.uber.org/fx"
)

// RepositoryModule provides repository dependencies
// Constructors already return interface types, so no fx.Annotate needed
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		pgRepo.NewBlogRepository,
		pgRepo.NewCategoryRepository,
		pgRepo.NewTagRepository,
		pgRepo.NewCommentRepository,
		pgRepo.NewSubscriptionRepository,
		pgRepo.NewBookmarkRepository,
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
		pgRepo.NewFraudDetectionRepository,
	),
)
