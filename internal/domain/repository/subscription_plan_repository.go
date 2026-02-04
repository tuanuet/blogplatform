package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// SubscriptionPlanRepository defines the interface for subscription plan data operations
type SubscriptionPlanRepository interface {
	// Create creates a new subscription plan
	Create(ctx context.Context, plan *entity.SubscriptionPlan) error

	// FindByID retrieves a subscription plan by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.SubscriptionPlan, error)

	// FindByAuthorAndTier retrieves a plan by author ID and tier
	FindByAuthorAndTier(ctx context.Context, authorID uuid.UUID, tier entity.SubscriptionTier) (*entity.SubscriptionPlan, error)

	// FindByAuthor retrieves all plans for an author
	FindByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.SubscriptionPlan, error)

	// FindActiveByAuthor retrieves all active plans for an author
	FindActiveByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.SubscriptionPlan, error)

	// Update updates an existing subscription plan
	Update(ctx context.Context, plan *entity.SubscriptionPlan) error

	// Upsert creates or updates a subscription plan based on (author_id, tier) uniqueness
	Upsert(ctx context.Context, plan *entity.SubscriptionPlan) error

	// Delete soft deletes a subscription plan
	Delete(ctx context.Context, id uuid.UUID) error

	// WithTx returns a new repository with the given transaction
	WithTx(tx interface{}) SubscriptionPlanRepository
}
