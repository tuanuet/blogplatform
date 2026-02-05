package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBlogVersionRepository_Create_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()
	editorID := uuid.New()
	categoryID := uuid.New()
	version := &entity.BlogVersion{
		ID:            uuid.New(),
		BlogID:        blogID,
		VersionNumber: 1,
		Title:         "Test Blog Version",
		Slug:          "test-blog-version",
		Content:       "Content of the version",
		Status:        entity.BlogStatusDraft,
		Visibility:    entity.BlogVisibilityPublic,
		CategoryID:    &categoryID,
		EditorID:      editorID,
		CreatedAt:     time.Now(),
		Tags: []entity.Tag{
			{ID: uuid.New(), Name: "Go"},
			{ID: uuid.New(), Name: "Testing"},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "blog_versions"`)).
		WithArgs(
			version.BlogID,
			version.VersionNumber,
			version.Title,
			version.Slug,
			sqlmock.AnyArg(), // Excerpt
			version.Content,
			sqlmock.AnyArg(), // ThumbnailURL
			version.Status,
			version.Visibility,
			version.CategoryID,
			version.EditorID,
			sqlmock.AnyArg(), // ChangeSummary
			version.ID,
			sqlmock.AnyArg(), // CreatedAt
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(version.ID))

	// Expect tags to be upserted
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(version.Tags[0].ID, time.Now(), time.Now()).
			AddRow(version.Tags[1].ID, time.Now(), time.Now()))

	// Expect tag associations to be inserted
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "blog_version_tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"blog_version_id", "tag_id"}).
			AddRow(version.ID, version.Tags[0].ID).
			AddRow(version.ID, version.Tags[1].ID))

	mock.ExpectCommit()

	err := repo.Create(ctx, version)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_FindByID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	id := uuid.New()
	blogID := uuid.New()
	editorID := uuid.New()
	categoryID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "blog_id", "version_number", "title", "editor_id", "category_id", "created_at"}).
		AddRow(id, blogID, 1, "Test Title", editorID, categoryID, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blog_versions" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	// Expect Preload Category
	categoryRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(categoryID, "Test Category")
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
		WithArgs(categoryID).
		WillReturnRows(categoryRows)

	// Expect Preload Editor
	editorRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(editorID, "Test Editor")
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(editorID).
		WillReturnRows(editorRows)

	// Expect Preload Tags (Many2Many)
	// GORM splits this into two queries:
	// 1. Get join table entries
	// 2. Get tags

	tagID1 := uuid.New()
	tagID2 := uuid.New()

	joinRows := sqlmock.NewRows([]string{"blog_version_id", "tag_id"}).
		AddRow(id, tagID1).
		AddRow(id, tagID2)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blog_version_tags" WHERE "blog_version_tags"."blog_version_id" = $1`)).
		WithArgs(id).
		WillReturnRows(joinRows)

	tagRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(tagID1, "Tag1").
		AddRow(tagID2, "Tag2")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tags" WHERE "tags"."id" IN ($1,$2)`)).
		WithArgs(tagID1, tagID2).
		WillReturnRows(tagRows)

	version, err := repo.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, version)
	assert.Equal(t, id, version.ID)
	assert.Equal(t, editorID, version.Editor.ID)
	assert.Equal(t, categoryID, version.Category.ID)
	assert.Len(t, version.Tags, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blog_versions" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	version, err := repo.FindByID(ctx, id)

	assert.NoError(t, err) // Should return nil, nil
	assert.Nil(t, version)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_FindByBlogID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()
	pagination := repository.Pagination{Page: 1, PageSize: 10}

	// 1. Count query
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blog_versions" WHERE blog_id = $1`)).
		WithArgs(blogID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// 2. Data query
	editorID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "blog_id", "version_number", "editor_id", "created_at"}).
		AddRow(uuid.New(), blogID, 2, editorID, time.Now()).
		AddRow(uuid.New(), blogID, 1, editorID, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blog_versions" WHERE blog_id = $1 ORDER BY version_number DESC LIMIT $2`)).
		WithArgs(blogID, 10).
		WillReturnRows(rows)

	// 3. Preload Editor
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(editorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(editorID, "Editor Name"))

	result, err := repo.FindByBlogID(ctx, blogID, pagination)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Data[0].VersionNumber)
	assert.Equal(t, 1, result.Data[1].VersionNumber)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_GetNextVersionNumber_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT COALESCE(MAX(version_number), 0) + 1 FROM "blog_versions" WHERE blog_id = $1`)).
		WithArgs(blogID).
		WillReturnRows(rows)

	nextVersion, err := repo.GetNextVersionNumber(ctx, blogID)

	assert.NoError(t, err)
	assert.Equal(t, 5, nextVersion)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_Delete_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "blog_versions" WHERE "blog_versions"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_CountByBlogID_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(10)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blog_versions" WHERE blog_id = $1`)).
		WithArgs(blogID).
		WillReturnRows(rows)

	count, err := repo.CountByBlogID(ctx, blogID)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_DeleteOldest_Success(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()
	keep := 5

	// 1. Count
	rows := sqlmock.NewRows([]string{"count"}).AddRow(12)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blog_versions" WHERE blog_id = $1`)).
		WithArgs(blogID).
		WillReturnRows(rows)

	// 2. Delete (12 - 5 = 7 items)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "blog_versions" WHERE id IN (SELECT id FROM "blog_versions" WHERE blog_id = $1 ORDER BY version_number ASC LIMIT $2)`)).
		WithArgs(blogID, 7).
		WillReturnResult(sqlmock.NewResult(0, 7))
	mock.ExpectCommit()

	err := repo.DeleteOldest(ctx, blogID, keep)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogVersionRepository_DeleteOldest_NoAction(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogVersionRepository(db)

	ctx := context.Background()
	blogID := uuid.New()
	keep := 10

	// 1. Count
	rows := sqlmock.NewRows([]string{"count"}).AddRow(5) // Total 5, Keep 10. No delete.
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blog_versions" WHERE blog_id = $1`)).
		WithArgs(blogID).
		WillReturnRows(rows)

	err := repo.DeleteOldest(ctx, blogID, keep)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
