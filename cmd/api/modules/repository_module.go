package modules

import (
	"github.com/aiagent/boilerplate/internal/infrastructure/adapter"
	"github.com/aiagent/boilerplate/internal/infrastructure/persistence"
	pgRepo "github.com/aiagent/boilerplate/internal/infrastructure/persistence/postgres/repository"
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
		pgRepo.NewUserRepository,
		pgRepo.NewRoleRepository,
		pgRepo.NewUserVelocityScoreRepository,
		pgRepo.NewUserRankingHistoryRepository,
		pgRepo.NewUserFollowerSnapshotRepository,
		adapter.NewSystemRepository,
		persistence.NewFraudDetectionRepository,
	),
)
