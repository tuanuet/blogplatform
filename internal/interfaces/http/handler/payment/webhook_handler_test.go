package payment_test

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
	"github.com/aiagent/internal/interfaces/http/handler/payment"
	"github.com/aiagent/internal/interfaces/http/handler/payment/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const validAPIKey = "valid-api-key"

func setupWebhookTest(t *testing.T) (*gomock.Controller, *mocks.MockProcessWebhookUseCase, *gin.Engine, payment.WebhookHandler) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockProcessWebhookUseCase(ctrl)
	h := payment.NewWebhookHandler(mockUC, validAPIKey)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	return ctrl, mockUC, r, h
}

func TestWebhookHandler_HandleSePayWebhook(t *testing.T) {
	ctrl, mockUC, r, h := setupWebhookTest(t)
	defer ctrl.Finish()

	r.POST("/api/v1/webhooks/sepay", h.HandleSePayWebhook)

	t.Run("success", func(t *testing.T) {
		req := dto.ProcessWebhookRequest{
			ID:              12345,
			Gateway:         "bank",
			TransactionDate: "2024-01-15 10:30:00",
			AccountNumber:   "123456789",
			Code:            "TXN123",
			Content:         "Payment for subscription",
			TransferType:    "in",
			TransferAmount:  decimal.NewFromInt(100000),
			Accumulated:     decimal.NewFromInt(100000),
			SubAccount:      "SUB001",
			ReferenceCode:   "REF123",
			Description:     "Monthly subscription",
		}
		body, _ := json.Marshal(req)

		expectedTx := &entity.Transaction{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			Amount:        decimal.NewFromInt(100000),
			Currency:      "VND",
			Provider:      entity.TransactionProviderSEPAY,
			Gateway:       func() *entity.TransactionGateway { g := entity.TransactionGatewayBankTransfer; return &g }(),
			Type:          entity.TransactionTypeSubscription,
			Status:        entity.TransactionStatusSuccess,
			SePayID:       "12345",
			ReferenceCode: "REF123",
			OrderID:       "ORD123",
		}

		mockUC.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r dto.ProcessWebhookRequest) (*entity.Transaction, error) {
				assert.Equal(t, int64(12345), r.ID)
				assert.Equal(t, "bank", r.Gateway)
				assert.Equal(t, decimal.NewFromInt(100000), r.TransferAmount)
				return expectedTx, nil
			})

		httpReq, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+validAPIKey)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp["success"])
	})

	t.Run("missing API key", func(t *testing.T) {
		req := dto.ProcessWebhookRequest{
			ID:             12345,
			TransferAmount: decimal.NewFromInt(100000),
		}
		body, _ := json.Marshal(req)

		httpReq, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		// No Authorization header
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
	})

	t.Run("invalid API key", func(t *testing.T) {
		req := dto.ProcessWebhookRequest{
			ID:             12345,
			TransferAmount: decimal.NewFromInt(100000),
		}
		body, _ := json.Marshal(req)

		httpReq, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer invalid-api-key")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		httpReq, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+validAPIKey)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
	})

	t.Run("use case returns error", func(t *testing.T) {
		req := dto.ProcessWebhookRequest{
			ID:             12345,
			Gateway:        "bank",
			TransferAmount: decimal.NewFromInt(100000),
			ReferenceCode:  "REF123",
		}
		body, _ := json.Marshal(req)

		mockUC.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("transaction processing failed"))

		httpReq, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+validAPIKey)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, false, resp["success"])
	})
}
