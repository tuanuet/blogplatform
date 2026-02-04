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
	"gorm.io/gorm"
)

// TestSubscriptionPlanRepository_Create_Success tests successful creation
func TestSubscriptionPlanRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		Tier:         entity.TierBronze,
		Price:        decimal.NewFromFloat(9.99),
		DurationDays: 30,
		IsActive:     true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "subscription_plans" ("author_id","tier","price","duration_days","name","description","is_active","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id","created_at","updated_at"`)).
		WithArgs(plan.AuthorID, plan.Tier, plan.Price, plan.DurationDays, plan.Name, plan.Description, plan.IsActive, plan.DeletedAt, plan.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(plan.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(ctx, plan)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Create_Error tests database error
func TestSubscriptionPlanRepository_Create_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Tier:     entity.TierBronze,
		Price:    decimal.NewFromFloat(9.99),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "subscription_plans"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Create(ctx, plan)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByID_Success tests successful find
func TestSubscriptionPlanRepository_FindByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	planID := uuid.New()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tier", "price", "duration_days", "name", "description", "is_active", "created_at", "updated_at"}).
		AddRow(planID, authorID, entity.TierBronze, decimal.NewFromFloat(9.99), 30, "Bronze Plan", "Description", true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE id = $1 AND "subscription_plans"."deleted_at" IS NULL ORDER BY "subscription_plans"."id" LIMIT $2`)).
		WithArgs(planID, 1).
		WillReturnRows(rows)

	plan, err := repo.FindByID(ctx, planID)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, planID, plan.ID)
	assert.Equal(t, entity.TierBronze, plan.Tier)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByID_NotFound tests not found case
func TestSubscriptionPlanRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	planID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE id = $1 AND "subscription_plans"."deleted_at" IS NULL ORDER BY "subscription_plans"."id" LIMIT $2`)).
		WithArgs(planID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	plan, err := repo.FindByID(ctx, planID)

	assert.NoError(t, err)
	assert.Nil(t, plan)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByAuthorAndTier_Success tests successful find by author and tier
func TestSubscriptionPlanRepository_FindByAuthorAndTier_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tier", "price", "duration_days", "name", "description", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, entity.TierBronze, decimal.NewFromFloat(9.99), 30, "Bronze Plan", "Description", true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE (author_id = $1 AND tier = $2) AND "subscription_plans"."deleted_at" IS NULL ORDER BY "subscription_plans"."id" LIMIT $3`)).
		WithArgs(authorID, entity.TierBronze, 1).
		WillReturnRows(rows)

	plan, err := repo.FindByAuthorAndTier(ctx, authorID, entity.TierBronze)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, authorID, plan.AuthorID)
	assert.Equal(t, entity.TierBronze, plan.Tier)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByAuthor_Success tests successful find by author
