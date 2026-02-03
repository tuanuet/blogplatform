package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type planManagementService struct {
	planRepo    repository.SubscriptionPlanRepository
	tagTierRepo repository.TagTierMappingRepository
}

// NewPlanManagementService creates a new PlanManagementService instance
func NewPlanManagementService(
	planRepo repository.SubscriptionPlanRepository,
	tagTierRepo repository.TagTierMappingRepository,
) PlanManagementService {
	return &planManagementService{
		planRepo:    planRepo,
		tagTierRepo: tagTierRepo,
	}
}

// UpsertPlans creates or updates multiple subscription plans
func (s *planManagementService) UpsertPlans(
	ctx context.Context,
	authorID uuid.UUID,
	plans []CreatePlanDTO,
) ([]entity.SubscriptionPlan, []string, error) {
	if len(plans) == 0 {
		return nil, nil, nil
	}

	// Validate tier values and check for duplicates
	seenTiers := make(map[entity.SubscriptionTier]bool)
	for _, p := range plans {
		if !p.Tier.IsValid() {
			return nil, nil, fmt.Errorf("invalid tier: %s", p.Tier)
		}
		if seenTiers[p.Tier] {
			return nil, nil, fmt.Errorf("duplicate tier: %s", p.Tier)
		}
		seenTiers[p.Tier] = true
	}

	// Validate price hierarchy (warnings only)
	warnings := validatePriceHierarchy(plans)

	// Create plan entities and upsert
	result := make([]entity.SubscriptionPlan, 0, len(plans))
	for _, p := range plans {
		plan := &entity.SubscriptionPlan{
			AuthorID:     authorID,
			Tier:         p.Tier,
			Price:        p.Price,
			DurationDays: p.DurationDays,
			Name:         p.Name,
			Description:  p.Description,
			IsActive:     true,
		}

		if err := s.planRepo.Upsert(ctx, plan); err != nil {
			return nil, nil, fmt.Errorf("failed to upsert plan for tier %s: %w", p.Tier, err)
		}
		result = append(result, *plan)
	}

	return result, warnings, nil
}

// GetAuthorPlans retrieves all plans for an author with tag information
func (s *planManagementService) GetAuthorPlans(
	ctx context.Context,
	authorID uuid.UUID,
) ([]PlanWithTags, error) {
	plans, err := s.planRepo.FindActiveByAuthor(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plans: %w", err)
	}

	// Get all tag-tier mappings for the author
	tagMappings, err := s.tagTierRepo.FindByAuthor(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tag mappings: %w", err)
	}

	// Build tag count per tier
	tagCountByTier := make(map[entity.SubscriptionTier]int)
	for _, m := range tagMappings {
		tagCountByTier[m.RequiredTier]++
	}

	// Build result with all 4 tiers including FREE
	allTiers := []entity.SubscriptionTier{
		entity.TierFree,
		entity.TierBronze,
		entity.TierSilver,
		entity.TierGold,
	}

	result := make([]PlanWithTags, 0, 4)
	for _, tier := range allTiers {
		var plan entity.SubscriptionPlan
		found := false

		for _, p := range plans {
			if p.Tier == tier {
				plan = p
				found = true
				break
			}
		}

		if !found {
			// Create empty plan for missing tiers (especially FREE)
			plan = entity.SubscriptionPlan{
				ID:           uuid.Nil,
				AuthorID:     authorID,
				Tier:         tier,
				Price:        decimal.Zero,
				DurationDays: 0,
				IsActive:     false,
			}
		}

		result = append(result, PlanWithTags{
			Plan:     plan,
			Tags:     []string{},
			TagCount: tagCountByTier[tier],
		})
	}

	// Sort by tier level
	slices.SortFunc(result, func(a, b PlanWithTags) int {
		return a.Plan.Tier.Level() - b.Plan.Tier.Level()
	})

	return result, nil
}

// DeactivatePlan deactivates a plan without deleting it
func (s *planManagementService) DeactivatePlan(
	ctx context.Context,
	authorID uuid.UUID,
	tier entity.SubscriptionTier,
) error {
	if tier == entity.TierFree {
		return fmt.Errorf("cannot deactivate FREE tier")
	}

	plan, err := s.planRepo.FindByAuthorAndTier(ctx, authorID, tier)
	if err != nil {
		return fmt.Errorf("failed to find plan: %w", err)
	}
	if plan == nil {
		return fmt.Errorf("plan not found for tier %s", tier)
	}

	plan.IsActive = false
	if err := s.planRepo.Update(ctx, plan); err != nil {
		return fmt.Errorf("failed to deactivate plan: %w", err)
	}

	return nil
}

// ActivatePlan activates a previously deactivated plan
func (s *planManagementService) ActivatePlan(
	ctx context.Context,
	authorID uuid.UUID,
	tier entity.SubscriptionTier,
) error {
	if tier == entity.TierFree {
		return fmt.Errorf("cannot activate FREE tier")
	}

	plan, err := s.planRepo.FindByAuthorAndTier(ctx, authorID, tier)
	if err != nil {
		return fmt.Errorf("failed to find plan: %w", err)
	}
	if plan == nil {
		return fmt.Errorf("plan not found for tier %s", tier)
	}

	plan.IsActive = true
	if err := s.planRepo.Update(ctx, plan); err != nil {
		return fmt.Errorf("failed to activate plan: %w", err)
	}

	return nil
}

// validatePriceHierarchy checks if price hierarchy is correct (BRONZE < SILVER < GOLD)
// Returns warnings but does not block the operation
func validatePriceHierarchy(plans []CreatePlanDTO) []string {
	warnings := []string{}
	priceMap := make(map[entity.SubscriptionTier]decimal.Decimal)

	for _, p := range plans {
		priceMap[p.Tier] = p.Price
	}

	// Check BRONZE < SILVER
	if bronze, okB := priceMap[entity.TierBronze]; okB {
		if silver, okS := priceMap[entity.TierSilver]; okS {
			if silver.LessThanOrEqual(bronze) {
				warnings = append(warnings, fmt.Sprintf("SILVER price (%s) should be > BRONZE (%s)", silver, bronze))
			}
		}
	}

	// Check SILVER < GOLD
	if silver, okS := priceMap[entity.TierSilver]; okS {
		if gold, okG := priceMap[entity.TierGold]; okG {
			if gold.LessThanOrEqual(silver) {
				warnings = append(warnings, fmt.Sprintf("GOLD price (%s) should be > SILVER (%s)", gold, silver))
			}
		}
	}

	return warnings
}
