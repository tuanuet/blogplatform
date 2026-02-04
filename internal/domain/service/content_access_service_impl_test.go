package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestContentAccessService_CheckBlogAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierService := serviceMocks.NewMockTagTierService(ctrl)
	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)

	svc := service.NewContentAccessService(mockTagTierService, mockSubRepo, mockPlanRepo)

	ctx := context.Background()
	blogID := uuid.New()
	authorID := uuid.New()
	userID := uuid.New()
	bronzePlanID := uuid.New()
	silverPlanID := uuid.New()
	goldPlanID := uuid.New()

	// Helper to create a subscription plan
	createPlan := func(id uuid.UUID, tier entity.SubscriptionTier, price string) entity.SubscriptionPlan {
		d, _ := decimal.NewFromString(price)
		return entity.SubscriptionPlan{
			ID:           id,
			AuthorID:     authorID,
			Tier:         tier,
			Price:        d,
			DurationDays: 30,
			IsActive:     true,
		}
	}

	// Define test access matrix: user tier × required tier
	t.Run("access_matrix", func(t *testing.T) {
		tests := []struct {
			name             string
			userTier         entity.SubscriptionTier
			requiredTier     entity.SubscriptionTier
			expectAccessible bool
		}{
			{"FREE user vs FREE requirement", entity.TierFree, entity.TierFree, true},
			{"FREE user vs BRONZE requirement", entity.TierFree, entity.TierBronze, false},
			{"FREE user vs SILVER requirement", entity.TierFree, entity.TierSilver, false},
			{"FREE user vs GOLD requirement", entity.TierFree, entity.TierGold, false},
			{"BRONZE user vs FREE requirement", entity.TierBronze, entity.TierFree, true},
			{"BRONZE user vs BRONZE requirement", entity.TierBronze, entity.TierBronze, true},
			{"BRONZE user vs SILVER requirement", entity.TierBronze, entity.TierSilver, false},
			{"BRONZE user vs GOLD requirement", entity.TierBronze, entity.TierGold, false},
			{"SILVER user vs FREE requirement", entity.TierSilver, entity.TierFree, true},
			{"SILVER user vs BRONZE requirement", entity.TierSilver, entity.TierBronze, true},
			{"SILVER user vs SILVER requirement", entity.TierSilver, entity.TierSilver, true},
			{"SILVER user vs GOLD requirement", entity.TierSilver, entity.TierGold, false},
			{"GOLD user vs FREE requirement", entity.TierGold, entity.TierFree, true},
			{"GOLD user vs BRONZE requirement", entity.TierGold, entity.TierBronze, true},
			{"GOLD user vs SILVER requirement", entity.TierGold, entity.TierSilver, true},
			{"GOLD user vs GOLD requirement", entity.TierGold, entity.TierGold, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Setup mocks
				mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(tt.requiredTier, authorID, nil)

				if tt.userTier != entity.TierFree {
					sub := &entity.Subscription{
						SubscriberID: userID,
						AuthorID:     authorID,
						Tier:         string(tt.userTier),
					}
					mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)
				} else {
					// For FREE tier, return nil subscription (or expired)
					mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)
				}

				// If content is blocked, mock plan repo for upgrade options
				if !tt.expectAccessible {
					plans := []entity.SubscriptionPlan{
						createPlan(bronzePlanID, entity.TierBronze, "9.99"),
						createPlan(silverPlanID, entity.TierSilver, "19.99"),
						createPlan(goldPlanID, entity.TierGold, "29.99"),
					}
					mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(plans, nil)
				}

				// Act
				result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

				// Assert
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectAccessible, result.Accessible)
				assert.Equal(t, tt.userTier, result.UserTier)
				assert.Equal(t, tt.requiredTier, result.RequiredTier)

				// Verify upgrade options when blocked
				if !tt.expectAccessible {
					assert.NotEmpty(t, result.UpgradeOptions)
				} else {
					assert.Empty(t, result.UpgradeOptions)
				}
			})
		}
	})

	t.Run("anonymous_user_gets_free_tier", func(t *testing.T) {
		// Setup: Blog requires BRONZE, user is anonymous
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierBronze, authorID, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return([]entity.SubscriptionPlan{
			createPlan(bronzePlanID, entity.TierBronze, "9.99"),
			createPlan(silverPlanID, entity.TierSilver, "19.99"),
			createPlan(goldPlanID, entity.TierGold, "29.99"),
		}, nil)

		// Act: Pass nil userID
		result, err := svc.CheckBlogAccess(ctx, blogID, nil)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, entity.TierFree, result.UserTier) // Anonymous users are FREE
		assert.Equal(t, entity.TierBronze, result.RequiredTier)
		assert.False(t, result.Accessible) // FREE < BRONZE
		assert.NotEmpty(t, result.UpgradeOptions)
	})

	t.Run("no_subscription_returns_free_tier", func(t *testing.T) {
		// Setup: User has no active subscription
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierSilver, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return([]entity.SubscriptionPlan{
			createPlan(silverPlanID, entity.TierSilver, "19.99"),
			createPlan(goldPlanID, entity.TierGold, "29.99"),
		}, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, result.UserTier)
		assert.Equal(t, entity.TierSilver, result.RequiredTier)
		assert.False(t, result.Accessible)
		assert.NotEmpty(t, result.UpgradeOptions)
	})

	t.Run("expired_subscription_treated_as_free", func(t *testing.T) {
		// Setup: User has expired subscription
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierGold, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return([]entity.SubscriptionPlan{
			createPlan(goldPlanID, entity.TierGold, "29.99"),
		}, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, result.UserTier) // No active subscription
		assert.Equal(t, entity.TierGold, result.RequiredTier)
		assert.False(t, result.Accessible)
	})

	t.Run("upgrade_options_generated_when_blocked", func(t *testing.T) {
		// Setup: FREE user trying to access GOLD content
		plans := []entity.SubscriptionPlan{
			createPlan(bronzePlanID, entity.TierBronze, "9.99"),
			createPlan(silverPlanID, entity.TierSilver, "19.99"),
			createPlan(goldPlanID, entity.TierGold, "29.99"),
		}

		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierGold, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(plans, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.False(t, result.Accessible)
		assert.NotEmpty(t, result.UpgradeOptions)
	})

	t.Run("no_upgrade_options_when_accessible", func(t *testing.T) {
		// Setup: GOLD user accessing GOLD content (accessible)
		sub := &entity.Subscription{
			SubscriberID: userID,
			AuthorID:     authorID,
			Tier:         string(entity.TierGold),
		}

		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierGold, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.True(t, result.Accessible)
		assert.Empty(t, result.UpgradeOptions) // No upgrade options when accessible
	})

	t.Run("error_getting_required_tier", func(t *testing.T) {
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierFree, uuid.Nil, errors.New("blog not found"))

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get required tier")
	})

	t.Run("error_getting_user_tier", func(t *testing.T) {
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierSilver, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, errors.New("database error"))

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get user tier")
	})

	t.Run("invalid_tier_returns_free", func(t *testing.T) {
		// Setup: Subscription has invalid tier value
		sub := &entity.Subscription{
			SubscriberID: userID,
			AuthorID:     authorID,
			Tier:         "INVALID_TIER",
		}

		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierBronze, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return([]entity.SubscriptionPlan{
			createPlan(bronzePlanID, entity.TierBronze, "9.99"),
		}, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, result.UserTier) // Invalid tier -> FREE
		assert.Equal(t, entity.TierBronze, result.RequiredTier)
		assert.False(t, result.Accessible)
		assert.NotEmpty(t, result.UpgradeOptions)
	})

	t.Run("no_upgrade_options_when_plan_fetch_fails", func(t *testing.T) {
		// Setup: User is blocked, but plan fetch fails
		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierSilver, authorID, nil)
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(nil, errors.New("database error"))

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, &userID)

		// Assert
		assert.NoError(t, err)
		assert.False(t, result.Accessible)
		assert.Empty(t, result.UpgradeOptions) // Empty on error
	})

	t.Run("upgrade_options_for_anonymous_user", func(t *testing.T) {
		// Setup: Anonymous user blocked
		plans := []entity.SubscriptionPlan{
			createPlan(bronzePlanID, entity.TierBronze, "9.99"),
			createPlan(silverPlanID, entity.TierSilver, "19.99"),
			createPlan(goldPlanID, entity.TierGold, "29.99"),
		}

		mockTagTierService.EXPECT().GetRequiredTierForBlog(ctx, blogID).Return(entity.TierSilver, authorID, nil)
		mockPlanRepo.EXPECT().FindActiveByAuthor(ctx, authorID).Return(plans, nil)

		// Act
		result, err := svc.CheckBlogAccess(ctx, blogID, nil)

		// Assert
		assert.NoError(t, err)
		assert.False(t, result.Accessible)
		assert.NotEmpty(t, result.UpgradeOptions)
		// For anonymous users, should show all tiers > FREE (all paid tiers)
	})
}

