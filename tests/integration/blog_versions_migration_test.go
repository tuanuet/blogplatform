package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestBlogVersionsMigration tests that the blog_versions and blog_version_tags tables
// are created correctly by the migration
func TestBlogVersionsMigration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("blog_versions table exists with correct structure", func(t *testing.T) {
		// Verify table exists by inserting a valid record
		blogID := uuid.New()
		editorID := uuid.New()
		categoryID := uuid.New()

		// Create prerequisite records (blog, user, category)
		createPrerequisites(t, db, blogID, editorID, categoryID)

		// Test that we can insert into blog_versions
		versionID := uuid.New()
		err := db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, excerpt, content,
				thumbnail_url, status, visibility, category_id, editor_id,
				change_summary, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID, blogID, 1, "Test Title", "test-slug", "Test excerpt",
			"Test content", "http://example.com/image.jpg", "draft", "public",
			categoryID, editorID, "Initial version").Error

		require.NoError(t, err, "should be able to insert into blog_versions table")

		// Verify the record exists
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM blog_versions WHERE id = ?", versionID).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "blog_versions record should exist")
	})

	t.Run("blog_version_tags table exists with correct structure", func(t *testing.T) {
		// Create prerequisites
		blogID := uuid.New()
		editorID := uuid.New()
		categoryID := uuid.New()
		tagID := uuid.New()

		createPrerequisites(t, db, blogID, editorID, categoryID)
		createTag(t, db, tagID)

		// Create a version first
		versionID := uuid.New()
		err := db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, content,
				status, visibility, editor_id, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID, blogID, 1, "Test Title", "test-slug",
			"Test content", "draft", "public", editorID).Error
		require.NoError(t, err)

		// Test that we can insert into blog_version_tags
		err = db.Exec(`
			INSERT INTO blog_version_tags (version_id, tag_id)
			VALUES (?, ?)
		`, versionID, tagID).Error

		require.NoError(t, err, "should be able to insert into blog_version_tags table")

		// Verify the record exists
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM blog_version_tags WHERE version_id = ? AND tag_id = ?",
			versionID, tagID).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "blog_version_tags record should exist")
	})

	t.Run("unique constraint on blog_id and version_number", func(t *testing.T) {
		blogID := uuid.New()
		editorID := uuid.New()
		categoryID := uuid.New()

		createPrerequisites(t, db, blogID, editorID, categoryID)

		// Create first version
		versionID1 := uuid.New()
		err := db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, content,
				status, visibility, editor_id, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID1, blogID, 1, "Version 1", "version-1",
			"Content v1", "draft", "public", editorID).Error
		require.NoError(t, err)

		// Try to create another version with the same blog_id and version_number
		versionID2 := uuid.New()
		err = db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, content,
				status, visibility, editor_id, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID2, blogID, 1, "Version 1 Duplicate", "version-1-dup",
			"Content v1 dup", "draft", "public", editorID).Error

		assert.Error(t, err, "should fail due to unique constraint violation")
	})

	t.Run("foreign key constraint on blog_id", func(t *testing.T) {
		nonExistentBlogID := uuid.New()
		editorID := uuid.New()

		createUser(t, db, editorID)

		// Try to create a version with non-existent blog_id
		versionID := uuid.New()
		err := db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, content,
				status, visibility, editor_id, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID, nonExistentBlogID, 1, "Test", "test",
			"Content", "draft", "public", editorID).Error

		assert.Error(t, err, "should fail due to foreign key constraint violation")
	})

	t.Run("cascade delete when blog is deleted", func(t *testing.T) {
		blogID := uuid.New()
		editorID := uuid.New()
		categoryID := uuid.New()

		createPrerequisites(t, db, blogID, editorID, categoryID)

		// Create a version
		versionID := uuid.New()
		err := db.Exec(`
			INSERT INTO blog_versions (
				id, blog_id, version_number, title, slug, content,
				status, visibility, editor_id, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`, versionID, blogID, 1, "Test", "test",
			"Content", "draft", "public", editorID).Error
		require.NoError(t, err)

		// Delete the blog
		err = db.Exec("DELETE FROM blogs WHERE id = ?", blogID).Error
		require.NoError(t, err)

		// Verify the version is also deleted (cascade)
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM blog_versions WHERE id = ?", versionID).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "version should be deleted when blog is deleted")
	})

	t.Run("indexes exist", func(t *testing.T) {
		// Verify indexes exist by checking if queries use them
		// This is a basic check - we verify the index names exist in pg_indexes
		var idxCount int64
		err := db.Raw(`
			SELECT COUNT(*) FROM pg_indexes 
			WHERE tablename = 'blog_versions' 
			AND indexname IN ('idx_blog_versions_blog_id', 'idx_blog_versions_blog_id_created_at')
		`).Scan(&idxCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), idxCount, "both indexes should exist")
	})
}

// Helper functions

func createPrerequisites(t *testing.T, db *gorm.DB, blogID, editorID, categoryID uuid.UUID) {
	t.Helper()
	createUser(t, db, editorID)
	createCategory(t, db, categoryID)
	createBlog(t, db, blogID, editorID, categoryID)
}

func createUser(t *testing.T, db *gorm.DB, userID uuid.UUID) {
	t.Helper()
	err := db.Exec(`
		INSERT INTO users (id, email, name, password_hash, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, true, NOW(), NOW())
	`, userID, userID.String()+"@test.com", "Test User", "hash").Error
	require.NoError(t, err, "failed to create user")
}

func createCategory(t *testing.T, db *gorm.DB, categoryID uuid.UUID) {
	t.Helper()
	err := db.Exec(`
		INSERT INTO categories (id, name, slug, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, categoryID, "Test Category "+categoryID.String(), "test-category-"+categoryID.String()).Error
	require.NoError(t, err, "failed to create category")
}

func createBlog(t *testing.T, db *gorm.DB, blogID, authorID, categoryID uuid.UUID) {
	t.Helper()
	err := db.Exec(`
		INSERT INTO blogs (id, author_id, category_id, title, slug, content, status, visibility, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, blogID, authorID, categoryID, "Test Blog", "test-blog-"+blogID.String(),
		"Test content", "draft", "public").Error
	require.NoError(t, err, "failed to create blog")
}

func createTag(t *testing.T, db *gorm.DB, tagID uuid.UUID) {
	t.Helper()
	err := db.Exec(`
		INSERT INTO tags (id, name, slug, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, tagID, "Test Tag "+tagID.String(), "test-tag-"+tagID.String()).Error
	require.NoError(t, err, "failed to create tag")
}
