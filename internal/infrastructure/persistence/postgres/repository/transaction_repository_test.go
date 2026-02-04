package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// JSON is a custom type to handle JSON columns
type JSON map[string]interface{}

// Value implements driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// setupTestDB creates a mock database using sqlmock
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	return gormDB, mock
}

// TestTransactionRepository_Create_Success tests successful transaction creation
func TestTransactionRepository_Create_Success(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	amount := decimal.NewFromInt(100000)

	tx := &entity.Transaction{
		UserID:        userID,
		Amount:        amount,
		Currency:      "VND",
		Provider:      entity.TransactionProviderSEPAY,
		Type:          entity.TransactionTypeSubscription,
		Status:        entity.TransactionStatusPending,
		SePayID:       "SEPAY123456",
		ReferenceCode: "REF789",
		OrderID:       "ORD123456789",
	}

	// Expect GORM to insert the transaction
	// Column order: user_id, amount, currency, provider, gateway, type, status, target_id, plan_id, content, se_pay_id, reference_code, order_id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "transactions"`)).
		WithArgs(
			userID,
			amount,
			"VND",
			entity.TransactionProviderSEPAY,
			nil, // gateway
			entity.TransactionTypeSubscription,
			entity.TransactionStatusPending,
			nil, // target_id
			nil, // plan_id
			"",  // content
			"SEPAY123456",
			"REF789",
			"ORD123456789",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "webhook_payload", "created_at", "updated_at"}).
			AddRow(uuid.New(), []byte("{}"), time.Now(), time.Now()))
	mock.ExpectCommit()

	// Act
	err := repo.Create(ctx, tx)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_Create_Error tests transaction creation failure
