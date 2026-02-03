package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/service/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePaymentUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := mocks.NewMockPaymentService(ctrl)
	useCase := NewCreatePaymentUseCase(mockPaymentService)

	ctx := context.Background()
	amount := decimal.NewFromInt(100000)
	userID := "user-uuid"
	targetID := "target-uuid"

	req := dto.CreatePaymentRequest{
		UserID:   userID,
		Amount:   amount,
		Type:     entity.TransactionTypeSubscription,
		Gateway:  entity.TransactionGatewayVietQR,
		TargetID: &targetID,
	}

	expectedServiceResp := &service.PaymentResponse{
		OrderID:   "order-123",
		Amount:    amount,
		Gateway:   entity.TransactionGatewayVietQR,
		QRDataURL: "https://qr.url",
		QRData:    "qr-code-data",
	}

	t.Run("success", func(t *testing.T) {
		// We expect InitPayment to be called with a request that has the same values
		mockPaymentService.EXPECT().
			InitPayment(ctx, gomock.Any()). // We verify details inside DoAndReturn or assume mapping is correct if test passes. Better to use a matcher.
			DoAndReturn(func(ctx context.Context, sReq service.CreatePaymentRequest) (*service.PaymentResponse, error) {
				assert.Equal(t, req.UserID, sReq.UserID)
				assert.Equal(t, req.Amount, sReq.Amount)
				assert.Equal(t, req.Type, sReq.Type)
				assert.Equal(t, req.Gateway, sReq.Gateway)
				assert.Equal(t, req.TargetID, sReq.TargetID)
				return expectedServiceResp, nil
			})

		resp, err := useCase.Execute(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedServiceResp.OrderID, resp.OrderID)
		assert.Equal(t, expectedServiceResp.QRDataURL, resp.QRDataURL)
	})

	t.Run("service error", func(t *testing.T) {
		mockPaymentService.EXPECT().
			InitPayment(ctx, gomock.Any()).
			Return(nil, errors.New("service error"))

		resp, err := useCase.Execute(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
