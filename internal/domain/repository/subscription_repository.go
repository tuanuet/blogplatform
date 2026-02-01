package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
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
}
