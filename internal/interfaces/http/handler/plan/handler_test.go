package plan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/service/mocks"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupRouter() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	w := httptest.NewRecorder()
	return r, w
}

func TestPlanHandler_UpsertPlans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/plans", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UpsertPlans(c)
		})

		reqBody := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{
					Tier:  "BRONZE",
					Price: decimal.NewFromInt(999),
				},
				{
					Tier:  "SILVER",
					Price: decimal.NewFromInt(1999),
				},
				{
					Tier:  "GOLD",
					Price: decimal.NewFromInt(2999),
				},
			},
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedPlans := []entity.SubscriptionPlan{
			{
				ID:           uuid.New(),
				Tier:         entity.TierBronze,
				Price:        decimal.NewFromInt(999),
				DurationDays: 30,
				IsActive:     true,
			},
		}
		expectedWarnings := []string{}

		mockPlanService.EXPECT().
			UpsertPlans(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedPlans, expectedWarnings, nil).Times(1)

		req, _ := http.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/plans", handler.UpsertPlans)

		reqBody := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "BRONZE", Price: decimal.NewFromInt(999)},
			},
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_request_tier_validation", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/plans", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UpsertPlans(c)
		})

		reqBody := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "INVALID", Price: decimal.NewFromInt(999)},
			},
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid_request_price_validation", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/plans", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UpsertPlans(c)
		})

		reqBody := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "BRONZE", Price: decimal.NewFromInt(-100)},
			},
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/plans", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UpsertPlans(c)
		})

		reqBody := dto.UpsertPlansRequest{
			Plans: []dto.CreatePlanRequest{
				{Tier: "BRONZE", Price: decimal.NewFromInt(999)},
			},
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockPlanService.EXPECT().
			UpsertPlans(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil, errors.New("database error")).Times(1)

		req, _ := http.NewRequest(http.MethodPost, "/plans", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPlanHandler_GetAuthorPlans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/authors/:authorId/plans", handler.GetAuthorPlans)

		authorID := uuid.New()

		expectedPlans := []service.PlanWithTags{
			{
				Plan: entity.SubscriptionPlan{
					ID:           uuid.New(),
					AuthorID:     authorID,
					Tier:         entity.TierBronze,
					Price:        decimal.NewFromInt(999),
					DurationDays: 30,
					IsActive:     true,
				},
				Tags:     []string{"tag1", "tag2"},
				TagCount: 2,
			},
		}

		mockPlanService.EXPECT().
			GetAuthorPlans(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, id uuid.UUID) ([]service.PlanWithTags, error) {
				assert.Equal(t, authorID, id)
				return expectedPlans, nil
			})

		req, _ := http.NewRequest(http.MethodGet, "/authors/"+authorID.String()+"/plans", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("invalid_author_id", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/authors/:authorId/plans", handler.GetAuthorPlans)

		req, _ := http.NewRequest(http.MethodGet, "/authors/invalid-id/plans", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/authors/:authorId/plans", handler.GetAuthorPlans)

		authorID := uuid.New()

		mockPlanService.EXPECT().
			GetAuthorPlans(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("database error"))

		req, _ := http.NewRequest(http.MethodGet, "/authors/"+authorID.String()+"/plans", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPlanHandler_AssignTagToTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.AssignTagToTier(c)
		})

		tagID := uuid.New()

		reqBody := dto.AssignTagTierRequest{
			RequiredTier: "BRONZE",
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedMapping := &entity.TagTierMapping{
			ID:           uuid.New(),
			AuthorID:     uuid.New(),
			TagID:        tagID,
			RequiredTier: entity.TierBronze,
		}

		mockTagService.EXPECT().
			AssignTagToTier(gomock.Any(), gomock.Any(), gomock.Any(), entity.TierBronze).
			DoAndReturn(func(_ context.Context, authorID, mappingTagID uuid.UUID, tier entity.SubscriptionTier) (*entity.TagTierMapping, int64, error) {
				assert.Equal(t, tagID, mappingTagID)
				assert.Equal(t, entity.TierBronze, tier)
				return expectedMapping, int64(5), nil
			})

		req, _ := http.NewRequest(http.MethodPost, "/tags/"+tagID.String()+"/tier", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/tags/:tagId/tier", handler.AssignTagToTier)

		tagID := uuid.New()

		reqBody := dto.AssignTagTierRequest{RequiredTier: "BRONZE"}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/tags/"+tagID.String()+"/tier", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_tag_id", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.AssignTagToTier(c)
		})

		reqBody := dto.AssignTagTierRequest{RequiredTier: "BRONZE"}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/tags/invalid-id/tier", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid_tier_validation", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.AssignTagToTier(c)
		})

		tagID := uuid.New()

		reqBody := dto.AssignTagTierRequest{RequiredTier: "INVALID"}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/tags/"+tagID.String()+"/tier", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.AssignTagToTier(c)
		})

		tagID := uuid.New()

		reqBody := dto.AssignTagTierRequest{RequiredTier: "BRONZE"}
		jsonBody, _ := json.Marshal(reqBody)

		mockTagService.EXPECT().
			AssignTagToTier(gomock.Any(), gomock.Any(), gomock.Any(), entity.TierBronze).
			Return(nil, int64(0), errors.New("database error"))

		req, _ := http.NewRequest(http.MethodPost, "/tags/"+tagID.String()+"/tier", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPlanHandler_UnassignTagFromTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.DELETE("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UnassignTagFromTier(c)
		})

		tagID := uuid.New()

		mockTagService.EXPECT().
			UnassignTagFromTier(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, authorID, mappingTagID uuid.UUID) (int64, error) {
				assert.Equal(t, tagID, mappingTagID)
				return int64(5), nil
			})

		req, _ := http.NewRequest(http.MethodDelete, "/tags/"+tagID.String()+"/tier", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		r, w := setupRouter()
		r.DELETE("/tags/:tagId/tier", handler.UnassignTagFromTier)

		tagID := uuid.New()

		req, _ := http.NewRequest(http.MethodDelete, "/tags/"+tagID.String()+"/tier", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid_tag_id", func(t *testing.T) {
		r, w := setupRouter()
		r.DELETE("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UnassignTagFromTier(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/tags/invalid-id/tier", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.DELETE("/tags/:tagId/tier", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.UnassignTagFromTier(c)
		})

		tagID := uuid.New()

		mockTagService.EXPECT().
			UnassignTagFromTier(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("database error"))

		req, _ := http.NewRequest(http.MethodDelete, "/tags/"+tagID.String()+"/tier", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPlanHandler_GetAuthorTagTiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/tag-tiers", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.GetAuthorTagTiers(c)
		})

		expectedMappings := []service.TagTierWithCount{
			{
				Mapping: entity.TagTierMapping{
					ID:           uuid.New(),
					AuthorID:     uuid.New(),
					TagID:        uuid.New(),
					RequiredTier: entity.TierBronze,
				},
				TagName:   "premium",
				BlogCount: int64(10),
			},
		}

		mockTagService.EXPECT().
			GetAuthorTagTiers(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, authorID uuid.UUID) ([]service.TagTierWithCount, error) {
				return expectedMappings, nil
			})

		req, _ := http.NewRequest(http.MethodGet, "/tag-tiers", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/tag-tiers", handler.GetAuthorTagTiers)

		req, _ := http.NewRequest(http.MethodGet, "/tag-tiers", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/tag-tiers", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.GetAuthorTagTiers(c)
		})

		mockTagService.EXPECT().
			GetAuthorTagTiers(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("database error"))

		req, _ := http.NewRequest(http.MethodGet, "/tag-tiers", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPlanHandler_CheckBlogAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlanService := mocks.NewMockPlanManagementService(ctrl)
	mockTagService := mocks.NewMockTagTierService(ctrl)
	mockAccessService := mocks.NewMockContentAccessService(ctrl)
	handler := NewPlanHandler(mockPlanService, mockTagService, mockAccessService)

	t.Run("success_authenticated_user_with_access", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/blogs/:blogId/access", func(c *gin.Context) {
			c.Set("userID", uuid.New())
			handler.CheckBlogAccess(c)
		})

		blogID := uuid.New()

		expectedResult := &service.AccessResult{
			Accessible:     true,
			UserTier:       entity.TierBronze,
			RequiredTier:   entity.TierBronze,
			Reason:         "User has access",
			UpgradeOptions: nil,
		}

		mockAccessService.EXPECT().
			CheckBlogAccess(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, id uuid.UUID, uid *uuid.UUID) (*service.AccessResult, error) {
				assert.Equal(t, blogID, id)
				return expectedResult, nil
			})

		req, _ := http.NewRequest(http.MethodGet, "/blogs/"+blogID.String()+"/access", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("success_anonymous_user_no_access", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/blogs/:blogId/access", handler.CheckBlogAccess)

		blogID := uuid.New()

		expectedResult := &service.AccessResult{
			Accessible:   false,
			UserTier:     entity.TierFree,
			RequiredTier: entity.TierSilver,
			Reason:       "Content requires SILVER tier or higher",
			UpgradeOptions: []service.UpgradeOption{
				{
					PlanID:       uuid.New(),
					Tier:         entity.TierSilver,
					Price:        "19.99",
					DurationDays: 30,
				},
			},
		}

		mockAccessService.EXPECT().
			CheckBlogAccess(gomock.Any(), gomock.Any(), (*uuid.UUID)(nil)).
			DoAndReturn(func(_ context.Context, id uuid.UUID, uid *uuid.UUID) (*service.AccessResult, error) {
				assert.Equal(t, blogID, id)
				return expectedResult, nil
			})

		req, _ := http.NewRequest(http.MethodGet, "/blogs/"+blogID.String()+"/access", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp.Success)
	})

	t.Run("invalid_blog_id", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/blogs/:blogId/access", handler.CheckBlogAccess)

		req, _ := http.NewRequest(http.MethodGet, "/blogs/invalid-id/access", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/blogs/:blogId/access", handler.CheckBlogAccess)

		blogID := uuid.New()

		mockAccessService.EXPECT().
			CheckBlogAccess(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("database error"))

		req, _ := http.NewRequest(http.MethodGet, "/blogs/"+blogID.String()+"/access", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
