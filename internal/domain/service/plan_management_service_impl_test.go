package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPlanManagementService_UpsertPlans_PriceHierarchyValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("valid_price_hierarchy_bronze_silver_gold", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(50000), DurationDays: 30},
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(100000), DurationDays: 30},
			{Tier: entity.TierGold, Price: decimal.NewFromInt(200000), DurationDays: 30},
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(3).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("warning_silver_price_less_than_or_equal_to_bronze", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(100000), DurationDays: 30},
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(100000), DurationDays: 30}, // Same price
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0], "SILVER price")
		assert.Contains(t, warnings[0], "BRONZE")
	})

	t.Run("warning_silver_price_less_than_bronze", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(150000), DurationDays: 30},
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(100000), DurationDays: 30}, // Lower price
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0], "SILVER price")
		assert.Contains(t, warnings[0], "BRONZE")
	})

	t.Run("warning_gold_price_less_than_or_equal_to_silver", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(100000), DurationDays: 30},
			{Tier: entity.TierGold, Price: decimal.NewFromInt(100000), DurationDays: 30}, // Same price
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0], "GOLD price")
		assert.Contains(t, warnings[0], "SILVER")
	})

	t.Run("multiple_warnings_for_hierarchy_violations", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(200000), DurationDays: 30},
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(150000), DurationDays: 30}, // Lower
			{Tier: entity.TierGold, Price: decimal.NewFromInt(100000), DurationDays: 30},   // Lower
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(3).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Len(t, warnings, 2)
	})
}

func TestPlanManagementService_UpsertPlans_TierValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("valid_tier_values", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierFree, Price: decimal.Zero, DurationDays: 0},
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(50000), DurationDays: 30},
			{Tier: entity.TierSilver, Price: decimal.NewFromInt(100000), DurationDays: 30},
			{Tier: entity.TierGold, Price: decimal.NewFromInt(200000), DurationDays: 30},
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Times(4).Return(nil)

		_, warnings, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("invalid_tier_value", func(t *testing.T) {
		invalidTier := entity.SubscriptionTier("INVALID")
		plansDTO := []service.CreatePlanDTO{
			{Tier: invalidTier, Price: decimal.NewFromInt(100000), DurationDays: 30},
		}

		_, _, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid tier")
	})

	t.Run("cannot_create_multiple_free_tier_plans", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierFree, Price: decimal.Zero, DurationDays: 0},
			{Tier: entity.TierFree, Price: decimal.Zero, DurationDays: 0}, // Duplicate
		}

		_, _, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate tier")
	})
}

func TestPlanManagementService_UpsertPlans_RepositoryErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("upsert_failure", func(t *testing.T) {
		plansDTO := []service.CreatePlanDTO{
			{Tier: entity.TierBronze, Price: decimal.NewFromInt(50000), DurationDays: 30},
		}

		mockPlanRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		_, _, err := svc.UpsertPlans(ctx, authorID, plansDTO)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to upsert plan")
	})
}

func TestPlanManagementService_GetAuthorPlans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("returns_all_four_tiers_including_free_as_empty", func(t *testing.T) {
		bronzeID := uuid.New()
		plans := []entity.SubscriptionPlan{
			{
				ID:           bronzeID,
				AuthorID:     authorID,
				Tier:         entity.TierBronze,
				Price:        decimal.NewFromInt(50000),
				DurationDays: 30,
				IsActive:     true,
			},
			{
				ID:           uuid.New(),
				AuthorID:     authorID,
				Tier:         entity.TierSilver,
				Price:        decimal.NewFromInt(100000),
				DurationDays: 30,
				IsActive:     true,
			},
			{
				ID:           uuid.New(),
				AuthorID:     authorID,
				Tier:         entity.TierGold,
				Price:        decimal.NewFromInt(200000),
				DurationDays: 30,
				IsActive:     true,
			},
		}

		tagMappings := []entity.TagTierMapping{
			{ID: uuid.New(), AuthorID: authorID, TagID: uuid.New(), RequiredTier: entity.TierBronze},
			{ID: uuid.New(), AuthorID: authorID, TagID: uuid.New(), RequiredTier: entity.TierBronze},
			{ID: uuid.New(), AuthorID: authorID, TagID: uuid.New(), RequiredTier: entity.TierSilver},
		}

		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(plans, nil)
		mockTagTierRepo.EXPECT().FindByAuthor(ctx, authorID).Return(tagMappings, nil)

		result, err := svc.GetAuthorPlans(ctx, authorID)

		assert.NoError(t, err)
		assert.Len(t, result, 4)

		// Find FREE tier (should be empty plan)
		var freePlan *service.PlanWithTags
		var bronzePlan *service.PlanWithTags
		for i, p := range result {
			if p.Plan.Tier == entity.TierFree {
				freePlan = &result[i]
			}
			if p.Plan.Tier == entity.TierBronze {
				bronzePlan = &result[i]
			}
		}

		assert.NotNil(t, freePlan)
		assert.Equal(t, entity.TierFree, freePlan.Plan.Tier)
		assert.Equal(t, decimal.Zero, freePlan.Plan.Price)
		assert.False(t, freePlan.Plan.IsActive) // Empty FREE plan is inactive
		assert.Equal(t, 0, freePlan.TagCount)
		assert.Empty(t, freePlan.Tags)

		assert.NotNil(t, bronzePlan)
		assert.Equal(t, entity.TierBronze, bronzePlan.Plan.Tier)
		assert.Equal(t, 2, bronzePlan.TagCount) // 2 mappings for BRONZE
	})

	t.Run("returns_empty_list_on_repository_error", func(t *testing.T) {
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(nil, errors.New("database error"))

		_, err := svc.GetAuthorPlans(ctx, authorID)

		assert.Error(t, err)
	})
}

