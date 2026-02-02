package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

// Cache defines the interface for cache operations
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	DeleteByPattern(ctx context.Context, pattern string) error
	Close() error
	HealthCheck(ctx context.Context) error
	Client() *redis.Client
}
