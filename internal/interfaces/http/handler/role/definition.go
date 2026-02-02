package role

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type RoleHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	SetPermission(c *gin.Context)
	GetUserRoles(c *gin.Context)
	AssignRole(c *gin.Context)
	RemoveRole(c *gin.Context)
	GetMyPermission(c *gin.Context)
}
