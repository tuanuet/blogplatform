package bookmark

import "github.com/gin-gonic/gin"

// BookmarkHandler defines the interface for bookmark HTTP handlers
type BookmarkHandler interface {
	// Bookmark handles POST /api/v1/blogs/:id/bookmarks
	Bookmark(c *gin.Context)

	// Unbookmark handles DELETE /api/v1/blogs/:id/bookmarks
	Unbookmark(c *gin.Context)

	// List handles GET /api/v1/me/bookmarks
	List(c *gin.Context)
}
