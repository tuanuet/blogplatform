package dto

import (
	"github.com/aiagent/internal/domain/entity"
	"github.com/shopspring/decimal"
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

// CreatePaymentResponse represents the response after initiating a payment
type CreatePaymentResponse struct {
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

// ProcessWebhookRequest represents the payload from SePay webhook
type ProcessWebhookRequest struct {
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
