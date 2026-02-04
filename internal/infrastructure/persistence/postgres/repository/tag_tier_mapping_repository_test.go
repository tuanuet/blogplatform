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
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// testTime is a helper variable for test time values
var testTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// TestTagTierMappingRepository_Create_Success tests successful creation
func TestTagTierMappingRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		TagID:        uuid.New(),
		RequiredTier: entity.TierBronze,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tag_tier_mappings" ("author_id","tag_id","required_tier","id") VALUES ($1,$2,$3,$4) RETURNING "id","created_at","updated_at"`)).
		WithArgs(mapping.AuthorID, mapping.TagID, mapping.RequiredTier, mapping.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(mapping.ID, testTime, testTime))
	mock.ExpectCommit()

	err := repo.Create(ctx, mapping)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Create_Error tests database error
func TestTagTierMappingRepository_Create_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		TagID:    uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tag_tier_mappings"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Create(ctx, mapping)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByID_Success tests successful find
func TestTagTierMappingRepository_FindByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	id := uuid.New()
	authorID := uuid.New()
	tagID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tag_id", "required_tier", "created_at", "updated_at"}).
		AddRow(id, authorID, tagID, entity.TierBronze, testTime, testTime)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE id = $1 ORDER BY "tag_tier_mappings"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	mapping, err := repo.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, id, mapping.ID)
	assert.Equal(t, entity.TierBronze, mapping.RequiredTier)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByID_NotFound tests not found case
func TestTagTierMappingRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE id = $1 ORDER BY "tag_tier_mappings"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	mapping, err := repo.FindByID(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, mapping)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByAuthorAndTag_Success tests successful find by author and tag
func TestTagTierMappingRepository_FindByAuthorAndTag_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tag_id", "required_tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, tagID, entity.TierBronze, testTime, testTime)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id = $2 ORDER BY "tag_tier_mappings"."id" LIMIT $3`)).
		WithArgs(authorID, tagID, 1).
		WillReturnRows(rows)

	mapping, err := repo.FindByAuthorAndTag(ctx, authorID, tagID)

	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, authorID, mapping.AuthorID)
	assert.Equal(t, tagID, mapping.TagID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByAuthorAndTag_Error tests find by author and tag error
func TestTagTierMappingRepository_FindByAuthorAndTag_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id = $2`)).
		WithArgs(authorID, tagID, 1).
		WillReturnError(assert.AnError)

	mapping, err := repo.FindByAuthorAndTag(ctx, authorID, tagID)

	assert.Error(t, err)
	assert.Nil(t, mapping)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByAuthor_Success tests successful find by author
