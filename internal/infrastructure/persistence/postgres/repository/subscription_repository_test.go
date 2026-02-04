package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionRepository_UpdateExpiry_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	tier := "PREMIUM"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscriptions" SET "expires_at"=$1,"tier"=$2,"updated_at"=$3 WHERE subscriber_id = $4 AND author_id = $5`)).
		WithArgs(expiresAt, tier, sqlmock.AnyArg(), userID, authorID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateExpiry(ctx, userID, authorID, expiresAt, tier)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_UpdateExpiry_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	tier := "PREMIUM"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscriptions" SET "expires_at"=$1,"tier"=$2,"updated_at"=$3 WHERE subscriber_id = $4 AND author_id = $5`)).
		WithArgs(expiresAt, tier, sqlmock.AnyArg(), userID, authorID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.UpdateExpiry(ctx, userID, authorID, expiresAt, tier)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_UpdateExpiry_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	tier := "PREMIUM"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscriptions" SET "expires_at"=$1,"tier"=$2,"updated_at"=$3 WHERE subscriber_id = $4 AND author_id = $5`)).
		WithArgs(expiresAt, tier, sqlmock.AnyArg(), userID, authorID).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.UpdateExpiry(ctx, userID, authorID, expiresAt, tier)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_FindActiveSubscription_Paid(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	rows := sqlmock.NewRows([]string{"id", "subscriber_id", "author_id", "expires_at", "tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), userID, authorID, expiresAt, "PREMIUM", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscriptions" WHERE (subscriber_id = $1 AND author_id = $2) AND (expires_at > $3 OR expires_at IS NULL) ORDER BY "subscriptions"."id" LIMIT $4`)).
		WithArgs(userID, authorID, sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	result, err := repo.FindActiveSubscription(ctx, userID, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.SubscriberID)
	assert.Equal(t, authorID, result.AuthorID)
	assert.Equal(t, "PREMIUM", result.Tier)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_FindActiveSubscription_SilverTier(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(time.Hour)

	rows := sqlmock.NewRows([]string{"id", "subscriber_id", "author_id", "expires_at", "tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), userID, authorID, expiresAt, "SILVER", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscriptions" WHERE (subscriber_id = $1 AND author_id = $2) AND (expires_at > $3 OR expires_at IS NULL) ORDER BY "subscriptions"."id" LIMIT $4`)).
		WithArgs(userID, authorID, sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	result, err := repo.FindActiveSubscription(ctx, userID, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "SILVER", result.Tier)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	subscription := &entity.Subscription{
		SubscriberID: userID,
		AuthorID:     authorID,
		ExpiresAt:    &expiresAt,
		Tier:         "SILVER",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "subscriptions" ("subscriber_id","author_id","expires_at","tier") VALUES ($1,$2,$3,$4) RETURNING "id","created_at","updated_at"`)).
		WithArgs(userID, authorID, &expiresAt, "SILVER").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(uuid.New(), time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(ctx, subscription)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_FindActiveSubscription_Free(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "subscriber_id", "author_id", "expires_at", "tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), userID, authorID, nil, "FREE", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscriptions" WHERE (subscriber_id = $1 AND author_id = $2) AND (expires_at > $3 OR expires_at IS NULL) ORDER BY "subscriptions"."id" LIMIT $4`)).
		WithArgs(userID, authorID, sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	result, err := repo.FindActiveSubscription(ctx, userID, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.ExpiresAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubscriptionRepository_FindActiveSubscription_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionRepository(db)

	ctx := context.Background()
	userID := uuid.New()
	authorID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscriptions" WHERE (subscriber_id = $1 AND author_id = $2) AND (expires_at > $3 OR expires_at IS NULL) ORDER BY "subscriptions"."id" LIMIT $4`)).
		WithArgs(userID, authorID, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	result, err := repo.FindActiveSubscription(ctx, userID, authorID)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
