package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// TagTierService defines the interface for tag-tier mapping business logic
type TagTierService interface {
	// AssignTagToTier assigns a tag to a subscription tier
	// Returns the mapping and count of affected blogs
	AssignTagToTier(ctx context.Context, authorID, tagID uuid.UUID, tier entity.SubscriptionTier) (*entity.TagTierMapping, int64, error)

	// UnassignTagFromTier removes tier requirement from a tag
	// Returns count of affected blogs
	UnassignTagFromTier(ctx context.Context, authorID, tagID uuid.UUID) (int64, error)

	// GetAuthorTagTiers retrieves all tag-tier mappings for an author
	GetAuthorTagTiers(ctx context.Context, authorID uuid.UUID) ([]TagTierWithCount, error)

	// GetRequiredTierForBlog determines the highest required tier for a blog based on its tags
	// Returns: required tier, blog author ID, error
	GetRequiredTierForBlog(ctx context.Context, blogID uuid.UUID) (entity.SubscriptionTier, uuid.UUID, error)
}

// TagTierWithCount represents a tag-tier mapping with blog count
type TagTierWithCount struct {
	Mapping   entity.TagTierMapping
	TagName   string
	BlogCount int64
}
