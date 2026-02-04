package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	planHandler "github.com/aiagent/internal/interfaces/http/handler/plan"
	"github.com/aiagent/internal/interfaces/http/router"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestResponse wraps the standardized API response
type TestResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	} `json:"error,omitempty"`
}

// PlanTestContext holds test dependencies
type PlanTestContext struct {
	DB           *gorm.DB
	AuthorID     uuid.UUID
	Server       *httptest.Server
	SessionToken string
	Cleanup      func()
}

// conditionalAuthMiddleware creates a mock auth middleware that only sets userID if auth header is present
func conditionalAuthMiddleware(authorID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			c.Set("userID", authorID)
		}
		c.Next()
	}
}

// setupPlanTestServer creates a test server with all dependencies
func setupPlanTestServer(t *testing.T) *PlanTestContext {
	db, cleanup := setupTestDB(t)

	// Setup repositories
	planRepo := repository.NewSubscriptionPlanRepository(db)
	tagRepo := repository.NewTagRepository(db)
	tagTierRepo := repository.NewTagTierMappingRepository(db)
	blogRepo := repository.NewBlogRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)

	// Setup services
	planService := service.NewPlanManagementService(planRepo, tagTierRepo)
	tagService := service.NewTagTierService(tagTierRepo, tagRepo, blogRepo)
	accessService := service.NewContentAccessService(tagService, subRepo, planRepo)

	// Setup handler
	handler := planHandler.NewPlanHandler(planService, tagService, accessService)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	// Create mock auth middleware that requires auth header
	authorID := uuid.New()
	sessionAuth := conditionalAuthMiddleware(authorID)

	// Register routes
	router.RegisterPlanRoutes(v1, handler, sessionAuth)

	server := httptest.NewServer(r)

	// Create author in database
	userRepo := repository.NewUserRepository(db)
	author := &entity.User{
		ID:       authorID,
		Email:    "author@example.com",
		Name:     "Test Author",
		IsActive: true,
	}
	err := userRepo.Create(context.Background(), author)
	require.NoError(t, err)

	return &PlanTestContext{
		DB:           db,
		AuthorID:     authorID,
		Server:       server,
		SessionToken: "mock-session-token",
		Cleanup: func() {
			server.Close()
			cleanup()
		},
	}
}

// TestPlanAPI_UpsertPlans tests the POST /api/v1/authors/me/plans endpoint
func TestPlanAPI_UpsertPlans(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	t.Run("Success - Create BRONZE plan", func(t *testing.T) {
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{
					Tier:  "BRONZE",
					Price: decimal.NewFromInt(50000),
				},
			},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.UpsertPlansResponse
		json.Unmarshal(dataBytes, &result)

		assert.Len(t, result.Plans, 1)
		assert.Equal(t, "BRONZE", result.Plans[0].Tier)
		assert.Equal(t, 0, result.Plans[0].Price.Cmp(decimal.NewFromInt(50000)))
		assert.Equal(t, 30, result.Plans[0].DurationDays)
		assert.True(t, result.Plans[0].IsActive)
	})

	t.Run("Success - Create multiple plans (BRONZE, SILVER, GOLD)", func(t *testing.T) {
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{
					Tier:  "BRONZE",
					Price: decimal.NewFromInt(50000),
					Name:  strPtr("Basic Plan"),
				},
				{
					Tier:  "SILVER",
					Price: decimal.NewFromInt(100000),
					Name:  strPtr("Premium Plan"),
				},
				{
					Tier:  "GOLD",
					Price: decimal.NewFromInt(200000),
					Name:  strPtr("VIP Plan"),
				},
			},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.UpsertPlansResponse
		json.Unmarshal(dataBytes, &result)

		assert.Len(t, result.Plans, 3)
	})

	t.Run("Authentication - Reject unauthenticated request", func(t *testing.T) {
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "BRONZE", Price: decimal.NewFromInt(50000)},
			},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
	})

	t.Run("Validation - Invalid tier returns 400", func(t *testing.T) {
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "PLATINUM", Price: decimal.NewFromInt(50000)},
			},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
		assert.Contains(t, resp.Error.Message, "invalid tier")
	})

	t.Run("Validation - Negative price returns 400", func(t *testing.T) {
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "BRONZE", Price: decimal.NewFromInt(-100)},
			},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
		assert.Contains(t, resp.Error.Message, "price")
	})

	t.Run("Validation - Empty plans returns success with empty list", func(t *testing.T) {
		// Note: Current handler implementation accepts empty arrays and returns empty result
		payload := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{},
		}

		reqBody, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ctx.Server.URL+"/api/v1/authors/me/plans", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		// Currently returns 200 with empty plans list
		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.UpsertPlansResponse
		json.Unmarshal(dataBytes, &result)

		assert.Len(t, result.Plans, 0)
	})
}

