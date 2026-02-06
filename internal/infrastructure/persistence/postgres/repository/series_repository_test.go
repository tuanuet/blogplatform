package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupSeriesTestDB creates a mock database using sqlmock
func setupSeriesTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestSeriesRepository_GetHighlighted_Success(t *testing.T) {
	// Arrange
	db, mock := setupSeriesTestDB(t)
	repo := NewSeriesRepository(db)

	ctx := context.Background()
	limit := 10

	seriesID := uuid.New()
	authorID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "description", "author_id", "created_at",
		"author_name", "author_avatar_url",
		"subscriber_count", "blog_count",
	}).AddRow(
		seriesID, "Test Series", "test-series", "Description", authorID, now,
		"Author Name", "http://avatar.url",
		100, 5,
	)

	query := `SELECT
          s.id, s.title, s.slug, s.description, s.author_id, s.created_at,
          u.name as author_name, u.avatar_url as author_avatar_url,
          COALESCE(sub.subscriber_count, 0) as subscriber_count,
          COALESCE(bc.blog_count, 0) as blog_count
      FROM series s
      LEFT JOIN users u ON s.author_id = u.id
      LEFT JOIN (
          SELECT series_id, COUNT(*) as subscriber_count
          FROM user_series_purchases
          GROUP BY series_id
      ) sub ON s.id = sub.series_id
      LEFT JOIN (
          SELECT sb.series_id, COUNT(*) as blog_count
          FROM series_blogs sb
          JOIN blogs b ON sb.blog_id = b.id
          WHERE b.deleted_at IS NULL
          GROUP BY sb.series_id
      ) bc ON s.id = bc.series_id
      WHERE s.deleted_at IS NULL
      ORDER BY subscriber_count DESC
      LIMIT $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(limit).
		WillReturnRows(rows)

	// Act
	result, err := repo.GetHighlighted(ctx, limit)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, seriesID, result[0].Series.ID)
	assert.Equal(t, "Test Series", result[0].Series.Title)
	assert.Equal(t, authorID, result[0].Author.ID)
	assert.Equal(t, "Author Name", result[0].Author.Name)
	assert.Equal(t, "http://avatar.url", *result[0].Author.AvatarURL)
	assert.Equal(t, 100, result[0].SubscriberCount)
	assert.Equal(t, 5, result[0].BlogCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSeriesRepository_GetHighlighted_Error(t *testing.T) {
	// Arrange
	db, mock := setupSeriesTestDB(t)
	repo := NewSeriesRepository(db)

	ctx := context.Background()
	limit := 10

	query := `SELECT
          s.id, s.title, s.slug, s.description, s.author_id, s.created_at,
          u.name as author_name, u.avatar_url as author_avatar_url,
          COALESCE(sub.subscriber_count, 0) as subscriber_count,
          COALESCE(bc.blog_count, 0) as blog_count
      FROM series s
      LEFT JOIN users u ON s.author_id = u.id
      LEFT JOIN (
          SELECT series_id, COUNT(*) as subscriber_count
          FROM user_series_purchases
          GROUP BY series_id
      ) sub ON s.id = sub.series_id
      LEFT JOIN (
          SELECT sb.series_id, COUNT(*) as blog_count
          FROM series_blogs sb
          JOIN blogs b ON sb.blog_id = b.id
          WHERE b.deleted_at IS NULL
          GROUP BY sb.series_id
      ) bc ON s.id = bc.series_id
      WHERE s.deleted_at IS NULL
      ORDER BY subscriber_count DESC
      LIMIT $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(limit).
		WillReturnError(assert.AnError)

	// Act
	result, err := repo.GetHighlighted(ctx, limit)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
