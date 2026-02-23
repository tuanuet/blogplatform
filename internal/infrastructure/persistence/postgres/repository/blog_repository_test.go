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
)

func TestBlogRepository_FindTopViewed(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogRepository(db)

	ctx := context.Background()
	publishedBefore := time.Now()
	pagination := repository.Pagination{Page: 1, PageSize: 10}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(*) FROM "blogs" WHERE deleted_at IS NULL AND status = $1 AND visibility = $2 AND published_at <= $3`)).
		WithArgs(entity.BlogStatusPublished, entity.BlogVisibilityPublic, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "blogs" WHERE deleted_at IS NULL AND status = $1 AND visibility = $2 AND published_at <= $3 AND "blogs"."deleted_at" IS NULL ORDER BY view_count DESC,published_at DESC LIMIT $4`)).
		WithArgs(entity.BlogStatusPublished, entity.BlogVisibilityPublic, sqlmock.AnyArg(), 10).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	result, err := repo.FindTopViewed(ctx, pagination, publishedBefore)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, 0, result.TotalPages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlogRepository_IncrementViewCount(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewBlogRepository(db)

	ctx := context.Background()
	blogID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "blogs" SET "view_count"=view_count + 1,"updated_at"=$1 WHERE (id = $2 AND deleted_at IS NULL) AND "blogs"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), blogID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.IncrementViewCount(ctx, blogID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
