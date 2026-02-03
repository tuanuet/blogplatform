package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUserSeriesPurchaseRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()
	amount := decimal.NewFromFloat(9.99)

	purchase := &entity.UserSeriesPurchase{
		UserID:   userID,
		SeriesID: seriesID,
		Amount:   amount,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "user_series_purchases" ("user_id","series_id","amount") VALUES ($1,$2,$3) RETURNING "created_at"`)).
		WithArgs(userID, seriesID, amount).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()))
	mock.ExpectCommit()

	err := repo.Create(ctx, purchase)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_HasPurchased_True(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "user_series_purchases" WHERE user_id = $1 AND series_id = $2`)).
		WithArgs(userID, seriesID).
		WillReturnRows(rows)

	hasPurchased, err := repo.HasPurchased(ctx, userID, seriesID)

	assert.NoError(t, err)
	assert.True(t, hasPurchased)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_HasPurchased_False(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "user_series_purchases" WHERE user_id = $1 AND series_id = $2`)).
		WithArgs(userID, seriesID).
		WillReturnRows(rows)

	hasPurchased, err := repo.HasPurchased(ctx, userID, seriesID)

	assert.NoError(t, err)
	assert.False(t, hasPurchased)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_HasPurchased_Failure(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "user_series_purchases" WHERE user_id = $1 AND series_id = $2`)).
		WithArgs(userID, seriesID).
		WillReturnError(assert.AnError)

	hasPurchased, err := repo.HasPurchased(ctx, userID, seriesID)

	assert.Error(t, err)
	assert.False(t, hasPurchased)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_Create_Failure(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()
	amount := decimal.NewFromFloat(9.99)

	purchase := &entity.UserSeriesPurchase{
		UserID:   userID,
		SeriesID: seriesID,
		Amount:   amount,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "user_series_purchases" ("user_id","series_id","amount") VALUES ($1,$2,$3) RETURNING "created_at"`)).
		WithArgs(userID, seriesID, amount).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Create(ctx, purchase)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_GetUserPurchases_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	seriesID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "series_id", "amount", "created_at"}).
		AddRow(userID, seriesID, decimal.NewFromFloat(9.99), time.Now())

	// Expect Preload("Series")
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_series_purchases" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	// Mock for Preload("Series") - GORM usually does this in a separate query
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "series" WHERE "series"."id" = $1`)).
		WithArgs(seriesID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(seriesID, "Test Series"))

	purchases, err := repo.GetUserPurchases(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, purchases, 1)
	assert.Equal(t, seriesID, purchases[0].SeriesID)
	assert.NotNil(t, purchases[0].Series)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_GetUserPurchases_Empty(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_series_purchases" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "series_id", "amount", "created_at"}))

	purchases, err := repo.GetUserPurchases(ctx, userID)

	assert.NoError(t, err)
	assert.Empty(t, purchases)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSeriesPurchaseRepository_GetUserPurchases_Failure(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserSeriesPurchaseRepository(db)

	ctx := context.Background()
	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user_series_purchases" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnError(assert.AnError)

	purchases, err := repo.GetUserPurchases(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, purchases)
	assert.NoError(t, mock.ExpectationsWereMet())
}