// TestPlanAPI_GetAuthorPlans tests the GET /api/v1/authors/:authorId/plans endpoint
func TestPlanAPI_GetAuthorPlans(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	// Create some plans first
	planRepo := repository.NewSubscriptionPlanRepository(ctx.DB)

	plans := []entity.SubscriptionPlan{
		{
			ID:           uuid.New(),
			AuthorID:     ctx.AuthorID,
			Tier:         entity.TierBronze,
			Price:        decimal.NewFromInt(50000),
			DurationDays: 30,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			AuthorID:     ctx.AuthorID,
			Tier:         entity.TierSilver,
			Price:        decimal.NewFromInt(100000),
			DurationDays: 30,
			IsActive:     true,
		},
	}

	for _, plan := range plans {
		err := planRepo.Create(context.Background(), &plan)
		require.NoError(t, err)
	}

	t.Run("Success - Get author's plans", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/"+ctx.AuthorID.String()+"/plans", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.GetAuthorPlansResponse
		json.Unmarshal(dataBytes, &result)

		assert.Equal(t, ctx.AuthorID, result.AuthorID)
		// GetAuthorPlans returns all 4 tiers (FREE, BRONZE, SILVER, GOLD)
		// FREE tier always included with price 0 and isActive=false
		assert.Len(t, result.Plans, 4)

		// Verify created plans have correct values
		var bronzePlan, silverPlan *dto.PlanWithTagsResponse
		for i := range result.Plans {
			if result.Plans[i].Tier == "BRONZE" {
				p := result.Plans[i]
				bronzePlan = &p
			} else if result.Plans[i].Tier == "SILVER" {
				p := result.Plans[i]
				silverPlan = &p
			}
		}

		require.NotNil(t, bronzePlan)
		assert.True(t, bronzePlan.Price.Equal(decimal.NewFromInt(50000)))

		require.NotNil(t, silverPlan)
		assert.True(t, silverPlan.Price.Equal(decimal.NewFromInt(100000)))
	})

	t.Run("Success - Public endpoint (no auth required)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/"+ctx.AuthorID.String()+"/plans", nil)
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Validation - Invalid author ID returns 400", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/invalid-uuid/plans", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	})

	t.Run("Success - Non-existent author returns empty plans list", func(t *testing.T) {
		unknownAuthorID := uuid.New()
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/"+unknownAuthorID.String()+"/plans", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.GetAuthorPlansResponse
		json.Unmarshal(dataBytes, &result)

		// GetAuthorPlans returns all 4 tiers even for non-existent authors
		// with price 0 for each tier (no IsActive field in PlanWithTagsResponse)
		assert.Len(t, result.Plans, 4)
		assert.Equal(t, unknownAuthorID, result.AuthorID)

		// All plans should have zero price
		for _, plan := range result.Plans {
			assert.True(t, plan.Price.IsZero())
		}
	})
}