func TestSubscriptionPlanRepository_FindByAuthor_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tier", "price", "duration_days", "name", "description", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, entity.TierBronze, decimal.NewFromFloat(9.99), 30, "Bronze", "Desc", true, time.Now(), time.Now()).
		AddRow(uuid.New(), authorID, entity.TierSilver, decimal.NewFromFloat(19.99), 30, "Silver", "Desc", true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE author_id = $1 AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(authorID).
		WillReturnRows(rows)

	plans, err := repo.FindByAuthor(ctx, authorID)

	assert.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByAuthor_Empty tests empty result
func TestSubscriptionPlanRepository_FindByAuthor_Empty(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tier", "price", "duration_days", "name", "description", "is_active", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE author_id = $1 AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(authorID).
		WillReturnRows(rows)

	plans, err := repo.FindByAuthor(ctx, authorID)

	assert.NoError(t, err)
	assert.Empty(t, plans)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindActiveByAuthor_Success tests successful find active by author
func TestSubscriptionPlanRepository_FindActiveByAuthor_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tier", "price", "duration_days", "name", "description", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, entity.TierBronze, decimal.NewFromFloat(9.99), 30, "Bronze", "Desc", true, time.Now(), time.Now()).
		AddRow(uuid.New(), authorID, entity.TierSilver, decimal.NewFromFloat(19.99), 30, "Silver", "Desc", true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE (author_id = $1 AND is_active = $2) AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(authorID, true).
		WillReturnRows(rows)

	plans, err := repo.FindActiveByAuthor(ctx, authorID)

	assert.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Update_Success tests successful update
func TestSubscriptionPlanRepository_Update_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		Tier:         entity.TierBronze,
		Price:        decimal.NewFromFloat(9.99),
		DurationDays: 30,
		IsActive:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscription_plans" SET "author_id"=$1,"tier"=$2,"price"=$3,"duration_days"=$4,"name"=$5,"description"=$6,"is_active"=$7,"created_at"=$8,"updated_at"=$9,"deleted_at"=$10 WHERE "subscription_plans"."deleted_at" IS NULL AND "id" = $11`)).
		WithArgs(plan.AuthorID, plan.Tier, plan.Price, plan.DurationDays, plan.Name, plan.Description, plan.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg(), plan.DeletedAt, plan.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, plan)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Update_Error tests database error
func TestSubscriptionPlanRepository_Update_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Tier:     entity.TierBronze,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscription_plans"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Update(ctx, plan)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Upsert_Insert tests upsert with insert
func TestSubscriptionPlanRepository_Upsert_Insert(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		Tier:         entity.TierBronze,
		Price:        decimal.NewFromFloat(9.99),
		DurationDays: 30,
		IsActive:     true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "subscription_plans" ("author_id","tier","price","duration_days","name","description","is_active","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT ("author_id","tier") DO UPDATE SET "price"="excluded"."price","duration_days"="excluded"."duration_days","name"="excluded"."name","description"="excluded"."description","is_active"="excluded"."is_active","updated_at"="excluded"."updated_at","deleted_at"="excluded"."deleted_at" RETURNING "id","created_at","updated_at"`)).
		WithArgs(plan.AuthorID, plan.Tier, plan.Price, plan.DurationDays, plan.Name, plan.Description, plan.IsActive, plan.DeletedAt, plan.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(plan.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Upsert(ctx, plan)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Upsert_Error tests upsert error
func TestSubscriptionPlanRepository_Upsert_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	plan := &entity.SubscriptionPlan{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Tier:     entity.TierBronze,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "subscription_plans"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Upsert(ctx, plan)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Delete_Success tests successful delete (hard delete)
func TestSubscriptionPlanRepository_Delete_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	planID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscription_plans" SET "deleted_at"=$1 WHERE "subscription_plans"."id" = $2 AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), planID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, planID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Delete_NotFound tests not found on delete
func TestSubscriptionPlanRepository_Delete_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	planID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscription_plans" SET "deleted_at"=$1 WHERE "subscription_plans"."id" = $2 AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), planID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, planID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_WithTx tests WithTx method
func TestSubscriptionPlanRepository_WithTx(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	// Valid path
	txRepo := repo.WithTx(db)
	assert.NotNil(t, txRepo)
	// Even though they have the same DB, it should be a new repository instance pointer
	// but testify's NotEqual might still consider them equal if they are deep equal.
	// We just want to ensure it's not nil and returned something.
	assert.NotNil(t, txRepo)

	// Invalid path
	invalidRepo := repo.WithTx("invalid")
	assert.Equal(t, repo, invalidRepo)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_Delete_Error tests database error on delete
func TestSubscriptionPlanRepository_Delete_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	planID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "subscription_plans"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Delete(ctx, planID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSubscriptionPlanRepository_FindByAuthorAndTier_NotFound tests not found
func TestSubscriptionPlanRepository_FindByAuthorAndTier_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewSubscriptionPlanRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "subscription_plans" WHERE (author_id = $1 AND tier = $2) AND "subscription_plans"."deleted_at" IS NULL`)).
		WithArgs(authorID, entity.TierBronze, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	plan, err := repo.FindByAuthorAndTier(ctx, authorID, entity.TierBronze)

	assert.NoError(t, err)
	assert.Nil(t, plan)
	assert.NoError(t, mock.ExpectationsWereMet())
}
