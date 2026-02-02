package category

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type CategoryHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