// TestPlanAPI_AssignTagToTier tests the POST /api/v1/authors/me/tags/:tagId/tier endpoint
func TestPlanAPI_AssignTagToTier(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	// Create a tag
	tagRepo := repository.NewTagRepository(ctx.DB)
	blogRepo := repository.NewBlogRepository(ctx.DB)

	tag := &entity.Tag{
		Name: "Premium",
		Slug: "premium",
	}
	err := tagRepo.Create(context.Background(), tag)
	require.NoError(t, err)

	// Create a blog with this tag
	blog := &entity.Blog{
		ID:         uuid.New(),
		AuthorID:   ctx.AuthorID,
		Title:      "Premium Content",
		Slug:       "premium-content",
		Content:    "This is premium content",
		Status:     entity.BlogStatusPublished,
		Visibility: entity.BlogVisibilityPublic,
	}
	err = blogRepo.Create(context.Background(), blog)
	require.NoError(t, err)

	err = blogRepo.AddTags(context.Background(), blog.ID, []uuid.UUID{tag.ID})
	require.NoError(t, err)

	t.Run("Success - Assign tag to BRONZE tier", func(t *testing.T) {
		payload := dto.AssignTagTierRequest{
			RequiredTier: "BRONZE",
		}

		reqBody, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/%s/tier", ctx.Server.URL, tag.ID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.AssignTagTierResponse
		json.Unmarshal(dataBytes, &result)

		assert.Equal(t, tag.ID, result.TagID)
		assert.Equal(t, "BRONZE", result.RequiredTier)
		assert.Equal(t, int64(1), result.AffectedBlogsCount)
	})

	t.Run("Authentication - Reject unauthenticated request", func(t *testing.T) {
		payload := dto.AssignTagTierRequest{
			RequiredTier: "BRONZE",
		}

		reqBody, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/%s/tier", ctx.Server.URL, tag.ID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
	})

	t.Run("Validation - Invalid tag ID returns 400", func(t *testing.T) {
		payload := dto.AssignTagTierRequest{
			RequiredTier: "BRONZE",
		}

		reqBody, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/invalid-uuid/tier", ctx.Server.URL)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	})

	t.Run("Validation - Invalid tier returns 400", func(t *testing.T) {
		payload := dto.AssignTagTierRequest{
			RequiredTier: "PLATINUM",
		}

		reqBody, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/%s/tier", ctx.Server.URL, tag.ID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
		assert.Contains(t, resp.Error.Message, "RequiredTier")
		assert.Contains(t, resp.Error.Message, "oneof")
	})
}

// TestPlanAPI_UnassignTagFromTier tests the DELETE /api/v1/authors/me/tags/:tagId/tier endpoint
func TestPlanAPI_UnassignTagFromTier(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	// Create a tag and assign it to a tier
	tagRepo := repository.NewTagRepository(ctx.DB)
	tagTierRepo := repository.NewTagTierMappingRepository(ctx.DB)
	blogRepo := repository.NewBlogRepository(ctx.DB)

	tag := &entity.Tag{
		Name: "Exclusive",
		Slug: "exclusive",
	}
	err := tagRepo.Create(context.Background(), tag)
	require.NoError(t, err)

	blog := &entity.Blog{
		ID:         uuid.New(),
		AuthorID:   ctx.AuthorID,
		Title:      "Exclusive Content",
		Slug:       "exclusive-content",
		Content:    "This is exclusive content",
		Status:     entity.BlogStatusPublished,
		Visibility: entity.BlogVisibilityPublic,
	}
	err = blogRepo.Create(context.Background(), blog)
	require.NoError(t, err)

	err = blogRepo.AddTags(context.Background(), blog.ID, []uuid.UUID{tag.ID})
	require.NoError(t, err)

	mapping := &entity.TagTierMapping{
		ID:           uuid.New(),
		AuthorID:     ctx.AuthorID,
		TagID:        tag.ID,
		RequiredTier: entity.TierSilver,
	}
	err = tagTierRepo.Create(context.Background(), mapping)
	require.NoError(t, err)

	t.Run("Success - Unassign tag from tier", func(t *testing.T) {
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/%s/tier", ctx.Server.URL, tag.ID)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.UnassignTagTierResponse
		json.Unmarshal(dataBytes, &result)

		assert.Equal(t, "Tag tier assignment removed successfully", result.Message)
		assert.Equal(t, int64(1), result.AffectedBlogsCount)
	})

	t.Run("Authentication - Reject unauthenticated request", func(t *testing.T) {
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/%s/tier", ctx.Server.URL, tag.ID)
		req, _ := http.NewRequest("DELETE", url, nil)
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
	})

	t.Run("Validation - Invalid tag ID returns 400", func(t *testing.T) {
		url := fmt.Sprintf("%s/api/v1/authors/me/tags/invalid-uuid/tier", ctx.Server.URL)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	})
}

// TestPlanAPI_GetAuthorTagTiers tests the GET /api/v1/authors/me/tag-tiers endpoint
func TestPlanAPI_GetAuthorTagTiers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	// Create some tag-tier mappings
	tagRepo := repository.NewTagRepository(ctx.DB)
	tagTierRepo := repository.NewTagTierMappingRepository(ctx.DB)

	tags := []*entity.Tag{
		{Name: "Premium", Slug: "premium"},
		{Name: "Exclusive", Slug: "exclusive"},
	}

	for _, tag := range tags {
		err := tagRepo.Create(context.Background(), tag)
		require.NoError(t, err)

		mapping := &entity.TagTierMapping{
			ID:           uuid.New(),
			AuthorID:     ctx.AuthorID,
			TagID:        tag.ID,
			RequiredTier: entity.TierSilver,
		}
		err = tagTierRepo.Create(context.Background(), mapping)
		require.NoError(t, err)
	}

	t.Run("Success - Get author's tag-tier mappings", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/me/tag-tiers", nil)
		req.Header.Set("Authorization", "Bearer "+ctx.SessionToken)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.GetTagTiersResponse
		json.Unmarshal(dataBytes, &result)

		assert.Len(t, result.Mappings, 2)
	})

	t.Run("Authentication - Reject unauthenticated request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/authors/me/tag-tiers", nil)
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
	})

	t.Run("Success - Empty list when no mappings exist for new author", func(t *testing.T) {
		// Create a new server for a different author with no mappings
		db2, cleanup2 := setupTestDB(t)
		defer cleanup2()

		newAuthorID := uuid.New()

		// Create repositories
		planRepo := repository.NewSubscriptionPlanRepository(db2)
		tagRepo := repository.NewTagRepository(db2)
		tagTierRepo := repository.NewTagTierMappingRepository(db2)
		blogRepo := repository.NewBlogRepository(db2)
		subRepo := repository.NewSubscriptionRepository(db2)

		// Create author
		userRepo := repository.NewUserRepository(db2)
		author := &entity.User{
			ID:       newAuthorID,
			Email:    "newauthor@example.com",
			Name:     "New Author",
			IsActive: true,
		}
		err := userRepo.Create(context.Background(), author)
		require.NoError(t, err)

		// Setup services
		planService := service.NewPlanManagementService(planRepo, tagTierRepo)
		tagService := service.NewTagTierService(tagTierRepo, tagRepo, blogRepo)
		accessService := service.NewContentAccessService(tagService, subRepo, planRepo)

		handler := planHandler.NewPlanHandler(planService, tagService, accessService)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		v1 := r.Group("/api/v1")

		sessionAuth := conditionalAuthMiddleware(newAuthorID)
		router.RegisterPlanRoutes(v1, handler, sessionAuth)

		req, _ := http.NewRequest("GET", "/api/v1/authors/me/tag-tiers", nil)
		req.Header.Set("Authorization", "Bearer mock-token")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.GetTagTiersResponse
		json.Unmarshal(dataBytes, &result)

		assert.Len(t, result.Mappings, 0)
	})
}

