package adapter

import (
	"context"

	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/cache"
	"github.com/aiagent/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type systemRepository struct {
	db    *gorm.DB
	redis *cache.RedisClient
}

// NewSystemRepository creates a new system repository
func NewSystemRepository(db *gorm.DB, redis *cache.RedisClient) repository.SystemRepository {
	return &systemRepository{
		db:    db,
		redis: redis,
	}
}

func (r *systemRepository) CheckDatabase(ctx context.Context) error {
	return postgres.HealthCheck(r.db)
}

func (r *systemRepository) CheckRedis(ctx context.Context) error {
	return r.redis.HealthCheck(ctx)
}
