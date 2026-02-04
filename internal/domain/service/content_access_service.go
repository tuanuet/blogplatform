package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// ContentAccessService defines the interface for content access control business logic
type ContentAccessService interface {
	// CheckBlogAccess determines if a user can access a specific blog
	// Returns access result with user tier, required tier, and upgrade options if blocked
	CheckBlogAccess(ctx context.Context, blogID uuid.UUID, userID *uuid.UUID) (*AccessResult, error)

	// GetUserTier retrieves the current subscription tier for a user and author
	// Returns FREE tier if no active subscription exists
	GetUserTier(ctx context.Context, userID, authorID uuid.UUID) (entity.SubscriptionTier, error)
}

// AccessResult represents the result of a blog access check
type AccessResult struct {
	Accessible     bool
	UserTier       entity.SubscriptionTier
	RequiredTier   entity.SubscriptionTier
	Reason         string
	UpgradeOptions []UpgradeOption
}

// UpgradeOption represents an available upgrade option for blocked content
type UpgradeOption struct {
	PlanID       uuid.UUID
	Tier         entity.SubscriptionTier
	Price        string
	DurationDays int
}
