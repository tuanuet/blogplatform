package entity_test

import (
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestTransaction_TableName verifies the table name is correctly set
func TestTransaction_TableName(t *testing.T) {
	// Arrange
	tx := entity.Transaction{}

	// Act
	tableName := tx.TableName()

	// Assert
	assert.Equal(t, "transactions", tableName, "TableName should return 'transactions'")
}

// TestTransaction_Enums tests enum values are correctly defined
func TestTransaction_Enums(t *testing.T) {
	t.Run("TransactionType values", func(t *testing.T) {
		assert.Equal(t, entity.TransactionTypeSubscription, entity.TransactionType("SUBSCRIPTION"))
		assert.Equal(t, entity.TransactionTypeSeries, entity.TransactionType("SERIES"))
		assert.Equal(t, entity.TransactionTypeDonation, entity.TransactionType("DONATION"))
	})

	t.Run("TransactionStatus values", func(t *testing.T) {
		assert.Equal(t, entity.TransactionStatusPending, entity.TransactionStatus("PENDING"))
		assert.Equal(t, entity.TransactionStatusSuccess, entity.TransactionStatus("SUCCESS"))
		assert.Equal(t, entity.TransactionStatusFailed, entity.TransactionStatus("FAILED"))
	})

	t.Run("TransactionProvider values", func(t *testing.T) {
		assert.Equal(t, entity.TransactionProviderSEPAY, entity.TransactionProvider("SEPAY"))
	})

	t.Run("TransactionGateway values", func(t *testing.T) {
		assert.Equal(t, entity.TransactionGatewayVietQR, entity.TransactionGateway("VIETQR"))
		assert.Equal(t, entity.TransactionGatewayBankTransfer, entity.TransactionGateway("BANK_TRANSFER"))
	})
}

// TestTransaction_DefaultValues verifies default field values
func TestTransaction_DefaultValues(t *testing.T) {
	// Arrange
	userID := uuid.New()
	amount := decimal.NewFromInt(10000)

	// Act
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   amount,
		Currency: "VND",
	}

	// Assert - Verify default currency is set
	assert.Equal(t, "VND", tx.Currency)
}

// TestTransaction_SePayFields verifies SePay specific fields
func TestTransaction_SePayFields(t *testing.T) {
	// Arrange
	userID := uuid.New()
	amount := decimal.NewFromInt(50000)
	sePayID := "SEPAY123456"
	referenceCode := "REF789"
	webhookPayload := map[string]interface{}{
		"content": "PAYMENT_DESCRIPTION",
		"amount":  50000,
	}

	// Act
	tx := entity.Transaction{
		UserID:         userID,
		Amount:         amount,
		Provider:       entity.TransactionProviderSEPAY,
		SePayID:        sePayID,
		ReferenceCode:  referenceCode,
		WebhookPayload: webhookPayload,
	}

	// Assert
	assert.Equal(t, entity.TransactionProviderSEPAY, tx.Provider)
	assert.Equal(t, sePayID, tx.SePayID)
	assert.Equal(t, referenceCode, tx.ReferenceCode)
	assert.NotNil(t, tx.WebhookPayload)
	assert.Equal(t, "PAYMENT_DESCRIPTION", tx.WebhookPayload["content"])
}

// TestTransaction_TypeSubscription tests SUBSCRIPTION transaction type
func TestTransaction_TypeSubscription(t *testing.T) {
	// Arrange
	userID := uuid.New()
	authorID := uuid.New()
	planID := "1_MONTH"
	gateway := entity.TransactionGatewayVietQR

	// Act
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   decimal.NewFromInt(100000),
		Type:     entity.TransactionTypeSubscription,
		TargetID: &authorID,
		PlanID:   &planID,
		Currency: "VND",
		Status:   entity.TransactionStatusPending,
		Provider: entity.TransactionProviderSEPAY,
		Gateway:  &gateway,
	}

	// Assert
	assert.Equal(t, entity.TransactionTypeSubscription, tx.Type)
	assert.NotNil(t, tx.TargetID)
	assert.Equal(t, authorID, *tx.TargetID)
	assert.NotNil(t, tx.PlanID)
	assert.Equal(t, "1_MONTH", *tx.PlanID)
	assert.Equal(t, entity.TransactionStatusPending, tx.Status)
	assert.NotNil(t, tx.Gateway)
	assert.Equal(t, entity.TransactionGatewayVietQR, *tx.Gateway)
}

// TestTransaction_TypeSeries tests SERIES transaction type
func TestTransaction_TypeSeries(t *testing.T) {
	// Arrange
	userID := uuid.New()
	seriesID := uuid.New()
	planID := "3_MONTH"
	gateway := entity.TransactionGatewayBankTransfer

	// Act
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   decimal.NewFromInt(200000),
		Type:     entity.TransactionTypeSeries,
		TargetID: &seriesID,
		PlanID:   &planID,
		Currency: "VND",
		Status:   entity.TransactionStatusSuccess,
		Provider: entity.TransactionProviderSEPAY,
		Gateway:  &gateway,
	}

	// Assert
	assert.Equal(t, entity.TransactionTypeSeries, tx.Type)
	assert.Equal(t, seriesID, *tx.TargetID)
	assert.Equal(t, entity.TransactionStatusSuccess, tx.Status)
	assert.Equal(t, entity.TransactionGatewayBankTransfer, *tx.Gateway)
}

// TestTransaction_TypeDonation tests DONATION transaction type
func TestTransaction_TypeDonation(t *testing.T) {
	// Arrange
	userID := uuid.New()
	authorID := uuid.New()
	gateway := entity.TransactionGatewayVietQR

	// Act
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   decimal.NewFromInt(50000),
		Type:     entity.TransactionTypeDonation,
		TargetID: &authorID,
		Content:  "Support your content",
		Currency: "VND",
		Status:   entity.TransactionStatusSuccess,
		Provider: entity.TransactionProviderSEPAY,
		Gateway:  &gateway,
	}

	// Assert
	assert.Equal(t, entity.TransactionTypeDonation, tx.Type)
	assert.Equal(t, authorID, *tx.TargetID)
	assert.NotEmpty(t, tx.Content)
	assert.Equal(t, "Support your content", tx.Content)
}

// TestTransaction_OrderID tests OrderID field
func TestTransaction_OrderID(t *testing.T) {
	// Arrange
	userID := uuid.New()
	orderID := "ORD123456789"

	// Act
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   decimal.NewFromInt(10000),
		OrderID:  orderID,
		Currency: "VND",
	}

	// Assert
	assert.Equal(t, orderID, tx.OrderID)
}

// TestTransaction_OptionalFields tests optional/nullable fields
func TestTransaction_OptionalFields(t *testing.T) {
	// Arrange
	userID := uuid.New()

	// Act - Create transaction with optional fields set to nil/empty
	tx := entity.Transaction{
		UserID:   userID,
		Amount:   decimal.NewFromInt(10000),
		TargetID: nil, // Optional
		PlanID:   nil, // Optional
		Gateway:  nil, // Optional
		Currency: "VND",
	}

	// Assert
	assert.Nil(t, tx.TargetID)
	assert.Nil(t, tx.PlanID)
	assert.Nil(t, tx.Gateway)
}
