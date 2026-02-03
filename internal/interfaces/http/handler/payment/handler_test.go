package payment_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockCreatePaymentUseCase, *gin.Engine, payment.PaymentHandler) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockCreatePaymentUseCase(ctrl)
	h := payment.NewPaymentHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Mock middleware to set user ID in context
	userID := uuid.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	return ctrl, mockUC, r, h
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.POST("/api/v1/payments", h.CreatePayment)

	t.Run("success", func(t *testing.T) {
		planID := "plan-123"
		req := dto.CreatePaymentRequest{
			Amount:  decimal.NewFromInt(100000),
			Type:    entity.TransactionTypeSubscription,
			Gateway: entity.TransactionGatewayVietQR,
			PlanID:  &planID,
		}
		body, _ := json.Marshal(req)

		expectedResp := &dto.CreatePaymentResponse{
			OrderID:       uuid.New().String(),
			Amount:        decimal.NewFromInt(100000),
			Gateway:       entity.TransactionGatewayVietQR,
			QRDataURL:     "https://example.com/qr.png",
			QRData:        "QR_DATA",
			BankName:      "VietQR Bank",
			AccountNo:     "123456789",
			AccountName:   "Test Account",
			ReferenceCode: "REF123",
		}

		mockUC.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error) {
				// Verify that user ID from context is used instead of request body
				assert.NotEmpty(t, r.UserID, "UserID should be extracted from context")
				return expectedResp, nil
			})

		httpReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp["success"])
	})

	t.Run("missing user in context", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rNoUser := gin.New()
		hNoUser := payment.NewPaymentHandler(mockUC)
		rNoUser.POST("/api/v1/payments", hNoUser.CreatePayment)

		req := dto.CreatePaymentRequest{
			Amount:  decimal.NewFromInt(100000),
			Type:    entity.TransactionTypeSubscription,
			Gateway: entity.TransactionGatewayVietQR,
		}
		body, _ := json.Marshal(req)

		httpReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rNoUser.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		rInvalid := gin.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		h2 := payment.NewPaymentHandler(mocks.NewMockCreatePaymentUseCase(ctrl))
		// Add middleware to set user ID
		rInvalid.Use(func(c *gin.Context) {
			c.Set("userID", uuid.New())
			c.Next()
		})
		rInvalid.POST("/api/v1/payments", h2.CreatePayment)

		httpReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rInvalid.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		req := dto.CreatePaymentRequest{
			Amount:  decimal.NewFromInt(100000),
			Type:    entity.TransactionTypeSubscription,
			Gateway: entity.TransactionGatewayVietQR,
		}
		body, _ := json.Marshal(req)

		mockUC.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		httpReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPaymentHandler_CreatePayment_Donation(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.POST("/api/v1/payments", h.CreatePayment)

	targetID := "author-123"
	req := dto.CreatePaymentRequest{
		Amount:   decimal.NewFromInt(50000),
		Type:     entity.TransactionTypeDonation,
		Gateway:  entity.TransactionGatewayBankTransfer,
		TargetID: &targetID,
	}
	body, _ := json.Marshal(req)

	expectedResp := &dto.CreatePaymentResponse{
		OrderID:       uuid.New().String(),
		Amount:        decimal.NewFromInt(50000),
		Gateway:       entity.TransactionGatewayBankTransfer,
		BankName:      "Test Bank",
		AccountNo:     "987654321",
		AccountName:   "Author Account",
		ReferenceCode: "DONATION123",
	}

	mockUC.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(expectedResp, nil)

	httpReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
}
