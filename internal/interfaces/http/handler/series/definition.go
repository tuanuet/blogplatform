package series

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

// SeriesHandler defines the interface for series handlers
type SeriesHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetBySlug(c *gin.Context)
	List(c *gin.Context)
	AddBlog(c *gin.Context)
	RemoveBlog(c *gin.Context)
	GetHighlightedSeries(c *gin.Context)
}
