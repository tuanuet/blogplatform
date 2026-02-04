package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/payment"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/adapter"
	"github.com/aiagent/internal/infrastructure/config"
	"github.com/aiagent/internal/infrastructure/persistence/postgres/repository"
	paymentHandler "github.com/aiagent/internal/interfaces/http/handler/payment"
	"github.com/aiagent/internal/interfaces/http/router"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	runMigrations(t, db)

	cleanup := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return db, cleanup
}

func runMigrations(t *testing.T, db *gorm.DB) {
	// Locate migrations directory relative to the test file
	// Assuming test is run from project root or tests/integration
	wd, err := os.Getwd()
	require.NoError(t, err)

	// Look for migrations folder walking up
	var migrationsDir string
	for {
		if _, err := os.Stat(filepath.Join(wd, "migrations")); err == nil {
			migrationsDir = filepath.Join(wd, "migrations")
			break
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("migrations directory not found")
		}
		wd = parent
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	require.NoError(t, err)
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		require.NoError(t, err)
		err = db.Exec(string(content)).Error
		require.NoError(t, err, "failed to execute migration %s", file)
	}
}

// Mock auth middleware for testing
func mockAuthMiddleware(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestPaymentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup Dependencies
	txRepo := repository.NewTransactionRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	planRepo := repository.NewSubscriptionPlanRepository(db)
	userRepo := repository.NewUserRepository(db)
	purchaseRepo := repository.NewUserSeriesPurchaseRepository(db)

	cfg := &config.SePayConfig{
		APIKey:       "test-api-key",
		BankName:     "Test Bank",
		BankAccount:  "123456789",
		BankOwner:    "TEST OWNER",
		BankBranch:   "Test Branch",
		WebhookToken: "test-webhook-token",
	}
	sepayAdapter := adapter.NewSePayAdapter(cfg)

	paymentSvc := service.NewPaymentService(db, txRepo, subRepo, purchaseRepo, planRepo, sepayAdapter)
	createPaymentUC := payment.NewCreatePaymentUseCase(paymentSvc)
	processWebhookUC := payment.NewProcessWebhookUseCase(paymentSvc)

	paymentH := paymentHandler.NewPaymentHandler(createPaymentUC)
	webhookH := paymentHandler.NewWebhookHandler(processWebhookUC, cfg.WebhookToken)

	// Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create test user and author
	userID := uuid.New()
	authorID := uuid.New()

	// Insert user and author directly into DB
	user := &entity.User{
		ID:       userID,
		Email:    "user@example.com",
		Name:     "Test User",
		IsActive: true,
	}
	err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	author := &entity.User{
		ID:       authorID,
		Email:    "author@example.com",
		Name:     "Test Author",
		IsActive: true,
	}
	err = userRepo.Create(context.Background(), author)
	require.NoError(t, err)

	// Create subscription
	sub := &entity.Subscription{
		ID:           uuid.New(),
		SubscriberID: userID,
		AuthorID:     authorID,
		Tier:         "FREE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = subRepo.Create(context.Background(), sub)
	require.NoError(t, err)

	// Create subscription plan for the author
	plan := &entity.SubscriptionPlan{
		ID:           uuid.New(),
		AuthorID:     authorID,
		Tier:         entity.TierSilver,
		DurationDays: 30,
		Price:        decimal.NewFromInt(50000),
		IsActive:     true,
	}
	err = planRepo.Create(context.Background(), plan)
	require.NoError(t, err)

	// Setup router with mock auth
	router.RegisterPaymentRoutes(r.Group("/api/v1"), paymentH, webhookH, mockAuthMiddleware(userID))

	t.Run("Scenario 1: Happy Path (Subscription)", func(t *testing.T) {
		// 1. Create payment request
		planID := plan.ID.String()
		targetID := authorID.String()
		reqBody := dto.CreatePaymentRequest{
			UserID:   userID.String(),
			Amount:   decimal.NewFromInt(50000),
			Type:     entity.TransactionTypeSubscription,
			Gateway:  entity.TransactionGatewayBankTransfer,
			TargetID: &targetID,
			PlanID:   &planID,
		}

		jsonBody, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var container struct {
			Success bool                      `json:"success"`
			Data    dto.CreatePaymentResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &container)
		require.NoError(t, err)
		require.True(t, container.Success)
		resp := container.Data

		assert.NotEmpty(t, resp.ReferenceCode)
		assert.Equal(t, reqBody.Amount.String(), resp.Amount.String())

		// 2. Verify PENDING transaction in DB
		tx, err := txRepo.FindByRefID(context.Background(), resp.ReferenceCode)
		require.NoError(t, err)
		assert.Equal(t, entity.TransactionStatusPending, tx.Status)
		assert.Equal(t, userID, tx.UserID)

		// 3. Simulate SePay Webhook
		sepayID := int64(123456)
		webhookPayload := dto.ProcessWebhookRequest{
			ID:              sepayID,
			Gateway:         "MBBank",
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			AccountNumber:   "123456789",
			Code:            "CODE123",
			Content:         resp.ReferenceCode, // Important: match the reference code
			TransferType:    "in",
			TransferAmount:  reqBody.Amount,
			Accumulated:     reqBody.Amount,
			SubAccount:      "",
			ReferenceCode:   "REF123",
			Description:     "Payment for sub",
		}

		webhookJson, _ := json.Marshal(webhookPayload)
		wWebhook := httptest.NewRecorder()
		reqWebhook, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(webhookJson))
		reqWebhook.Header.Set("Content-Type", "application/json")
		// Set authorization header with token matching config
		reqWebhook.Header.Set("Authorization", "Bearer test-webhook-token")

		r.ServeHTTP(wWebhook, reqWebhook)

		require.Equal(t, http.StatusOK, wWebhook.Code)

		// 4. Verify transaction SUCCESS
		txUpdated, err := txRepo.FindByRefID(context.Background(), resp.ReferenceCode)
		require.NoError(t, err)
		assert.Equal(t, entity.TransactionStatusSuccess, txUpdated.Status)
		assert.Equal(t, fmt.Sprintf("%d", sepayID), txUpdated.SePayID)

		// 5. Verify Subscription expiry updated
		subUpdated, err := subRepo.FindActiveSubscription(context.Background(), userID, authorID)
		require.NoError(t, err)
		require.NotNil(t, subUpdated)
		assert.True(t, subUpdated.ExpiresAt.After(time.Now()))
		// Should be roughly 30 days from now
		assert.True(t, subUpdated.ExpiresAt.After(time.Now().AddDate(0, 0, 29)))
	})

	t.Run("Scenario 2: Idempotency", func(t *testing.T) {
		// Use the same transaction as Scenario 1 (which is already SUCCESS)
		// We need the reference code from scenario 1, but variables are scoped.
		// I'll re-query the transaction or create a new one.
		// For simplicity, let's just create a NEW transaction and succeed it, then try again.

		// Create a manual pending transaction
		orderID := "ORDER-IDEMP-TEST"
		tx := &entity.Transaction{
			ID:            uuid.New(),
			UserID:        userID,
			Amount:        decimal.NewFromInt(100000),
			Currency:      "VND",
			Provider:      entity.TransactionProviderSEPAY,
			Type:          entity.TransactionTypeDonation,
			Status:        entity.TransactionStatusPending,
			OrderID:       orderID,
			ReferenceCode: orderID,
		}
		err := txRepo.Create(context.Background(), tx)
		require.NoError(t, err)

		// Send Webhook 1st time
		sepayID := int64(999999)
		webhookPayload := dto.ProcessWebhookRequest{
			ID:              sepayID,
			Gateway:         "MBBank",
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			AccountNumber:   "123456789",
			Code:            "CODE999",
			Content:         orderID,
			TransferType:    "in",
			TransferAmount:  tx.Amount,
			Accumulated:     tx.Amount,
			SubAccount:      "",
			ReferenceCode:   "REF999",
			Description:     "Donation",
		}

		webhookJson, _ := json.Marshal(webhookPayload)

		// Helper to send webhook
		sendWebhook := func() *httptest.ResponseRecorder {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(webhookJson))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-webhook-token")
			r.ServeHTTP(w, req)
			return w
		}

		// 1st call
		w1 := sendWebhook()
		require.Equal(t, http.StatusOK, w1.Code)

		// Verify SUCCESS
		tx1, _ := txRepo.FindByRefID(context.Background(), orderID)
		assert.Equal(t, entity.TransactionStatusSuccess, tx1.Status)

		// 2. Send same webhook twice
		w2 := sendWebhook()
		require.Equal(t, http.StatusOK, w2.Code)

		// Verify status doesn't change and logic runs once
		// (We can't easily verify logic ran once without logs mock, but we check response is OK and state is stable)
		tx2, _ := txRepo.FindByRefID(context.Background(), orderID)
		assert.Equal(t, entity.TransactionStatusSuccess, tx2.Status)
		assert.Equal(t, tx1.UpdatedAt.Unix(), tx2.UpdatedAt.Unix()) // Should not have updated timestamp if skipped?
		// Actually, if it returns early, UpdatedAt might be same.
	})

	t.Run("Scenario 3: Error Cases", func(t *testing.T) {
		// 1. Invalid API Key on webhook
		webhookPayload := dto.ProcessWebhookRequest{
			ID:      int64(88888),
			Content: "SOME-CONTENT",
		}
		webhookJson, _ := json.Marshal(webhookPayload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(webhookJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer WRONG-TOKEN")

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusUnauthorized, w.Code)

		// 2. Invalid JSON
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBufferString("{invalid-json"))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer test-webhook-token")

		r.ServeHTTP(w2, req2)
		require.Equal(t, http.StatusBadRequest, w2.Code)
	})
}
