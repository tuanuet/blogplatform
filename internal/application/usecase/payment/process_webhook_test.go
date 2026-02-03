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

func TestProcessWebhookUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentService := mocks.NewMockPaymentService(ctrl)
	useCase := NewProcessWebhookUseCase(mockPaymentService)

	ctx := context.Background()
	req := dto.ProcessWebhookRequest{
		ID:             123,
		Gateway:        "VIETQR",
		Content:        "ORDER-123",
		TransferAmount: decimal.NewFromInt(100000),
	}

	expectedTx := &entity.Transaction{
		ID:      [16]byte{}, // UUID
		OrderID: "ORDER-123",
		Status:  entity.TransactionStatusSuccess,
	}

	t.Run("success", func(t *testing.T) {
		mockPaymentService.EXPECT().
			HandleSePayWebhook(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, payload service.SePayWebhookPayload) (*entity.Transaction, error) {
				assert.Equal(t, req.ID, payload.ID)
				assert.Equal(t, req.Gateway, payload.Gateway)
				assert.Equal(t, req.Content, payload.Content)
				assert.Equal(t, req.TransferAmount, payload.TransferAmount)
				return expectedTx, nil
			})

		resp, err := useCase.Execute(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expectedTx, resp)
	})

	t.Run("service error", func(t *testing.T) {
		mockPaymentService.EXPECT().
			HandleSePayWebhook(ctx, gomock.Any()).
			Return(nil, errors.New("service error"))

		resp, err := useCase.Execute(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