// TestPlanAPI_CheckBlogAccess tests the GET /api/v1/blogs/:blogId/access endpoint
func TestPlanAPI_CheckBlogAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupPlanTestServer(t)
	defer ctx.Cleanup()

	// Create a blog without tags (public content)
	blogRepo := repository.NewBlogRepository(ctx.DB)
	publicBlog := &entity.Blog{
		ID:         uuid.New(),
		AuthorID:   ctx.AuthorID,
		Title:      "Public Blog",
		Slug:       "public-blog",
		Content:    "This is public content",
		Status:     entity.BlogStatusPublished,
		Visibility: entity.BlogVisibilityPublic,
	}
	err := blogRepo.Create(context.Background(), publicBlog)
	require.NoError(t, err)

	// Create a blog with a premium tag
	tagRepo := repository.NewTagRepository(ctx.DB)
	tag := &entity.Tag{
		Name: "Premium",
		Slug: "premium",
	}
	err = tagRepo.Create(context.Background(), tag)
	require.NoError(t, err)

	premiumBlog := &entity.Blog{
		ID:         uuid.New(),
		AuthorID:   ctx.AuthorID,
		Title:      "Premium Blog",
		Slug:       "premium-blog",
		Content:    "This is premium content",
		Status:     entity.BlogStatusPublished,
		Visibility: entity.BlogVisibilityPublic,
	}
	err = blogRepo.Create(context.Background(), premiumBlog)
	require.NoError(t, err)

	err = blogRepo.AddTags(context.Background(), premiumBlog.ID, []uuid.UUID{tag.ID})
	require.NoError(t, err)

	// Assign the tag to SILVER tier
	tagTierRepo := repository.NewTagTierMappingRepository(ctx.DB)
	mapping := &entity.TagTierMapping{
		ID:           uuid.New(),
		AuthorID:     ctx.AuthorID,
		TagID:        tag.ID,
		RequiredTier: entity.TierSilver,
	}
	err = tagTierRepo.Create(context.Background(), mapping)
	require.NoError(t, err)

	// Create subscription plans for upgrade options
	planRepo := repository.NewSubscriptionPlanRepository(ctx.DB)
	silverPlan := &entity.SubscriptionPlan{
		ID:           uuid.New(),
		AuthorID:     ctx.AuthorID,
		Tier:         entity.TierSilver,
		Price:        decimal.NewFromInt(100000),
		DurationDays: 30,
		IsActive:     true,
	}
	err = planRepo.Create(context.Background(), silverPlan)
	require.NoError(t, err)

	t.Run("Success - Anonymous user can access public content", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/blogs/"+publicBlog.ID.String()+"/access", nil)
		// No authorization header

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.CheckBlogAccessResponse
		json.Unmarshal(dataBytes, &result)

		assert.True(t, result.Accessible)
		assert.Equal(t, "FREE", result.UserTier)
	})

	t.Run("Success - Public endpoint (no auth required)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/blogs/"+publicBlog.ID.String()+"/access", nil)
		// No authorization header - should still work

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success - Anonymous user blocked from premium content with upgrade options", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/blogs/"+premiumBlog.ID.String()+"/access", nil)
		// No authorization header - anonymous user

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		dataBytes, _ := json.Marshal(resp.Data)
		var result dto.CheckBlogAccessResponse
		json.Unmarshal(dataBytes, &result)

		assert.False(t, result.Accessible)
		assert.Equal(t, "FREE", result.UserTier)
		assert.Equal(t, "SILVER", result.RequiredTier)
		assert.NotEmpty(t, result.UpgradeOptions)
		// Verify SILVER plan is in upgrade options
		hasSilverUpgrade := false
		for _, opt := range result.UpgradeOptions {
			if opt.Tier == "SILVER" {
				hasSilverUpgrade = true
				break
			}
		}
		assert.True(t, hasSilverUpgrade, "SILVER upgrade option should be available")
	})

	t.Run("Validation - Invalid blog ID returns 400", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/blogs/invalid-uuid/access", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	})

	t.Run("Validation - Non-existent blog ID returns 500 (service error)", func(t *testing.T) {
		// When a blog doesn't exist, the service returns an error
		unknownBlogID := uuid.New()
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/blogs/"+unknownBlogID.String()+"/access", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		// Returns 500 due to service error
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp TestResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
	})
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}
