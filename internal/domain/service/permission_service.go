package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// PermissionService handles permission-related domain logic
// This is a pure domain service with no infrastructure dependencies
type PermissionService interface {
	// CheckPermission checks if a user has a specific permission on a resource
	CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permission entity.Permission) (bool, error)

	// GetUserPermission returns the combined permission for a user on a resource
	GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (entity.Permission, error)

	// HasPermission checks if the given permission value includes the required permission
	HasPermission(permission, required entity.Permission) bool

	// CombinePermissions combines multiple permissions using bitwise OR
	CombinePermissions(permissions ...entity.Permission) entity.Permission
}

type permissionService struct {
	roleRepo repository.RoleRepository
}

// NewPermissionService creates a new permission domain service
func NewPermissionService(roleRepo repository.RoleRepository) PermissionService {
	return &permissionService{
		roleRepo: roleRepo,
	}
}

// CheckPermission checks if a user has a specific permission on a resource
func (s *permissionService) CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permission entity.Permission) (bool, error) {
	userPerm, err := s.roleRepo.GetUserPermission(ctx, userID, resource)
	if err != nil {
		return false, err
	}
	return s.HasPermission(userPerm, permission), nil
}

// GetUserPermission returns the combined permission for a user on a resource
func (s *permissionService) GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (entity.Permission, error) {
	return s.roleRepo.GetUserPermission(ctx, userID, resource)
}

// HasPermission checks if the given permission value includes the required permission
// This is a pure domain logic function
func (s *permissionService) HasPermission(permission, required entity.Permission) bool {
	return permission.Has(required)
}

// CombinePermissions combines multiple permissions using bitwise OR
// This is a pure domain logic function
func (s *permissionService) CombinePermissions(permissions ...entity.Permission) entity.Permission {
	var result entity.Permission
	for _, p := range permissions {
		result |= p
	}
	return result
}
