package payment

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
)

type ProcessWebhookUseCase interface {
	Execute(ctx context.Context, req dto.ProcessWebhookRequest) (*entity.Transaction, error)
}

type processWebhookUseCase struct {
	paymentService service.PaymentService
}

func NewProcessWebhookUseCase(paymentService service.PaymentService) ProcessWebhookUseCase {
	return &processWebhookUseCase{
		paymentService: paymentService,
	}
}

func (u *processWebhookUseCase) Execute(ctx context.Context, req dto.ProcessWebhookRequest) (*entity.Transaction, error) {
	payload := service.SePayWebhookPayload{
		ID:              req.ID,
		Gateway:         req.Gateway,
		TransactionDate: req.TransactionDate,
		AccountNumber:   req.AccountNumber,
		Code:            req.Code,
		Content:         req.Content,
		TransferType:    req.TransferType,
		TransferAmount:  req.TransferAmount,
		Accumulated:     req.Accumulated,
		SubAccount:      req.SubAccount,
		ReferenceCode:   req.ReferenceCode,
		Description:     req.Description,
	}

	return u.paymentService.HandleSePayWebhook(ctx, payload)
}
