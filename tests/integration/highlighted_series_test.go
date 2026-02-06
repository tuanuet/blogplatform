package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/series"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	seriesHandler "github.com/aiagent/internal/interfaces/http/handler/series"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type SeriesTestContext struct {
	DB      *gorm.DB
	Server  *httptest.Server
	Cleanup func()
}

func setupSeriesTestServer(t *testing.T) *SeriesTestContext {
	db, cleanup := setupTestDB(t)

	// Setup repositories
	seriesRepo := repository.NewSeriesRepository(db)

	// Setup usecases
	seriesUC := series.NewSeriesUseCase(seriesRepo)

	// Setup handler
	handler := seriesHandler.NewSeriesHandler(seriesUC)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	// Register only the highlighted route for this specific test
	// We manually register to avoid setting up complex auth middleware dependencies
	// required by the full RegisterSeriesRoutes function.
	v1.GET("/series/highlighted", handler.GetHighlightedSeries)

	server := httptest.NewServer(r)

	return &SeriesTestContext{
		DB:     db,
		Server: server,
		Cleanup: func() {
			server.Close()
			cleanup()
		},
	}
}

func TestHighlightedSeriesAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := setupSeriesTestServer(t)
	defer ctx.Cleanup()

	// Repositories for seeding
	userRepo := repository.NewUserRepository(ctx.DB)
	seriesRepo := repository.NewSeriesRepository(ctx.DB)
	purchaseRepo := repository.NewUserSeriesPurchaseRepository(ctx.DB)

	// 1. Create Author
	authorID := uuid.New()
	author := &entity.User{
		ID:       authorID,
		Email:    "author@example.com",
		Name:     "Test Author",
		IsActive: true,
	}
	err := userRepo.Create(context.Background(), author)
	require.NoError(t, err)

	// 2. Create Series A
	seriesA := &entity.Series{
		ID:          uuid.New(),
		AuthorID:    authorID,
		Title:       "Series A - Popular",
		Slug:        "series-a-popular",
		Description: "This series has many subscribers",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = seriesRepo.Create(context.Background(), seriesA)
	require.NoError(t, err)

	// 3. Create Series B
	seriesB := &entity.Series{
		ID:          uuid.New(),
		AuthorID:    authorID,
		Title:       "Series B - New",
		Slug:        "series-b-new",
		Description: "This series has no subscribers",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = seriesRepo.Create(context.Background(), seriesB)
	require.NoError(t, err)

	// 4. Create 5 purchases for Series A
	for i := 0; i < 5; i++ {
		userID := uuid.New()
		// We might need to create the user first due to FK constraints
		user := &entity.User{
			ID:       userID,
			Email:    "subscriber" + uuid.New().String() + "@example.com",
			Name:     fmt.Sprintf("Subscriber %d", i),
			IsActive: true,
		}
		err = userRepo.Create(context.Background(), user)
		require.NoError(t, err)

		purchase := &entity.UserSeriesPurchase{
			UserID:    userID,
			SeriesID:  seriesA.ID,
			Amount:    decimal.NewFromInt(100),
			CreatedAt: time.Now(),
		}
		err = purchaseRepo.Create(context.Background(), purchase)
		require.NoError(t, err)
	}

	t.Run("Success - Get highlighted series sorted by subscribers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.Server.URL+"/api/v1/series/highlighted", nil)

		w := httptest.NewRecorder()
		ctx.Server.Config.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool                             `json:"success"`
			Data    []*dto.HighlightedSeriesResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)

		// Should return both series (limit is 10)
		assert.Len(t, resp.Data, 2)

		// First one should be Series A (5 subscribers)
		assert.Equal(t, seriesA.ID, resp.Data[0].ID)
		assert.Equal(t, 5, resp.Data[0].SubscriberCount)

		// Second one should be Series B (0 subscribers)
		assert.Equal(t, seriesB.ID, resp.Data[1].ID)
		assert.Equal(t, 0, resp.Data[1].SubscriberCount)
	})
}
