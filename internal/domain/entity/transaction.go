package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TransactionType represents the type of transaction
type TransactionType string

// TransactionType enum values
const (
	TransactionTypeSubscription TransactionType = "SUBSCRIPTION"
	TransactionTypeSeries       TransactionType = "SERIES"
	TransactionTypeDonation     TransactionType = "DONATION"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

// TransactionStatus enum values
const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailed  TransactionStatus = "FAILED"
)

// TransactionProvider represents the payment provider
type TransactionProvider string

// TransactionProvider enum values
const (
	TransactionProviderSEPAY TransactionProvider = "SEPAY"
)

// TransactionGateway represents the payment gateway method
type TransactionGateway string

// TransactionGateway enum values
const (
	TransactionGatewayVietQR       TransactionGateway = "VIETQR"
	TransactionGatewayBankTransfer TransactionGateway = "BANK_TRANSFER"
)

// Transaction represents a payment transaction
type Transaction struct {
	ID             uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID              `gorm:"type:uuid;not null;index" json:"userId"`
	Amount         decimal.Decimal        `gorm:"type:decimal(19,4);not null" json:"amount"`
	Currency       string                 `gorm:"size:3;not null;default:'VND'" json:"currency"`
	Provider       TransactionProvider    `gorm:"size:20;not null" json:"provider"`
	Gateway        *TransactionGateway    `gorm:"size:20" json:"gateway,omitempty"`
	Type           TransactionType        `gorm:"size:20;not null" json:"type"`
	Status         TransactionStatus      `gorm:"size:20;not null;default:'PENDING'" json:"status"`
	TargetID       *uuid.UUID             `gorm:"type:uuid" json:"targetId,omitempty"`
	PlanID         *string                `gorm:"size:50" json:"planId,omitempty"`
	Content        string                 `gorm:"type:text" json:"content,omitempty"`
	SePayID        string                 `gorm:"size:255;not null;unique;index" json:"sepayId"`
	ReferenceCode  string                 `gorm:"size:255;not null;index" json:"referenceCode"`
	OrderID        string                 `gorm:"size:255;not null;index" json:"orderId"`
	WebhookPayload map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"webhookPayload"`
	CreatedAt      time.Time              `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName returns the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}
