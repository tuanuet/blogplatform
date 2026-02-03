package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/adapter"
	"github.com/aiagent/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// CreatePaymentRequest represents the request to initiate a payment
type CreatePaymentRequest struct {
	UserID   string                    `json:"userId" validate:"required"`
	Amount   decimal.Decimal           `json:"amount" validate:"required,gt=0"`
	Type     entity.TransactionType    `json:"type" validate:"required"`
	Gateway  entity.TransactionGateway `json:"gateway" validate:"required"`
	TargetID *string                   `json:"targetId,omitempty"` // ID of Subscription Author or Series
	PlanID   *string                   `json:"planId,omitempty"`   // Subscription plan ID
}

// PaymentResponse represents the response after initiating a payment
type PaymentResponse struct {
	OrderID       string                    `json:"orderId"`
	Amount        decimal.Decimal           `json:"amount"`
	Gateway       entity.TransactionGateway `json:"gateway"`
	QRDataURL     string                    `json:"qrDataUrl,omitempty"`
	QRData        string                    `json:"qrData,omitempty"`
	BankName      string                    `json:"bankName,omitempty"`
	AccountNo     string                    `json:"accountNo,omitempty"`
	AccountName   string                    `json:"accountName,omitempty"`
	ReferenceCode string                    `json:"referenceCode"`
}

// SePayWebhookPayload represents the payload from SePay webhook
type SePayWebhookPayload struct {
	ID              int64           `json:"id"`
	Gateway         string          `json:"gateway"`
	TransactionDate string          `json:"transactionDate"`
	AccountNumber   string          `json:"accountNumber"`
	Code            string          `json:"code"`
	Content         string          `json:"content"`
	TransferType    string          `json:"transferType"`
	TransferAmount  decimal.Decimal `json:"transferAmount"`
	Accumulated     decimal.Decimal `json:"accumulated"`
	SubAccount      string          `json:"subAccount"`
	ReferenceCode   string          `json:"referenceCode"`
	Description     string          `json:"description"`
}

// PaymentService defines the interface for payment operations
type PaymentService interface {
	InitPayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResponse, error)
	HandleSePayWebhook(ctx context.Context, payload SePayWebhookPayload) (*entity.Transaction, error)
	GetTransactionStatus(ctx context.Context, orderID string) (*entity.Transaction, error)
	VerifySePayWebhookSignature(payload map[string]interface{}, signature string) bool
}

type paymentService struct {
	db           *gorm.DB
	txRepo       repository.TransactionRepository
	subRepo      repository.SubscriptionRepository
	purchaseRepo repository.UserSeriesPurchaseRepository
	sepayAdapter adapter.SePayAdapter
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(
	db *gorm.DB,
	txRepo repository.TransactionRepository,
	subRepo repository.SubscriptionRepository,
	purchaseRepo repository.UserSeriesPurchaseRepository,
	sepayAdapter adapter.SePayAdapter,
) PaymentService {
	return &paymentService{
		db:           db,
		txRepo:       txRepo,
		subRepo:      subRepo,
		purchaseRepo: purchaseRepo,
		sepayAdapter: sepayAdapter,
	}
}

// InitPayment initiates a payment process
func (s *paymentService) InitPayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResponse, error) {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	orderID := fmt.Sprintf("ORDER-SEPAY-%s", uuid.New().String())

	tx := &entity.Transaction{
		ID:            uuid.New(),
		UserID:        userUUID,
		Amount:        req.Amount,
		Currency:      "VND",
		Provider:      entity.TransactionProviderSEPAY,
		Gateway:       &req.Gateway,
		Type:          req.Type,
		Status:        entity.TransactionStatusPending,
		OrderID:       orderID,
		ReferenceCode: orderID,
	}

	if req.TargetID != nil {
		targetUUID, err := uuid.Parse(*req.TargetID)
		if err == nil {
			tx.TargetID = &targetUUID
		}
	}
	tx.PlanID = req.PlanID

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	resp := &PaymentResponse{
		OrderID:       orderID,
		Amount:        req.Amount,
		Gateway:       req.Gateway,
		ReferenceCode: orderID,
	}

	if req.Gateway == entity.TransactionGatewayVietQR {
		qrReq := adapter.CreateVietQRRequest{
			Amount:  int(req.Amount.IntPart()),
			AddInfo: orderID,
		}
		qrResp, err := s.sepayAdapter.CreateVietQR(ctx, qrReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create VietQR: %w", err)
		}
		resp.QRData = qrResp.Data.QrCode
		resp.QRDataURL = qrResp.Data.QrDataURL
	} else if req.Gateway == entity.TransactionGatewayBankTransfer {
		bankInfo := s.sepayAdapter.GetBankTransferInfo()
		resp.BankName = bankInfo.BankName
		resp.AccountNo = bankInfo.AccountNo
		resp.AccountName = bankInfo.AccountName
	}

	return resp, nil
}

