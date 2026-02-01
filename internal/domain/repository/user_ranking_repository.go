package repository

import (
	"context"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// UserVelocityScoreRepository defines the interface for user velocity score operations
type UserVelocityScoreRepository interface {
	// FindByUserID finds a velocity score by user ID
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserVelocityScore, error)

	// FindByRank finds a velocity score by rank position
	FindByRank(ctx context.Context, rank int) (*entity.UserVelocityScore, error)

	// ListTopRanked lists the top N ranked users
	ListTopRanked(ctx context.Context, limit int) ([]entity.UserVelocityScore, error)

	// ListRanked lists ranked users with pagination
	ListRanked(ctx context.Context, pagination Pagination) (*PaginatedResult[entity.UserVelocityScore], error)

	// Save saves or updates a velocity score
	Save(ctx context.Context, score *entity.UserVelocityScore) error

	// UpdateRankPosition updates the rank position for a user
	UpdateRankPosition(ctx context.Context, userID uuid.UUID, rank int) error

	// DeleteByUserID deletes a velocity score by user ID
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// Count returns the total count of ranked users
	Count(ctx context.Context) (int64, error)
}

// UserRankingHistoryRepository defines the interface for user ranking history operations
type UserRankingHistoryRepository interface {
	// Create creates a new ranking history entry
	Create(ctx context.Context, history *entity.UserRankingHistory) error

	// FindByUserID finds ranking history for a user
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]entity.UserRankingHistory, error)

	// FindLatestByUserID finds the most recent ranking history entry for a user
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserRankingHistory, error)

	// ListByDateRange lists ranking history within a date range
	ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.UserRankingHistory, error)
}

// UserFollowerSnapshotRepository defines the interface for follower snapshot operations
type UserFollowerSnapshotRepository interface {
	// Create creates a new follower snapshot
	Create(ctx context.Context, snapshot *entity.UserFollowerSnapshot) error

	// FindByUserIDAndDate finds a snapshot by user ID and date
	FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.UserFollowerSnapshot, error)

	// FindLatestByUserID finds the most recent snapshot for a user
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserFollowerSnapshot, error)

	// FindByUserIDAndDateRange finds snapshots within a date range
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.UserFollowerSnapshot, error)

	// CountFollowers counts current followers for a user
	CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error)

	// CountBlogs counts published blogs for a user in a date range
	CountBlogs(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (int64, error)
}
