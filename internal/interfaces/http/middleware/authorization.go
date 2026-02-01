package middleware

import (
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Authorization middleware for role-based access control
type Authorization struct {
	roleUseCase usecase.RoleUseCase
}

// NewAuthorization creates a new authorization middleware
func NewAuthorization(roleUseCase usecase.RoleUseCase) *Authorization {
	return &Authorization{
		roleUseCase: roleUseCase,
	}
}

// RequirePermission returns a middleware that checks if the user has the required permission
func (a *Authorization) RequirePermission(resource string, permission entity.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			response.Unauthorized(c, "Invalid user ID")
			c.Abort()
			return
		}

		hasPermission, err := a.roleUseCase.CheckPermission(c.Request.Context(), uid, resource, permission)
		if err != nil {
			response.InternalServerError(c, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRead returns a middleware that checks if the user has READ permission
func (a *Authorization) RequireRead(resource string) gin.HandlerFunc {
	return a.RequirePermission(resource, entity.PermissionRead)
}

// RequireCreate returns a middleware that checks if the user has CREATE permission
func (a *Authorization) RequireCreate(resource string) gin.HandlerFunc {
	return a.RequirePermission(resource, entity.PermissionCreate)
}

// RequireUpdate returns a middleware that checks if the user has UPDATE permission
func (a *Authorization) RequireUpdate(resource string) gin.HandlerFunc {
	return a.RequirePermission(resource, entity.PermissionUpdate)
}

// RequireDelete returns a middleware that checks if the user has DELETE permission
func (a *Authorization) RequireDelete(resource string) gin.HandlerFunc {
	return a.RequirePermission(resource, entity.PermissionDelete)
}

// RequireAdmin returns a middleware that checks if the user has full permissions on the resource
func (a *Authorization) RequireAdmin(resource string) gin.HandlerFunc {
	return a.RequirePermission(resource, entity.PermissionAll)
}