func TestTagTierMappingRepository_FindByAuthor_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "author_id", "tag_id", "required_tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, uuid.New(), entity.TierBronze, testTime, testTime).
		AddRow(uuid.New(), authorID, uuid.New(), entity.TierSilver, testTime, testTime)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE author_id = $1`)).
		WithArgs(authorID).
		WillReturnRows(rows)

	mappings, err := repo.FindByAuthor(ctx, authorID)

	assert.NoError(t, err)
	assert.Len(t, mappings, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByTagIDs_Success tests successful find by tag IDs
func TestTagTierMappingRepository_FindByTagIDs_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagIDs := []uuid.UUID{uuid.New(), uuid.New()}

	rows := sqlmock.NewRows([]string{"id", "author_id", "tag_id", "required_tier", "created_at", "updated_at"}).
		AddRow(uuid.New(), authorID, tagIDs[0], entity.TierBronze, testTime, testTime).
		AddRow(uuid.New(), authorID, tagIDs[1], entity.TierSilver, testTime, testTime)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id IN ($2,$3)`)).
		WithArgs(authorID, tagIDs[0], tagIDs[1]).
		WillReturnRows(rows)

	mappings, err := repo.FindByTagIDs(ctx, authorID, tagIDs)

	assert.NoError(t, err)
	assert.Len(t, mappings, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_FindByTagIDs_Empty tests with empty tag IDs
func TestTagTierMappingRepository_FindByTagIDs_Empty(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagIDs := []uuid.UUID{}

	rows := sqlmock.NewRows([]string{"id", "author_id", "tag_id", "required_tier", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag_tier_mappings" WHERE author_id = $1`)).
		WithArgs(authorID).
		WillReturnRows(rows)

	mappings, err := repo.FindByTagIDs(ctx, authorID, tagIDs)

	assert.NoError(t, err)
	assert.Empty(t, mappings)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Update_Success tests successful update
func TestTagTierMappingRepository_Update_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		TagID:        uuid.New(),
		RequiredTier: entity.TierBronze,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "tag_tier_mappings" SET "author_id"=$1,"tag_id"=$2,"required_tier"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`)).
		WithArgs(mapping.AuthorID, mapping.TagID, mapping.RequiredTier, sqlmock.AnyArg(), sqlmock.AnyArg(), mapping.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, mapping)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Update_Error tests database error
func TestTagTierMappingRepository_Update_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		TagID:    uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "tag_tier_mappings"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Update(ctx, mapping)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Upsert_Insert tests upsert with insert
func TestTagTierMappingRepository_Upsert_Insert(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:           uuid.New(),
		AuthorID:     uuid.New(),
		TagID:        uuid.New(),
		RequiredTier: entity.TierBronze,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tag_tier_mappings" ("author_id","tag_id","required_tier","id") VALUES ($1,$2,$3,$4) ON CONFLICT ("author_id","tag_id") DO UPDATE SET "required_tier"="excluded"."required_tier","updated_at"="excluded"."updated_at" RETURNING "id","created_at","updated_at"`)).
		WithArgs(mapping.AuthorID, mapping.TagID, mapping.RequiredTier, mapping.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(mapping.ID, testTime, testTime))
	mock.ExpectCommit()

	err := repo.Upsert(ctx, mapping)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Upsert_Error tests upsert error
func TestTagTierMappingRepository_Upsert_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	mapping := &entity.TagTierMapping{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		TagID:    uuid.New(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tag_tier_mappings"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Upsert(ctx, mapping)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Delete_Success tests successful delete (hard delete)
func TestTagTierMappingRepository_Delete_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id = $2`)).
		WithArgs(authorID, tagID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, authorID, tagID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Delete_NotFound tests not found on delete
func TestTagTierMappingRepository_Delete_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id = $2`)).
		WithArgs(authorID, tagID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, authorID, tagID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_Delete_Error tests database error on delete
func TestTagTierMappingRepository_Delete_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE author_id = $1 AND tag_id = $2`)).
		WithArgs(authorID, tagID).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Delete(ctx, authorID, tagID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_DeleteByID_Success tests successful delete by ID
func TestTagTierMappingRepository_DeleteByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE "tag_tier_mappings"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(ctx, id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_DeleteByID_NotFound tests not found on delete by ID
func TestTagTierMappingRepository_DeleteByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE "tag_tier_mappings"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.DeleteByID(ctx, id)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_DeleteByID_Error tests database error on delete by ID
func TestTagTierMappingRepository_DeleteByID_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag_tier_mappings" WHERE "tag_tier_mappings"."id" = $1`)).
		WithArgs(id).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.DeleteByID(ctx, id)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_CountBlogsByTagAndAuthor_Success tests successful count
func TestTagTierMappingRepository_CountBlogsByTagAndAuthor_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blogs" JOIN blog_tags ON blogs.id = blog_tags.blog_id WHERE (blogs.author_id = $1 AND blog_tags.tag_id = $2) AND "blogs"."deleted_at" IS NULL`)).
		WithArgs(authorID, tagID).
		WillReturnRows(rows)

	count, err := repo.CountBlogsByTagAndAuthor(ctx, authorID, tagID)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_CountBlogsByTagAndAuthor_Error tests error
func TestTagTierMappingRepository_CountBlogsByTagAndAuthor_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT count(*) FROM \"blogs\"")).
		WillReturnError(assert.AnError)

	count, err := repo.CountBlogsByTagAndAuthor(ctx, authorID, tagID)

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_WithTx tests WithTx method
func TestTagTierMappingRepository_WithTx(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock tx: %v", err)
	}
	defer sqlDB.Close()

	gormDB, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	txRepo := repo.WithTx(gormDB)

	assert.NotNil(t, txRepo)
	assert.NotEqual(t, repo, txRepo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestTagTierMappingRepository_WithTx_NonGorm tests WithTx with non-gorm db
func TestTagTierMappingRepository_WithTx_NonGorm(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := NewTagTierMappingRepository(db)

	txRepo := repo.WithTx("not a gorm db")

	assert.Equal(t, repo, txRepo)
}
