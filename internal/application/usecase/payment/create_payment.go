package payment

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/service"
)

type CreatePaymentUseCase interface {
	Execute(ctx context.Context, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error)
}

type createPaymentUseCase struct {
	paymentService service.PaymentService
}

func NewCreatePaymentUseCase(paymentService service.PaymentService) CreatePaymentUseCase {
	return &createPaymentUseCase{
		paymentService: paymentService,
	}
}

func (u *createPaymentUseCase) Execute(ctx context.Context, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error) {
	serviceReq := service.CreatePaymentRequest{
		UserID:   req.UserID,
		Amount:   req.Amount,
		Type:     req.Type,
		Gateway:  req.Gateway,
		TargetID: req.TargetID,
		PlanID:   req.PlanID,
	}

	resp, err := u.paymentService.InitPayment(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	return &dto.CreatePaymentResponse{
		OrderID:       resp.OrderID,
		Amount:        resp.Amount,
		Gateway:       resp.Gateway,
		QRDataURL:     resp.QRDataURL,
		QRData:        resp.QRData,
		BankName:      resp.BankName,
		AccountNo:     resp.AccountNo,
		AccountName:   resp.AccountName,
		ReferenceCode: resp.ReferenceCode,
	}, nil
}
