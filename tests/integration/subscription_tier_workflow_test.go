package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/payment"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/adapter"
	"github.com/aiagent/internal/infrastructure/config"
	pgRepo "github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionTierWorkflow tests the complete subscription tier workflow
// from plan creation through purchase to access check
func TestSubscriptionTierWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Setup repositories
	userRepo := pgRepo.NewUserRepository(db)
	subRepo := pgRepo.NewSubscriptionRepository(db)
	planRepo := pgRepo.NewSubscriptionPlanRepository(db)
	tagRepo := pgRepo.NewTagRepository(db)
	tagTierRepo := pgRepo.NewTagTierMappingRepository(db)
	blogRepo := pgRepo.NewBlogRepository(db)

	// Setup services
	planMgmtSvc := service.NewPlanManagementService(planRepo, tagTierRepo)
	tagTierSvc := service.NewTagTierService(tagTierRepo, tagRepo, blogRepo)
	contentAccessSvc := service.NewContentAccessService(tagTierSvc, subRepo, planRepo)
	subscriptionSvc := service.NewSubscriptionService(subRepo)

	// Setup payment service (for simulating webhooks)
	txRepo := pgRepo.NewTransactionRepository(db)
	cfg := &config.SePayConfig{
		APIKey:       "test-api-key",
		BankName:     "Test Bank",
		BankAccount:  "123456789",
		BankOwner:    "TEST OWNER",
		BankBranch:   "Test Branch",
		WebhookToken: "test-webhook-token",
	}
	sepayAdapter := adapter.NewSePayAdapter(cfg)
	paymentSvc := service.NewPaymentService(db, txRepo, subRepo, pgRepo.NewUserSeriesPurchaseRepository(db), planRepo, sepayAdapter)
	createPaymentUC := payment.NewCreatePaymentUseCase(paymentSvc)
	processWebhookUC := payment.NewProcessWebhookUseCase(paymentSvc)

	// Create test users
	authorID := uuid.New()
	readerID := uuid.New()

	author := &entity.User{
		ID:       authorID,
		Email:    "author@example.com",
		Name:     "Test Author",
		IsActive: true,
	}
	err := userRepo.Create(ctx, author)
	require.NoError(t, err)

	reader := &entity.User{
		ID:       readerID,
		Email:    "reader@example.com",
		Name:     "Test Reader",
		IsActive: true,
	}
	err = userRepo.Create(ctx, reader)
	require.NoError(t, err)

	// Test 1: Author creates BRONZE plan → Reader purchases → Verify subscription
	t.Run("Author creates BRONZE plan and Reader purchases", func(t *testing.T) {
		// Arrange: Author creates BRONZE plan
		plansDTO := []service.CreatePlanDTO{
			{
				Tier:         entity.TierBronze,
				Price:        decimal.NewFromInt(50000),
				DurationDays: 30,
			},
		}

		createdPlans, warnings, err := planMgmtSvc.UpsertPlans(ctx, authorID, plansDTO)
		require.NoError(t, err)
		assert.Empty(t, warnings)
		assert.Len(t, createdPlans, 1)
		assert.Equal(t, entity.TierBronze, createdPlans[0].Tier)

		bronzePlanID := createdPlans[0].ID

		// Act: Reader creates a FREE subscription first
		subscription, err := subscriptionSvc.Subscribe(ctx, readerID, authorID)
		require.NoError(t, err)
		assert.Equal(t, "FREE", subscription.Tier)
		assert.Nil(t, subscription.ExpiresAt)

		// Act: Simulate BRONZE purchase via webhook
		targetID := authorID.String()
		planIDStr := bronzePlanID.String()

		paymentReq := dto.CreatePaymentRequest{
			UserID:   readerID.String(),
			Amount:   decimal.NewFromInt(50000),
			Type:     entity.TransactionTypeSubscription,
			Gateway:  entity.TransactionGatewayBankTransfer,
			TargetID: &targetID,
			PlanID:   &planIDStr,
		}

		paymentResp, err := createPaymentUC.Execute(ctx, paymentReq)
		require.NoError(t, err)
		assert.NotEmpty(t, paymentResp.ReferenceCode)

		// Simulate webhook
		webhookPayload := dto.ProcessWebhookRequest{
			ID:              int64(100001),
			Gateway:         "TestBank",
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			AccountNumber:   "123456789",
			Code:            "CODE001",
			Content:         paymentResp.ReferenceCode,
			TransferType:    "in",
			TransferAmount:  paymentReq.Amount,
			Accumulated:     paymentReq.Amount,
			SubAccount:      "",
			ReferenceCode:   "REF001",
			Description:     "Payment for subscription",
		}

		tx, err := processWebhookUC.Execute(ctx, webhookPayload)
		require.NoError(t, err)
		assert.NotNil(t, tx)

		// Assert: Verify subscription updated to BRONZE with expiry
		updatedSub, err := subRepo.FindActiveSubscription(ctx, readerID, authorID)
		require.NoError(t, err)
		require.NotNil(t, updatedSub)
		assert.Equal(t, entity.TierBronze.String(), updatedSub.Tier)
		assert.NotNil(t, updatedSub.ExpiresAt)
		assert.True(t, updatedSub.ExpiresAt.After(time.Now()))

		// Assert: Verify expiry is roughly 30 days from now
		// Note: Use large variance to account for potential timezone differences
		expectedExpiry := time.Now().AddDate(0, 0, 30)
		allowableVariance := 24 * time.Hour // 24 hours tolerance for timezone differences
		assert.WithinDuration(t, expectedExpiry, *updatedSub.ExpiresAt, allowableVariance)
	})

	// Test 2: Author assigns tag to SILVER → Blog with tag requires SILVER → BRONZE user blocked
	t.Run("Tag-tier mapping blocks lower tier users", func(t *testing.T) {
		// Arrange: Create SILVER plan
		silverPlanDTO := []service.CreatePlanDTO{
			{
				Tier:         entity.TierSilver,
				Price:        decimal.NewFromInt(100000),
				DurationDays: 30,
			},
		}

		silverPlans, warnings, err := planMgmtSvc.UpsertPlans(ctx, authorID, silverPlanDTO)
		require.NoError(t, err)
		assert.Empty(t, warnings)
		require.Len(t, silverPlans, 1)
		silverPlanID := silverPlans[0].ID

		// Arrange: Create and assign a tag to SILVER tier
		tag := &entity.Tag{
			Name: "Premium",
			Slug: "premium",
		}
		err = tagRepo.Create(ctx, tag)
		require.NoError(t, err)

		// Create a blog with the tag first
		blog := &entity.Blog{
			ID:         uuid.New(),
			AuthorID:   authorID,
			Title:      "Premium Content",
			Slug:       "premium-content",
			Content:    "This is premium content",
			Status:     entity.BlogStatusPublished,
			Visibility: entity.BlogVisibilityPublic,
		}
		err = blogRepo.Create(ctx, blog)
		require.NoError(t, err)

		// Assign the tag to the blog
		err = blogRepo.AddTags(ctx, blog.ID, []uuid.UUID{tag.ID})
		require.NoError(t, err)

		// Now assign the tag to SILVER tier
		mapping, blogCount, err := tagTierSvc.AssignTagToTier(ctx, authorID, tag.ID, entity.TierSilver)
		require.NoError(t, err)
		assert.NotNil(t, mapping)
		assert.Equal(t, int64(1), blogCount) // One blog with the tag
		assert.Equal(t, entity.TierSilver, mapping.RequiredTier)

		// Act: Check access with BRONZE tier
		accessResult, err := contentAccessSvc.CheckBlogAccess(ctx, blog.ID, &readerID)
		require.NoError(t, err)

		// Assert: BRONZE user should be blocked
		assert.False(t, accessResult.Accessible)
		assert.Equal(t, entity.TierBronze, accessResult.UserTier)
		assert.Equal(t, entity.TierSilver, accessResult.RequiredTier)
		assert.NotEmpty(t, accessResult.UpgradeOptions)

		// Assert: SILVER plan should be in upgrade options
		hasSilverOption := false
		for _, opt := range accessResult.UpgradeOptions {
			if opt.Tier == entity.TierSilver {
				hasSilverOption = true
				assert.Equal(t, silverPlanID, opt.PlanID)
			}
		}
		assert.True(t, hasSilverOption, "SILVER upgrade option should be available")
	})

	// Test 3: Reader upgrades BRONZE→SILVER → ExpiresAt resets, old tier replaced
	t.Run("Reader upgrades from BRONZE to SILVER", func(t *testing.T) {
		// Arrange: Get current BRONZE subscription expiry
		currentSub, err := subRepo.FindActiveSubscription(ctx, readerID, authorID)
		require.NoError(t, err)
		require.NotNil(t, currentSub)
		require.NotNil(t, currentSub.ExpiresAt)
		originalExpiry := *currentSub.ExpiresAt

		// Wait a moment to ensure time difference
		time.Sleep(100 * time.Millisecond)

		// Act: Get SILVER plan ID
		silverPlans, err := planRepo.FindActiveByAuthor(ctx, authorID)
		require.NoError(t, err)
		var silverPlan *entity.SubscriptionPlan
		for _, p := range silverPlans {
			if p.Tier == entity.TierSilver {
				silverPlan = &p
				break
			}
		}
		require.NotNil(t, silverPlan, "SILVER plan not found")

		// Simulate SILVER upgrade via webhook
		targetID := authorID.String()
		planIDStr := silverPlan.ID.String()

		paymentReq := dto.CreatePaymentRequest{
			UserID:   readerID.String(),
			Amount:   silverPlan.Price,
			Type:     entity.TransactionTypeSubscription,
			Gateway:  entity.TransactionGatewayBankTransfer,
			TargetID: &targetID,
			PlanID:   &planIDStr,
		}

		paymentResp, err := createPaymentUC.Execute(ctx, paymentReq)
		require.NoError(t, err)
		assert.NotEmpty(t, paymentResp.ReferenceCode)

		webhookPayload := dto.ProcessWebhookRequest{
			ID:              int64(100002),
			Gateway:         "TestBank",
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			AccountNumber:   "123456789",
			Code:            "CODE002",
			Content:         paymentResp.ReferenceCode,
			TransferType:    "in",
			TransferAmount:  paymentReq.Amount,
			Accumulated:     paymentReq.Amount,
			SubAccount:      "",
			ReferenceCode:   "REF002",
			Description:     "Upgrade to SILVER",
		}

		tx, err := processWebhookUC.Execute(ctx, webhookPayload)
		require.NoError(t, err)
		assert.NotNil(t, tx)

		// Assert: Verify subscription upgraded to SILVER with new expiry
		updatedSub, err := subRepo.FindActiveSubscription(ctx, readerID, authorID)
		require.NoError(t, err)
		require.NotNil(t, updatedSub)
		assert.Equal(t, entity.TierSilver.String(), updatedSub.Tier)
		assert.NotNil(t, updatedSub.ExpiresAt)

		// Assert: Expiry should be reset (new from now, not extended from old)
		newExpiry := *updatedSub.ExpiresAt
		assert.True(t, newExpiry.After(originalExpiry), "New expiry should be after original expiry for upgrade")

		// Verify new expiry is roughly 30 days from now
		expectedExpiry := time.Now().AddDate(0, 0, 30)
		allowableVariance := 24 * time.Hour // 24 hours tolerance for timezone differences
		assert.WithinDuration(t, expectedExpiry, newExpiry, allowableVariance, "New expiry should be 30 days from now")
	})

	// Test 4: Anonymous user checks access → returns FREE tier
	t.Run("Anonymous user gets FREE tier and is blocked from premium content", func(t *testing.T) {
		// Arrange: Use the premium blog from Test 2
		// Find a blog that requires SILVER tier
		_, _ = planRepo.FindActiveByAuthor(ctx, authorID)
		silverTagMaps, _ := tagTierRepo.FindByAuthor(ctx, authorID)
		var requiredBlogID uuid.UUID
		if len(silverTagMaps) > 0 {
			tagID := silverTagMaps[0].TagID
			// Find a blog with this tag
			blogs, err := blogRepo.FindAll(ctx, repository.BlogFilter{
				AuthorID: &authorID,
			}, repository.Pagination{Page: 1, PageSize: 10})
			require.NoError(t, err)
			if len(blogs.Data) > 0 {
				// Get the blog with tags
				blogWithTags, err := blogRepo.FindByID(ctx, blogs.Data[0].ID)
				require.NoError(t, err)
				if len(blogWithTags.Tags) > 0 && blogWithTags.Tags[0].ID == tagID {
					requiredBlogID = blogWithTags.ID
				}
			}
		}

		if requiredBlogID == uuid.Nil {
			t.Skip("No premium blog found for testing")
		}

		// Act: Check access with nil user ID (anonymous)
		accessResult, err := contentAccessSvc.CheckBlogAccess(ctx, requiredBlogID, nil)
		require.NoError(t, err)

		// Assert: Anonymous user should have FREE tier
		assert.Equal(t, entity.TierFree, accessResult.UserTier)
		assert.False(t, accessResult.Accessible, "Anonymous user should not have access to premium content")
		assert.NotEmpty(t, accessResult.UpgradeOptions)
	})

	// Test 5: Price hierarchy warning returned (not error)
	t.Run("Price hierarchy warning returned for invalid prices", func(t *testing.T) {
		// Arrange: Create plans with SILVER <= BRONZE price
		invalidPlans := []service.CreatePlanDTO{
			{
				Tier:         entity.TierBronze,
				Price:        decimal.NewFromInt(100000),
				DurationDays: 30,
			},
			{
				Tier:         entity.TierSilver,
				Price:        decimal.NewFromInt(100000), // Same price - should warn
				DurationDays: 30,
			},
		}

		// Act: Upsert plans
		createdPlans, warnings, err := planMgmtSvc.UpsertPlans(ctx, authorID, invalidPlans)

		// Assert: Should succeed but with warnings
		require.NoError(t, err, "Should not error on price hierarchy violation")
		assert.NotEmpty(t, warnings, "Should return warnings for price hierarchy violation")
		assert.Len(t, createdPlans, 2)

		// Assert: Warning message should mention SILVER and BRONZE
		hasWarning := false
		for _, warning := range warnings {
			if (strings.Contains(warning, "SILVER") || strings.Contains(warning, entity.TierSilver.String())) &&
				(strings.Contains(warning, "BRONZE") || strings.Contains(warning, entity.TierBronze.String())) {
				hasWarning = true
				break
			}
		}
		assert.True(t, hasWarning, "Warning should mention SILVER and BRONZE tiers")
	})

	// Test 6: Tag-tier mapping affects blog count correctly
	t.Run("Tag-tier mapping shows correct blog count", func(t *testing.T) {
		// Arrange: Create additional tags and blogs
		tag2 := &entity.Tag{
			Name: "Exclusive",
			Slug: "exclusive",
		}
		err := tagRepo.Create(ctx, tag2)
		require.NoError(t, err)

		// Create multiple blogs with the tag
		blogIDs := make([]uuid.UUID, 3)
		for i := 0; i < 3; i++ {
			blog := &entity.Blog{
				ID:         uuid.New(),
				AuthorID:   authorID,
				Title:      "Exclusive Content " + string(rune('1'+i)),
				Slug:       "exclusive-content-" + string(rune('1'+i)),
				Content:    "This is exclusive content",
				Status:     entity.BlogStatusPublished,
				Visibility: entity.BlogVisibilityPublic,
			}
			err := blogRepo.Create(ctx, blog)
			require.NoError(t, err)
			blogIDs[i] = blog.ID

			// Add the tag
			err = blogRepo.AddTags(ctx, blog.ID, []uuid.UUID{tag2.ID})
			require.NoError(t, err)
		}

		// Assign tag2 to GOLD tier
		goldPlanDTO := []service.CreatePlanDTO{
			{
				Tier:         entity.TierGold,
				Price:        decimal.NewFromInt(200000),
				DurationDays: 30,
			},
		}
		_, _, err = planMgmtSvc.UpsertPlans(ctx, authorID, goldPlanDTO)
		require.NoError(t, err)

		mapping, blogCount, err := tagTierSvc.AssignTagToTier(ctx, authorID, tag2.ID, entity.TierGold)
		require.NoError(t, err)
		assert.NotNil(t, mapping)
		assert.Equal(t, int64(3), blogCount, "Should count 3 blogs with the tag")

		// Act: Get author tag tiers
		tagTiers, err := tagTierSvc.GetAuthorTagTiers(ctx, authorID)
		require.NoError(t, err)

		// Assert: Find the tag2 mapping
		var tag2WithCount *service.TagTierWithCount
		for i := range tagTiers {
			if tagTiers[i].TagName == "Exclusive" {
				tag2WithCount = &tagTiers[i]
				break
			}
		}

		require.NotNil(t, tag2WithCount)
		assert.Equal(t, int64(3), tag2WithCount.BlogCount, "Blog count should be 3")
		assert.Equal(t, entity.TierGold, tag2WithCount.Mapping.RequiredTier)
	})

	// Cleanup: All data should be cleaned up by the test DB cleanup function
}
