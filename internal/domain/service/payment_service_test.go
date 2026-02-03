package service_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/adapter"
	adapter_mocks "github.com/aiagent/internal/infrastructure/adapter/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPaymentService_InitPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockPurchaseRepo := mocks.NewMockUserSeriesPurchaseRepository(ctrl)
	mockSePayAdapter := adapter_mocks.NewMockSePayAdapter(ctrl)

	db, _, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	svc := service.NewPaymentService(
		gormDB,
		mockTxRepo,
		mockSubRepo,
		mockPurchaseRepo,
		mockSePayAdapter,
	)

	ctx := context.Background()
	userID := uuid.New()
	amount := decimal.NewFromInt(100000)

	t.Run("success_vietqr", func(t *testing.T) {
		req := service.CreatePaymentRequest{
			UserID:  userID.String(),
			Amount:  amount,
			Type:    entity.TransactionTypeSubscription,
			Gateway: entity.TransactionGatewayVietQR,
		}

		mockTxRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockSePayAdapter.EXPECT().CreateVietQR(gomock.Any(), gomock.Any()).Return(&adapter.VietQRResponse{
			Status: 200,
			Data: struct {
				QrCode    string `json:"qrCode"`
				QrDataURL string `json:"qrDataURL"`
			}{
				QrCode:    "qr-code-data",
				QrDataURL: "qr-data-url",
			},
		}, nil)

		resp, err := svc.InitPayment(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, amount, resp.Amount)
		assert.Equal(t, entity.TransactionGatewayVietQR, resp.Gateway)
		assert.Equal(t, "qr-data-url", resp.QRDataURL)
	})

	t.Run("success_bank_transfer", func(t *testing.T) {
		req := service.CreatePaymentRequest{
			UserID:  userID.String(),
			Amount:  amount,
			Type:    entity.TransactionTypeSubscription,
			Gateway: entity.TransactionGatewayBankTransfer,
		}

		mockTxRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockSePayAdapter.EXPECT().GetBankTransferInfo().Return(adapter.BankTransferInfo{
			BankName:    "MB Bank",
			AccountNo:   "123456789",
			AccountName: "AGENT",
		})

		resp, err := svc.InitPayment(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "MB Bank", resp.BankName)
		assert.Equal(t, "123456789", resp.AccountNo)
	})
}

func TestPaymentService_HandleSePayWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockPurchaseRepo := mocks.NewMockUserSeriesPurchaseRepository(ctrl)
	mockSePayAdapter := adapter_mocks.NewMockSePayAdapter(ctrl)

	db, sqlMock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	svc := service.NewPaymentService(
		gormDB,
		mockTxRepo,
		mockSubRepo,
		mockPurchaseRepo,
		mockSePayAdapter,
	)

	ctx := context.Background()
	orderID := "ORDER-SEPAY-123456"
	sePayID := "1001"
	userID := uuid.New()
	amount := decimal.NewFromInt(100000)

	payload := service.SePayWebhookPayload{
		ID:             1001,
		Content:        "ORDER-SEPAY-123456",
		TransferAmount: amount,
	}

	t.Run("transaction_not_found", func(t *testing.T) {
		mockTxRepo.EXPECT().FindBySePayID(ctx, sePayID).Return(nil, nil)
		mockTxRepo.EXPECT().FindByRefID(ctx, orderID).Return(nil, nil)

		tx, err := svc.HandleSePayWebhook(ctx, payload)

		assert.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "transaction not found")
	})

	t.Run("idempotency_check", func(t *testing.T) {
		existingTx := &entity.Transaction{
			ID:      uuid.New(),
			SePayID: sePayID,
			Status:  entity.TransactionStatusSuccess,
		}
		mockTxRepo.EXPECT().FindBySePayID(ctx, sePayID).Return(existingTx, nil)

		tx, err := svc.HandleSePayWebhook(ctx, payload)

		assert.NoError(t, err)
		assert.Equal(t, existingTx, tx)
	})

	t.Run("success_subscription", func(t *testing.T) {
		tx := &entity.Transaction{
			ID:      uuid.New(),
			UserID:  userID,
			OrderID: orderID,
			Amount:  amount,
			Type:    entity.TransactionTypeSubscription,
			Status:  entity.TransactionStatusPending,
		}
		targetID := uuid.New()
		planID := "premium"
		tx.TargetID = &targetID
		tx.PlanID = &planID

		mockTxRepo.EXPECT().FindBySePayID(ctx, sePayID).Return(nil, nil)
		mockTxRepo.EXPECT().FindByRefID(ctx, orderID).Return(tx, nil)

		// Expectations for transaction and WithTx
		sqlMock.ExpectBegin()
		mockTxRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxRepo)
		mockSubRepo.EXPECT().WithTx(gomock.Any()).Return(mockSubRepo)
		mockPurchaseRepo.EXPECT().WithTx(gomock.Any()).Return(mockPurchaseRepo)

		mockTxRepo.EXPECT().Update(ctx, tx).Return(nil)
		mockSubRepo.EXPECT().UpdateExpiry(ctx, userID, targetID, gomock.Any(), planID).Return(nil)
		sqlMock.ExpectCommit()

		result, err := svc.HandleSePayWebhook(ctx, payload)

		assert.NoError(t, err)
		assert.Equal(t, entity.TransactionStatusSuccess, result.Status)
		assert.Equal(t, sePayID, result.SePayID)
	})

	t.Run("success_series_purchase", func(t *testing.T) {
		seriesID := uuid.New()
		tx := &entity.Transaction{
			ID:       uuid.New(),
			UserID:   userID,
			OrderID:  orderID,
			Amount:   amount,
			Type:     entity.TransactionTypeSeries,
			Status:   entity.TransactionStatusPending,
			TargetID: &seriesID,
		}

		mockTxRepo.EXPECT().FindBySePayID(ctx, sePayID).Return(nil, nil)
		mockTxRepo.EXPECT().FindByRefID(ctx, orderID).Return(tx, nil)

		// Expectations for transaction and WithTx
		sqlMock.ExpectBegin()
		mockTxRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxRepo)
		mockSubRepo.EXPECT().WithTx(gomock.Any()).Return(mockSubRepo)
		mockPurchaseRepo.EXPECT().WithTx(gomock.Any()).Return(mockPurchaseRepo)

		mockTxRepo.EXPECT().Update(ctx, tx).Return(nil)
		mockPurchaseRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		sqlMock.ExpectCommit()

		result, err := svc.HandleSePayWebhook(ctx, payload)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestPaymentService_GetTransactionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)
	db, _, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	svc := service.NewPaymentService(gormDB, mockTxRepo, nil, nil, nil)

	ctx := context.Background()
	orderID := "ORDER-123"
	expectedTx := &entity.Transaction{OrderID: orderID, Status: entity.TransactionStatusSuccess}

	mockTxRepo.EXPECT().FindByRefID(ctx, orderID).Return(expectedTx, nil)

	tx, err := svc.GetTransactionStatus(ctx, orderID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
}

func TestPaymentService_VerifySePayWebhookSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := adapter_mocks.NewMockSePayAdapter(ctrl)
	db, _, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	svc := service.NewPaymentService(gormDB, nil, nil, nil, mockAdapter)

	payload := map[string]interface{}{"foo": "bar"}
	signature := "valid-sig"

	mockAdapter.EXPECT().VerifyWebhookSignature(payload, signature).Return(true)

	result := svc.VerifySePayWebhookSignature(payload, signature)

	assert.True(t, result)
}