func TestContentAccessService_GetUserTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierService := serviceMocks.NewMockTagTierService(ctrl)
	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockPlanRepo := mocks.NewMockSubscriptionPlanRepository(ctrl)

	svc := service.NewContentAccessService(mockTagTierService, mockSubRepo, mockPlanRepo)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()

	t.Run("returns_free_tier_when_no_subscription", func(t *testing.T) {
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, nil)

		// Act
		tier, err := svc.GetUserTier(ctx, userID, authorID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, tier)
	})

	t.Run("returns_user_tier_from_active_subscription", func(t *testing.T) {
		sub := &entity.Subscription{
			SubscriberID: userID,
			AuthorID:     authorID,
			Tier:         string(entity.TierGold),
		}

		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)

		// Act
		tier, err := svc.GetUserTier(ctx, userID, authorID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierGold, tier)
	})

	t.Run("returns_free_on_invalid_tier_string", func(t *testing.T) {
		sub := &entity.Subscription{
			SubscriberID: userID,
			AuthorID:     authorID,
			Tier:         "INVALID",
		}

		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)

		// Act
		tier, err := svc.GetUserTier(ctx, userID, authorID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, tier)
	})

	t.Run("all_valid_tier_levels", func(t *testing.T) {
		tests := []struct {
			tierString string
			expected   entity.SubscriptionTier
		}{
			{"FREE", entity.TierFree},
			{"BRONZE", entity.TierBronze},
			{"SILVER", entity.TierSilver},
			{"GOLD", entity.TierGold},
		}

		for _, tt := range tests {
			t.Run(tt.tierString, func(t *testing.T) {
				sub := &entity.Subscription{
					SubscriberID: userID,
					AuthorID:     authorID,
					Tier:         tt.tierString,
				}

				mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(sub, nil)

				// Act
				tier, err := svc.GetUserTier(ctx, userID, authorID)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tier)
			})
		}
	})

	t.Run("error_on_repository_failure", func(t *testing.T) {
		mockSubRepo.EXPECT().FindActiveSubscription(ctx, userID, authorID).Return(nil, errors.New("database error"))

		// Act
		tier, err := svc.GetUserTier(ctx, userID, authorID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.TierFree, tier) // Returns FREE even on error
		assert.Contains(t, err.Error(), "failed to get user tier")
	})
}