func TestTransactionRepository_Create_Error(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	amount := decimal.NewFromInt(100000)

	tx := &entity.Transaction{
		UserID: userID,
		Amount: amount,
	}

	// Expect GORM to return an error
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "transactions"`)).
		WithArgs(
			userID,
			amount,
			"VND",     // default currency
			"",        // default provider
			nil,       // gateway
			"",        // type
			"PENDING", // default status
			nil,       // target_id
			nil,       // plan_id
			"",        // content
			"",        // se_pay_id
			"",        // reference_code
			"",        // order_id
		).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	// Act
	err := repo.Create(ctx, tx)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindByRefID_Success tests finding transaction by reference code
func TestTransactionRepository_FindByRefID_Success(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	refID := "REF789"
	userID := uuid.New()

	// Expect GORM to query by reference_code
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "user_id", "amount", "currency",
		"provider", "gateway", "type", "status", "target_id", "plan_id",
		"content", "sepay_id", "reference_code", "order_id", "webhook_payload",
	}).AddRow(
		uuid.New(),
		time.Now(),
		time.Now(),
		nil, // deleted_at (not deleted)
		userID,
		"100000",
		"VND",
		entity.TransactionProviderSEPAY,
		nil,
		entity.TransactionTypeSubscription,
		entity.TransactionStatusPending,
		nil,
		nil,
		"",
		"SEPAY123456",
		refID,
		"ORD123456789", // order_id
		[]byte("{}"),   // webhook_payload as JSON bytes
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE reference_code = $1 ORDER BY "transactions"."id" LIMIT $2`)).
		WithArgs(refID, 1).
		WillReturnRows(rows)

	// Act
	result, err := repo.FindByRefID(ctx, refID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, refID, result.ReferenceCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindByRefID_NotFound tests finding non-existent transaction by reference code
func TestTransactionRepository_FindByRefID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	refID := "NONEXISTENT"

	// Expect GORM to return record not found error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE reference_code = $1 ORDER BY "transactions"."id" LIMIT $2`)).
		WithArgs(refID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.FindByRefID(ctx, refID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindBySePayID_Success tests finding transaction by SePay ID
func TestTransactionRepository_FindBySePayID_Success(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	sePayID := "SEPAY123456"
	userID := uuid.New()

	// Expect GORM to query by sepay_id
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "user_id", "amount", "currency",
		"provider", "gateway", "type", "status", "target_id", "plan_id",
		"content", "sepay_id", "reference_code", "order_id", "webhook_payload",
	}).AddRow(
		uuid.New(),
		time.Now(),
		time.Now(),
		nil, // deleted_at (not deleted)
		userID,
		"100000",
		"VND",
		entity.TransactionProviderSEPAY,
		nil,
		entity.TransactionTypeSubscription,
		entity.TransactionStatusPending,
		nil,
		nil,
		"",
		sePayID,
		"REF789",
		"ORD123456789", // order_id
		[]byte("{}"),   // webhook_payload as JSON bytes
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE sepay_id = $1 ORDER BY "transactions"."id" LIMIT $2`)).
		WithArgs(sePayID, 1).
		WillReturnRows(rows)

	// Act
	result, err := repo.FindBySePayID(ctx, sePayID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sePayID, result.SePayID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindBySePayID_NotFound tests SePay idempotency - returns nil when not found
func TestTransactionRepository_FindBySePayID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	sePayID := "NONEXISTENT"

	// Expect GORM to return record not found error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE sepay_id = $1 ORDER BY "transactions"."id" LIMIT $2`)).
		WithArgs(sePayID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.FindBySePayID(ctx, sePayID)

	// Assert
	// FIXED: Implementation returns nil, nil for not found
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_UpdateStatus_Success tests successful status update
func TestTransactionRepository_UpdateStatus_Success(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	id := uuid.New()
	newStatus := entity.TransactionStatusSuccess

	// Expect GORM to update the status
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "transactions" SET "status"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(newStatus, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := repo.UpdateStatus(ctx, id, newStatus)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_UpdateStatus_Error tests status update failure
func TestTransactionRepository_UpdateStatus_Error(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	id := uuid.New()

	// Expect GORM to return an error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "transactions" SET "status"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(entity.TransactionStatusSuccess, sqlmock.AnyArg(), id).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	// Act
	err := repo.UpdateStatus(ctx, id, entity.TransactionStatusSuccess)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindByUserID_Success tests finding transactions by user ID
func TestTransactionRepository_FindByUserID_Success(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	userID := uuid.New()

	// Expect GORM to query by user_id
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "user_id", "amount", "currency",
		"provider", "gateway", "type", "status", "target_id", "plan_id",
		"content", "sepay_id", "reference_code", "order_id", "webhook_payload",
	}).AddRow(
		uuid.New(),
		time.Now(),
		time.Now(),
		nil, // deleted_at (not deleted)
		userID,
		"100000",
		"VND",
		entity.TransactionProviderSEPAY,
		nil,
		entity.TransactionTypeSubscription,
		entity.TransactionStatusPending,
		nil,
		nil,
		"",
		"SEPAY123456",
		"REF789",
		"ORD123456789", // order_id
		[]byte("{}"),   // webhook_payload as JSON bytes
	).AddRow(
		uuid.New(),
		time.Now().Add(-1*time.Hour),
		time.Now(),
		nil, // deleted_at (not deleted)
		userID,
		"50000",
		"VND",
		entity.TransactionProviderSEPAY,
		nil,
		entity.TransactionTypeDonation,
		entity.TransactionStatusSuccess,
		nil,
		nil,
		"",
		"SEPAY654321",
		"REF999",
		"ORD987654321", // order_id
		[]byte("{}"),   // webhook_payload as JSON bytes
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(rows)

	// Act
	result, err := repo.FindByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindByUserID_NotFound tests finding transactions for user with no transactions
func TestTransactionRepository_FindByUserID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	userID := uuid.New()

	// Expect GORM to return empty result
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "user_id", "amount", "currency",
		"provider", "gateway", "type", "status", "target_id", "plan_id",
		"content", "sepay_id", "reference_code", "order_id", "webhook_payload",
	})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnRows(rows)

	// Act
	result, err := repo.FindByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTransactionRepository_FindByUserID_Error tests error when querying by user ID
func TestTransactionRepository_FindByUserID_Error(t *testing.T) {
	// Arrange
	db, mock := setupTestDB(t)
	repo := NewTransactionRepository(db)

	ctx := context.Background()
	userID := uuid.New()

	// Expect GORM to return an error
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "transactions" WHERE user_id = $1 ORDER BY created_at DESC`)).
		WithArgs(userID).
		WillReturnError(assert.AnError)

	// Act
	result, err := repo.FindByUserID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