func TestPlanManagementService_DeactivatePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("successfully_deactivate_plan", func(t *testing.T) {
		existingPlan := &entity.SubscriptionPlan{
			ID:       uuid.New(),
			AuthorID: authorID,
			Tier:     entity.TierBronze,
			IsActive: true,
		}

		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierBronze).Return(existingPlan, nil)
		mockPlanRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Do(func(_ context.Context, plan *entity.SubscriptionPlan) {
			assert.False(t, plan.IsActive)
		}).Return(nil)

		err := svc.DeactivatePlan(ctx, authorID, entity.TierBronze)

		assert.NoError(t, err)
	})

	t.Run("cannot_deactivate_free_tier", func(t *testing.T) {
		err := svc.DeactivatePlan(ctx, authorID, entity.TierFree)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot deactivate FREE tier")
	})

	t.Run("plan_not_found", func(t *testing.T) {
		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierSilver).Return(nil, nil)

		err := svc.DeactivatePlan(ctx, authorID, entity.TierSilver)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("update_failure", func(t *testing.T) {
		existingPlan := &entity.SubscriptionPlan{
			ID:       uuid.New(),
			AuthorID: authorID,
			Tier:     entity.TierBronze,
			IsActive: true,
		}

		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierBronze).Return(existingPlan, nil)
		mockPlanRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		err := svc.DeactivatePlan(ctx, authorID, entity.TierBronze)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deactivate plan")
	})
}

func TestPlanManagementService_ActivatePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)
	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)

	svc := service.NewPlanManagementService(mockPlanRepo, mockTagTierRepo)

	ctx := context.Background()
	authorID := uuid.New()

	t.Run("successfully_activate_plan", func(t *testing.T) {
		existingPlan := &entity.SubscriptionPlan{
			ID:       uuid.New(),
			AuthorID: authorID,
			Tier:     entity.TierBronze,
			IsActive: false,
		}

		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierBronze).Return(existingPlan, nil)
		mockPlanRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Do(func(_ context.Context, plan *entity.SubscriptionPlan) {
			assert.True(t, plan.IsActive)
		}).Return(nil)

		err := svc.ActivatePlan(ctx, authorID, entity.TierBronze)

		assert.NoError(t, err)
	})

	t.Run("cannot_activate_free_tier", func(t *testing.T) {
		err := svc.ActivatePlan(ctx, authorID, entity.TierFree)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot activate FREE tier")
	})

	t.Run("plan_not_found", func(t *testing.T) {
		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierSilver).Return(nil, nil)

		err := svc.ActivatePlan(ctx, authorID, entity.TierSilver)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("update_failure", func(t *testing.T) {
		existingPlan := &entity.SubscriptionPlan{
			ID:       uuid.New(),
			AuthorID: authorID,
			Tier:     entity.TierBronze,
			IsActive: false,
		}

		mockPlanRepo.EXPECT().FindByAuthorAndTier(ctx, authorID, entity.TierBronze).Return(existingPlan, nil)
		mockPlanRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		err := svc.ActivatePlan(ctx, authorID, entity.TierBronze)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to activate plan")
	})
}
