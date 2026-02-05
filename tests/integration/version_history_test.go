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
	"github.com/aiagent/internal/application/usecase/blog"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	postgresRepository "github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	blogHandler "github.com/aiagent/internal/interfaces/http/handler/blog"
	versionHandler "github.com/aiagent/internal/interfaces/http/handler/version"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type VersionTestContext struct {
	DB             *gin.Engine
	BlogHandler    blogHandler.BlogHandler
	VersionHandler *versionHandler.VersionHandler
	TagRepo        repository.TagRepository
	UserAReader    *TestUser
	UserBReader    *TestUser
	Cleanup        func()
}

type TestUser struct {
	ID    uuid.UUID
	Email string
	Token string
}

func setupVersionTestServer(t *testing.T) *VersionTestContext {
	db, cleanup := setupTestDB(t)

	// Repositories
	blogRepo := postgresRepository.NewBlogRepository(db)
	versionRepo := postgresRepository.NewBlogVersionRepository(db)
	userRepo := postgresRepository.NewUserRepository(db)
	subRepo := postgresRepository.NewSubscriptionRepository(db)
	tagRepo := postgresRepository.NewTagRepository(db)

	// Services
	versionSvc := service.NewVersionService(versionRepo, blogRepo)
	// blogService needs Redis for ReactionBatcher.
	// We might need to mock Redis or use a real one if available.
	// Looking at payment_test.go, it doesn't seem to use Redis directly in setup.
	// Wait, blogService constructor requires Redis.
	// I'll use nil for Redis if it allows it, or I'll see how other tests handle it.
	// Actually, I'll use a nil redis for now and see if it crashes.
	blogSvc := service.NewBlogService(blogRepo, subRepo, tagRepo, nil, versionSvc)

	// UseCases
	blogUC := blog.NewBlogUseCase(blogSvc)

	// Handlers
	bHandler := blogHandler.NewBlogHandler(blogUC)
	vHandler := versionHandler.NewVersionHandler(versionSvc, blogSvc)

	// Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Mock Auth
	v1 := r.Group("/api/v1")

	// Users
	userA := &entity.User{
		ID:       uuid.New(),
		Email:    "usera@example.com",
		Name:     "User A",
		IsActive: true,
	}
	userB := &entity.User{
		ID:       uuid.New(),
		Email:    "userb@example.com",
		Name:     "User B",
		IsActive: true,
	}
	require.NoError(t, userRepo.Create(context.Background(), userA))
	require.NoError(t, userRepo.Create(context.Background(), userB))

	// Dynamic auth middleware that uses Authorization header to decide which user is logged in
	authMiddleware := func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "Bearer userA" {
			c.Set("userID", userA.ID)
		} else if token == "Bearer userB" {
			c.Set("userID", userB.ID)
		}
		c.Next()
	}

	// Register routes manually to avoid missing handler panics
	blogs := v1.Group("/blogs")
	{
		blogs.POST("", authMiddleware, bHandler.Create)
		blogs.PUT("/:id", authMiddleware, bHandler.Update)
	}

	versions := v1.Group("/blogs/:id/versions")
	versions.Use(authMiddleware)
	{
		versions.GET("", vHandler.List)
		versions.GET("/:versionId", vHandler.Get)
		versions.POST("", vHandler.Create)
		versions.POST("/:versionId/restore", vHandler.Restore)
		versions.DELETE("/:versionId", vHandler.Delete)
	}

	return &VersionTestContext{
		DB:             r,
		BlogHandler:    bHandler,
		VersionHandler: vHandler,
		TagRepo:        tagRepo,
		UserAReader: &TestUser{
			ID:    userA.ID,
			Email: userA.Email,
			Token: "Bearer userA",
		},
		UserBReader: &TestUser{
			ID:    userB.ID,
			Email: userB.Email,
			Token: "Bearer userB",
		},
		Cleanup: cleanup,
	}
}

func TestVersionHistoryWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupVersionTestServer(t)
	defer ctx.Cleanup()

	var blogID string
	var tag1ID uuid.UUID

	t.Run("0. Create a tag", func(t *testing.T) {
		tag := &entity.Tag{
			ID:   uuid.New(),
			Name: "Tag 1",
			Slug: "tag-1",
		}
		require.NoError(t, ctx.TagRepo.Create(context.Background(), tag))
		tag1ID = tag.ID
	})

	t.Run("1. Create blog (verify Version 1 exists and has tags)", func(t *testing.T) {
		reqBody := dto.CreateBlogRequest{
			Title:   "Initial Title",
			Slug:    "initial-slug",
			Content: "Initial Content",
			TagIDs:  []string{tag1ID.String()},
		}
		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/blogs", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp struct {
			Success bool             `json:"success"`
			Data    dto.BlogResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		blogID = resp.Data.ID.String()

		// Verify Version 1 exists
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var vResp struct {
			Success bool                  `json:"success"`
			Data    []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Len(t, vResp.Data, 1)
		assert.Equal(t, 1, vResp.Data[0].VersionNumber)

		// Verify Version 1 detail has tags
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions/%s", blogID, vResp.Data[0].ID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var vDetail struct {
			Data dto.VersionDetailResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vDetail)
		assert.Len(t, vDetail.Data.Tags, 1)
		assert.Equal(t, "Tag 1", vDetail.Data.Tags[0].Name)
	})

	t.Run("2. Update blog (verify Version 2 exists)", func(t *testing.T) {
		title := "Updated Title"
		reqBody := dto.UpdateBlogRequest{
			Title: &title,
		}
		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/blogs/%s", blogID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Version 2 exists
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		var vResp struct {
			Success bool                  `json:"success"`
			Data    []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Len(t, vResp.Data, 2)
		// Usually ordered DESC
		assert.Equal(t, 2, vResp.Data[0].VersionNumber)
	})

	t.Run("3. Create manual version (verify Version 3 exists)", func(t *testing.T) {
		reqBody := dto.CreateVersionRequest{
			ChangeSummary: "Manual checkpoint",
		}
		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify Version 3 exists
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		var vResp struct {
			Success bool                  `json:"success"`
			Data    []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Len(t, vResp.Data, 3)
		assert.Equal(t, 3, vResp.Data[0].VersionNumber)
	})

	t.Run("4. List versions (verify order and count)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions?pageSize=2", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var vResp struct {
			Success bool                  `json:"success"`
			Data    []dto.VersionResponse `json:"data"`
			Meta    struct {
				Total int64 `json:"total"`
			} `json:"meta"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Len(t, vResp.Data, 2)
		assert.Equal(t, int64(3), vResp.Meta.Total)
		assert.Equal(t, 3, vResp.Data[0].VersionNumber)
		assert.Equal(t, 2, vResp.Data[1].VersionNumber)
	})

	t.Run("5. Get version detail", func(t *testing.T) {
		// First get the ID of version 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var listResp struct {
			Data []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &listResp)
		var v1ID uuid.UUID
		for _, v := range listResp.Data {
			if v.VersionNumber == 1 {
				v1ID = v.ID
			}
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions/%s", blogID, v1ID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var vDetail struct {
			Data dto.VersionDetailResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vDetail)
		assert.Equal(t, "Initial Title", vDetail.Data.Title)
		assert.Equal(t, "Initial Content", vDetail.Data.Content)
	})

	t.Run("6. Restore to Version 1 (verify blog content, tags and Version 4 exists)", func(t *testing.T) {
		// Get ID of version 1
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var listResp struct {
			Data []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &listResp)
		var v1ID uuid.UUID
		for _, v := range listResp.Data {
			if v.VersionNumber == 1 {
				v1ID = v.ID
			}
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", fmt.Sprintf("/api/v1/blogs/%s/versions/%s/restore", blogID, v1ID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var bResp struct {
			Data dto.BlogResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &bResp)
		assert.Equal(t, "Initial Title", bResp.Data.Title)
		assert.Len(t, bResp.Data.Tags, 1)
		assert.Equal(t, "Tag 1", bResp.Data.Tags[0].Name)

		// Verify Version 4 exists
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var vResp struct {
			Data []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Len(t, vResp.Data, 4)
		assert.Equal(t, 4, vResp.Data[0].VersionNumber)
	})

	t.Run("7. Delete a version", func(t *testing.T) {
		// Get ID of version 2
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var listResp struct {
			Data []dto.VersionResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &listResp)
		var v2ID uuid.UUID
		for _, v := range listResp.Data {
			if v.VersionNumber == 2 {
				v2ID = v.ID
			}
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/blogs/%s/versions/%s", blogID, v2ID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify version count is 3
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		var vResp struct {
			Meta struct {
				Total int64 `json:"total"`
			} `json:"meta"`
		}
		json.Unmarshal(w.Body.Bytes(), &vResp)
		assert.Equal(t, int64(3), vResp.Meta.Total)
	})
}

func TestVersionHistorySecurity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupVersionTestServer(t)
	defer ctx.Cleanup()

	// 1. User A creates a blog
	reqBody := dto.CreateBlogRequest{
		Title:   "User A's Blog",
		Slug:    "usera-blog",
		Content: "Content A",
	}
	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.UserAReader.Token)
	ctx.DB.ServeHTTP(w, req)
	var resp struct {
		Data dto.BlogResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	blogID := resp.Data.ID

	// Get Version 1 ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
	req.Header.Set("Authorization", ctx.UserAReader.Token)
	ctx.DB.ServeHTTP(w, req)
	var listResp struct {
		Data []dto.VersionResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	v1ID := listResp.Data[0].ID

	t.Run("User B cannot list versions of User A's blog", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
		req.Header.Set("Authorization", ctx.UserBReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("User B cannot get version detail of User A's blog", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions/%s", blogID, v1ID), nil)
		req.Header.Set("Authorization", ctx.UserBReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("User B cannot restore version of User A's blog", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/blogs/%s/versions/%s/restore", blogID, v1ID), nil)
		req.Header.Set("Authorization", ctx.UserBReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("User B cannot create version for User A's blog", func(t *testing.T) {
		reqBody := dto.CreateVersionRequest{
			ChangeSummary: "User B trying to create version",
		}
		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.UserBReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("User B cannot delete User A's version", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/blogs/%s/versions/%s", blogID, v1ID), nil)
		req.Header.Set("Authorization", ctx.UserBReader.Token)
		ctx.DB.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestVersionHistoryLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupVersionTestServer(t)
	defer ctx.Cleanup()

	// 1. User A creates a blog
	reqBody := dto.CreateBlogRequest{
		Title:   "Limit Test Blog",
		Slug:    "limit-test",
		Content: "Initial",
	}
	jsonBody, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/blogs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.UserAReader.Token)
	ctx.DB.ServeHTTP(w, req)
	var resp struct {
		Data dto.BlogResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	blogID := resp.Data.ID

	// 2. Create 50 more versions (total 51)
	for i := 0; i < 50; i++ {
		reqBody := dto.CreateVersionRequest{
			ChangeSummary: fmt.Sprintf("Manual version %d", i+1),
		}
		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.UserAReader.Token)
		ctx.DB.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// 3. Verify count is 50 (not 51)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%s/versions", blogID), nil)
	req.Header.Set("Authorization", ctx.UserAReader.Token)
	ctx.DB.ServeHTTP(w, req)

	var vResp struct {
		Meta struct {
			Total int64 `json:"total"`
		} `json:"meta"`
	}
	json.Unmarshal(w.Body.Bytes(), &vResp)
	assert.Equal(t, int64(50), vResp.Meta.Total)
}
