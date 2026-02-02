package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// ReadingHistoryRepository defines the interface for reading history persistence
type ReadingHistoryRepository interface {
	// Upsert records a view for a user on a blog. If the record exists, it updates LastReadAt.
	Upsert(ctx context.Context, history *entity.UserReadingHistory) error

	// GetRecentByUserID retrieves the recently viewed blogs for a user.
	// It returns a list of UserReadingHistory with the associated Blog preloaded.
	GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.UserReadingHistory, error)
}
