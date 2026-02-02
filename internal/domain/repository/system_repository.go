package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "context"

// SystemRepository defines the interface for system health checks
type SystemRepository interface {
	CheckDatabase(ctx context.Context) error
	CheckRedis(ctx context.Context) error
}