// HandleSePayWebhook processes the incoming webhook from SePay
func (s *paymentService) HandleSePayWebhook(ctx context.Context, payload SePayWebhookPayload) (*entity.Transaction, error) {
	sePayID := strconv.FormatInt(payload.ID, 10)

	// Idempotency check
	existingTx, err := s.txRepo.FindBySePayID(ctx, sePayID)
	if err == nil && existingTx != nil {
		return existingTx, nil
	}

	// Find original transaction by order_id (from content)
	tx, err := s.txRepo.FindByRefID(ctx, payload.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	if tx.Status == entity.TransactionStatusSuccess {
		return tx, nil
	}

	// Use transaction to ensure atomicity
	var resultTx *entity.Transaction
	err = s.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		txRepo := s.txRepo.WithTx(dbTx)
		subRepo := s.subRepo.WithTx(dbTx)
		purchaseRepo := s.purchaseRepo.WithTx(dbTx)

		// Update transaction with SePayID for idempotency
		tx.SePayID = sePayID
		tx.Status = entity.TransactionStatusSuccess
		if err := txRepo.Update(ctx, tx); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		// Grant benefits
		switch tx.Type {
		case entity.TransactionTypeSubscription:
			if tx.TargetID != nil && tx.PlanID != nil {
				// Determine expiry based on PlanID
				days := 30 // Default to 30 days
				if *tx.PlanID == "1_YEAR" {
					days = 365
				} else if *tx.PlanID == "1_MONTH" {
					days = 30
				}
				expiry := time.Now().AddDate(0, 0, days)
				if err := subRepo.UpdateExpiry(ctx, tx.UserID, *tx.TargetID, expiry, *tx.PlanID); err != nil {
					return fmt.Errorf("failed to update subscription: %w", err)
				}
			}
		case entity.TransactionTypeSeries:
			if tx.TargetID != nil {
				purchase := &entity.UserSeriesPurchase{
					UserID:   tx.UserID,
					SeriesID: *tx.TargetID,
					Amount:   tx.Amount,
				}
				if err := purchaseRepo.Create(ctx, purchase); err != nil {
					return fmt.Errorf("failed to record series purchase: %w", err)
				}
			}
		case entity.TransactionTypeDonation:
			// Just log
			logger.Info("Donation received", map[string]interface{}{
				"amount": tx.Amount,
				"userId": tx.UserID,
			})
		}
		resultTx = tx
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultTx, nil
}

// GetTransactionStatus retrieves the status of a transaction
func (s *paymentService) GetTransactionStatus(ctx context.Context, orderID string) (*entity.Transaction, error) {
	return s.txRepo.FindByRefID(ctx, orderID)
}

// VerifySePayWebhookSignature verifies the webhook signature
func (s *paymentService) VerifySePayWebhookSignature(payload map[string]interface{}, signature string) bool {
	return s.sepayAdapter.VerifyWebhookSignature(payload, signature)
}
