package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBlogVersion_TableName(t *testing.T) {
	bv := BlogVersion{}
	expected := "blog_versions"
	if bv.TableName() != expected {
		t.Errorf("TableName() = %v, expected %v", bv.TableName(), expected)
	}
}

func TestBlogVersion_Instantiation(t *testing.T) {
	id := uuid.New()
	blogID := uuid.New()
	categoryID := uuid.New()
	editorID := uuid.New()
	now := time.Now()
	excerpt := "Test excerpt"
	thumbnailURL := "https://example.com/image.jpg"
	changeSummary := "Initial version"

	tests := []struct {
		name string
		bv   BlogVersion
	}{
		{
			name: "can be instantiated with all required fields",
			bv: BlogVersion{
				ID:            id,
				BlogID:        blogID,
				VersionNumber: 1,
				Title:         "Test Blog Title",
				Slug:          "test-blog-title",
				Content:       "This is the blog content",
				Status:        BlogStatusDraft,
				Visibility:    BlogVisibilityPublic,
				EditorID:      editorID,
				CreatedAt:     now,
			},
		},
		{
			name: "can be instantiated with optional fields",
			bv: BlogVersion{
				ID:            uuid.New(),
				BlogID:        uuid.New(),
				VersionNumber: 2,
				Title:         "Another Blog",
				Slug:          "another-blog",
				Excerpt:       &excerpt,
				Content:       "Content with excerpt",
				ThumbnailURL:  &thumbnailURL,
				Status:        BlogStatusPublished,
				Visibility:    BlogVisibilitySubscribersOnly,
				CategoryID:    &categoryID,
				EditorID:      uuid.New(),
				ChangeSummary: &changeSummary,
				CreatedAt:     now,
			},
		},
		{
			name: "can be instantiated with minimal fields",
			bv: BlogVersion{
				BlogID:        blogID,
				VersionNumber: 1,
				Title:         "Minimal Blog",
				Slug:          "minimal-blog",
				Content:       "Minimal content",
				Status:        BlogStatusDraft,
				Visibility:    BlogVisibilityPublic,
				EditorID:      editorID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the struct fields can be accessed
			if tt.bv.ID != id && tt.name == "can be instantiated with all required fields" {
				t.Errorf("ID mismatch")
			}
			if tt.bv.BlogID != blogID && tt.name == "can be instantiated with all required fields" {
				t.Errorf("BlogID mismatch")
			}
			if tt.bv.VersionNumber != 1 && tt.name == "can be instantiated with all required fields" {
				t.Errorf("VersionNumber mismatch")
			}
			if tt.bv.Title != "Test Blog Title" && tt.name == "can be instantiated with all required fields" {
				t.Errorf("Title mismatch")
			}
			if tt.bv.Slug != "test-blog-title" && tt.name == "can be instantiated with all required fields" {
				t.Errorf("Slug mismatch")
			}
			if tt.bv.Content != "This is the blog content" && tt.name == "can be instantiated with all required fields" {
				t.Errorf("Content mismatch")
			}
			if tt.bv.Status != BlogStatusDraft && tt.name == "can be instantiated with all required fields" {
				t.Errorf("Status mismatch")
			}
			if tt.bv.Visibility != BlogVisibilityPublic && tt.name == "can be instantiated with all required fields" {
				t.Errorf("Visibility mismatch")
			}
			if tt.bv.EditorID != editorID && tt.name == "can be instantiated with all required fields" {
				t.Errorf("EditorID mismatch")
			}

			// Check optional fields
			if tt.name == "can be instantiated with optional fields" {
				if tt.bv.Excerpt == nil || *tt.bv.Excerpt != excerpt {
					t.Errorf("Excerpt mismatch")
				}
				if tt.bv.ThumbnailURL == nil || *tt.bv.ThumbnailURL != thumbnailURL {
					t.Errorf("ThumbnailURL mismatch")
				}
				if tt.bv.CategoryID == nil || *tt.bv.CategoryID != categoryID {
					t.Errorf("CategoryID mismatch")
				}
				if tt.bv.ChangeSummary == nil || *tt.bv.ChangeSummary != changeSummary {
					t.Errorf("ChangeSummary mismatch")
				}
			}
		})
	}
}

func TestBlogVersion_Fields(t *testing.T) {
	// Test that all expected fields exist and are accessible
	bv := BlogVersion{}

	// Required fields
	bv.ID = uuid.New()
	bv.BlogID = uuid.New()
	bv.VersionNumber = 1
	bv.Title = "Test"
	bv.Slug = "test"
	bv.Content = "Content"
	bv.Status = BlogStatusPublished
	bv.Visibility = BlogVisibilityPublic
	bv.EditorID = uuid.New()
	bv.CreatedAt = time.Now()

	// Optional fields (nullable)
	excerpt := "Excerpt"
	bv.Excerpt = &excerpt

	thumbnail := "https://example.com/image.jpg"
	bv.ThumbnailURL = &thumbnail

	categoryID := uuid.New()
	bv.CategoryID = &categoryID

	changeSummary := "Created"
	bv.ChangeSummary = &changeSummary

	// Relationships
	bv.Category = &Category{ID: uuid.New(), Name: "Tech"}
	bv.Editor = &User{ID: uuid.New(), Name: "John"}
	bv.Tags = []Tag{{ID: uuid.New(), Name: "Go"}}

	// If we reach here, all fields are accessible
}
