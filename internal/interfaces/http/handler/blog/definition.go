package blog

import "github.com/gin-gonic/gin"

type BlogHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Publish(c *gin.Context)
	Unpublish(c *gin.Context)
	React(c *gin.Context)
}
