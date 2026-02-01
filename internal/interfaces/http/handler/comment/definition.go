package comment

import "github.com/gin-gonic/gin"

type CommentHandler interface {
	Create(c *gin.Context)
	GetByBlogID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
