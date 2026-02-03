package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.Subscription) error
	Delete(ctx context.Context, subscriberID, authorID uuid.UUID) error
	Exists(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error)
	FindBySubscriber(ctx context.Context, subscriberID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.Subscription], error)
	FindByAuthor(ctx context.Context, authorID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.Subscription], error)
	CountSubscribers(ctx context.Context, authorID uuid.UUID) (int64, error)
	CountBySubscriber(ctx context.Context, subscriberID uuid.UUID) (int64, error)
	UpdateExpiry(ctx context.Context, userID, authorID uuid.UUID, expiresAt time.Time, tier string) error
	FindActiveSubscription(ctx context.Context, userID, authorID uuid.UUID) (*entity.Subscription, error)

	// WithTx returns a new repository with the given transaction
	WithTx(tx interface{}) SubscriptionRepository
}
