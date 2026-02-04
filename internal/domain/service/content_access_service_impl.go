package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

type contentAccessService struct {
	tagTierService TagTierService
	subRepo        repository.SubscriptionRepository
	planRepo       repository.SubscriptionPlanRepository
}

// NewContentAccessService creates a new ContentAccessService instance
func NewContentAccessService(
	tagTierService TagTierService,
	subRepo repository.SubscriptionRepository,
	planRepo repository.SubscriptionPlanRepository,
) ContentAccessService {
	return &contentAccessService{
		tagTierService: tagTierService,
		subRepo:        subRepo,
		planRepo:       planRepo,
	}
}

// CheckBlogAccess determines if a user can access a specific blog
func (s *contentAccessService) CheckBlogAccess(
	ctx context.Context,
	blogID uuid.UUID,
	userID *uuid.UUID,
) (*AccessResult, error) {
	// Get required tier for the blog (also returns authorID)
	requiredTier, authorID, err := s.tagTierService.GetRequiredTierForBlog(ctx, blogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get required tier: %w", err)
	}

	// Determine user tier (FREE for anonymous users)
	userTier := entity.TierFree
	if userID != nil {
		// For logged-in users, get their subscription tier for the blog's author
		subscription, err := s.subRepo.FindActiveSubscription(ctx, *userID, authorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user tier: %w", err)
		}

		if subscription != nil {
			// Convert tier string to enum
			tier := entity.SubscriptionTier(subscription.Tier)
			if tier.IsValid() {
				userTier = tier
			}
		}
	}

	// Determine if accessible based on tier levels
	accessible := userTier.Level() >= requiredTier.Level()

	result := &AccessResult{
		Accessible:   accessible,
		UserTier:     userTier,
		RequiredTier: requiredTier,
	}

	// Generate upgrade options if blocked
	if !accessible && userID != nil {
		result.UpgradeOptions = s.generateUpgradeOptions(ctx, authorID, requiredTier)
	} else if !accessible && userID == nil {
		// For anonymous users, show all available plans
		result.UpgradeOptions = s.generateUpgradeOptions(ctx, authorID, entity.TierFree)
	}

	return result, nil
}

// GetUserTier retrieves the current subscription tier for a user and author
func (s *contentAccessService) GetUserTier(
	ctx context.Context,
	userID, authorID uuid.UUID,
) (entity.SubscriptionTier, error) {
	subscription, err := s.subRepo.FindActiveSubscription(ctx, userID, authorID)
	if err != nil {
		return entity.TierFree, fmt.Errorf("failed to get user tier: %w", err)
	}

	// If no subscription, return FREE tier
	if subscription == nil {
		return entity.TierFree, nil
	}

	// Convert tier string to enum
	tier := entity.SubscriptionTier(subscription.Tier)
	if !tier.IsValid() {
		// Invalid tier - return FREE as default
		return entity.TierFree, nil
	}

	return tier, nil
}

// generateUpgradeOptions generates upgrade options from higher tiers when blocked
func (s *contentAccessService) generateUpgradeOptions(
	ctx context.Context,
	authorID uuid.UUID,
	requiredTier entity.SubscriptionTier,
) []UpgradeOption {
	// Get all active plans for the author
	plans, err := s.planRepo.FindActiveByAuthor(ctx, authorID)
	if err != nil || len(plans) == 0 {
		return []UpgradeOption{}
	}

	// Filter plans with tier level >= required tier level (plans that would grant access)
	var options []UpgradeOption
	for _, plan := range plans {
		if plan.Tier.Level() >= requiredTier.Level() && plan.IsActive {
			options = append(options, UpgradeOption{
				PlanID:       plan.ID,
				Tier:         plan.Tier,
				Price:        plan.Price.String(),
				DurationDays: plan.DurationDays,
			})
		}
	}

	// Sort by tier level (ascending)
	slices.SortFunc(options, func(a, b UpgradeOption) int {
		return a.Tier.Level() - b.Tier.Level()
	})

	return options
}
