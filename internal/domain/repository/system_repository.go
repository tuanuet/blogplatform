package repository

import "context"

// SystemRepository defines the interface for system health checks
type SystemRepository interface {
	CheckDatabase(ctx context.Context) error
	CheckRedis(ctx context.Context) error
}
