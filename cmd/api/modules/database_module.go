package modules

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/infrastructure/cache"
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/internal/infrastructure/persistence/postgres"
	"github.com/aiagent/boilerplate/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DatabaseModule provides database and cache dependencies with lifecycle management
var DatabaseModule = fx.Module("database",
	fx.Provide(newDatabase, newRedisClient),
)

// newDatabase creates DB connection with cleanup on shutdown
func newDatabase(lc fx.Lifecycle, cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := postgres.NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// Auto-migrate fraud detection models
	if err := autoMigrateFraudModels(db); err != nil {
		logger.Error("Failed to auto-migrate fraud detection models: %v", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection...")
			return postgres.Close(db)
		},
	})

	logger.Info("Database connected successfully")
	return db, nil
}

// autoMigrateFraudModels migrates fraud detection related database tables
func autoMigrateFraudModels(db *gorm.DB) error {
	logger.Info("Auto-migrating fraud detection models...")

	models := []interface{}{
		&entity.FollowerEvent{},
		&entity.BotDetectionSignal{},
		&entity.UserRiskScore{},
		&entity.UserBadgeStatus{},
		&entity.AdminReview{},
		&entity.BotFollowerNotification{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	logger.Info("Fraud detection models migrated successfully")
	return nil
}

// newRedisClient creates Redis client with cleanup on shutdown
func newRedisClient(lc fx.Lifecycle, cfg *config.RedisConfig) (*cache.RedisClient, error) {
	client, err := cache.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection...")
			return client.Close()
		},
	})

	logger.Info("Redis connected successfully")
	return client, nil
}
