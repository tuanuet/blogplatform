package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/cache"
	"github.com/aiagent/pkg/logger"
	"github.com/google/uuid"
)

// Cache keys and TTL
const (
	cacheKeyUserRoles = "rbac:user:%s:roles"
	cacheKeyRole      = "rbac:role:%s"
	cacheKeyRolesList = "rbac:roles:list"
	roleCacheTTL      = 10 * time.Minute
)

// Use case errors (re-exported from domain)
var (
	ErrRoleNotFound      = domainService.ErrRoleNotFound
	ErrRoleAlreadyExists = domainService.ErrRoleAlreadyExists
)

// RoleUseCase handles role-related application logic
// Orchestrates domain service + caching + DTO mapping
type RoleUseCase interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error)

	// GetRole retrieves a role by ID (with caching)
	GetRole(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error)

	// ListRoles returns all roles (with caching)
	ListRoles(ctx context.Context) ([]dto.RoleResponse, error)

	// UpdateRole updates a role
	UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)

	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// SetPermission sets a permission for a role on a resource
	SetPermission(ctx context.Context, roleID uuid.UUID, req dto.SetPermissionRequest) error

	// GetUserRoles returns all roles for a user (with caching)
	GetUserRoles(ctx context.Context, userID uuid.UUID) (*dto.UserRolesResponse, error)

	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error

	// RemoveRole removes a role from a user
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error

	// CheckPermission checks if a user has a specific permission
	CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permission entity.Permission) (bool, error)

	// GetUserPermission returns user's permission on a resource
	GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (*dto.UserPermissionResponse, error)
}

type roleUseCase struct {
	roleSvc       domainService.RoleService
	permissionSvc domainService.PermissionService
	permissionUC  PermissionUseCase
	cache         *cache.RedisClient
}

// NewRoleUseCase creates a new role use case
func NewRoleUseCase(
	roleSvc domainService.RoleService,
	permissionSvc domainService.PermissionService,
	cache *cache.RedisClient,
) RoleUseCase {
	permissionUC := NewPermissionUseCase(permissionSvc, cache)

	return &roleUseCase{
		roleSvc:       roleSvc,
		permissionSvc: permissionSvc,
		permissionUC:  permissionUC,
		cache:         cache,
	}
}

func (uc *roleUseCase) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role, err := uc.roleSvc.CreateRole(ctx, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	// Invalidate roles list cache
	uc.invalidateRolesListCache(ctx)

	return uc.toRoleResponse(role), nil
}

func (uc *roleUseCase) GetRole(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error) {
	cacheKey := fmt.Sprintf(cacheKeyRole, id.String())

	// Try cache first
	var cached dto.RoleResponse
	if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Cache miss - get from domain service
	role, err := uc.roleSvc.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := uc.toRoleResponse(role)

	// Store in cache
	if err := uc.cache.Set(ctx, cacheKey, resp, roleCacheTTL); err != nil {
		logger.Error("Failed to cache role", err, nil)
	}

	return resp, nil
}

func (uc *roleUseCase) ListRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	// Try cache first
	var cached []dto.RoleResponse
	if err := uc.cache.Get(ctx, cacheKeyRolesList, &cached); err == nil {
		return cached, nil
	}

	// Cache miss - get from domain service
	roles, err := uc.roleSvc.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = *uc.toRoleResponse(&role)
	}

	// Store in cache
	if err := uc.cache.Set(ctx, cacheKeyRolesList, responses, roleCacheTTL); err != nil {
		logger.Error("Failed to cache roles list", err, nil)
	}

	return responses, nil
}

func (uc *roleUseCase) UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := uc.roleSvc.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := uc.roleSvc.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	// Invalidate caches
	uc.invalidateRoleCache(ctx, id)
	uc.invalidateRolesListCache(ctx)

	return uc.toRoleResponse(role), nil
}

func (uc *roleUseCase) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := uc.roleSvc.DeleteRole(ctx, id); err != nil {
		return err
	}

	// Invalidate caches
	uc.invalidateRoleCache(ctx, id)
	uc.invalidateRolesListCache(ctx)
	uc.invalidateAllUserPermissionCaches(ctx)

	return nil
}

func (uc *roleUseCase) SetPermission(ctx context.Context, roleID uuid.UUID, req dto.SetPermissionRequest) error {
	if err := uc.roleSvc.SetPermission(ctx, roleID, req.Resource, entity.Permission(req.Permissions)); err != nil {
		return err
	}

	// Invalidate caches
	uc.invalidateRoleCache(ctx, roleID)
	if err := uc.permissionUC.InvalidateResourcePermissions(ctx, req.Resource); err != nil {
		logger.Error("Failed to invalidate resource permissions", err, nil)
	}

	return nil
}

func (uc *roleUseCase) GetUserRoles(ctx context.Context, userID uuid.UUID) (*dto.UserRolesResponse, error) {
	cacheKey := fmt.Sprintf(cacheKeyUserRoles, userID.String())

	// Try cache first
	var cached dto.UserRolesResponse
	if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// Cache miss - get from domain service
	roles, err := uc.roleSvc.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = *uc.toRoleResponse(&role)
	}

	resp := &dto.UserRolesResponse{
		UserID: userID,
		Roles:  roleResponses,
	}

	// Store in cache
	if err := uc.cache.Set(ctx, cacheKey, resp, roleCacheTTL); err != nil {
		logger.Error("Failed to cache user roles", err, nil)
	}

	return resp, nil
}

func (uc *roleUseCase) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	if err := uc.roleSvc.AssignRoleToUser(ctx, userID, roleID); err != nil {
		return err
	}

	// Invalidate user's caches
	uc.invalidateUserCaches(ctx, userID)

	return nil
}

func (uc *roleUseCase) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	if err := uc.roleSvc.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		return err
	}

	// Invalidate user's caches
	uc.invalidateUserCaches(ctx, userID)

	return nil
}

func (uc *roleUseCase) CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permission entity.Permission) (bool, error) {
	perm, err := uc.permissionSvc.GetUserPermission(ctx, userID, resource)
	if err != nil {
		return false, err
	}
	return uc.permissionSvc.HasPermission(perm, permission), nil
}

func (uc *roleUseCase) GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (*dto.UserPermissionResponse, error) {
	return uc.permissionUC.GetUserPermission(ctx, userID, resource)
}

// DTO mapping
func (uc *roleUseCase) toRoleResponse(role *entity.Role) *dto.RoleResponse {
	resp := &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
	}

	if len(role.Permissions) > 0 {
		resp.Permissions = make([]dto.PermissionResponse, len(role.Permissions))
		for i, perm := range role.Permissions {
			resp.Permissions[i] = dto.PermissionResponse{
				Resource:    perm.Resource,
				Permissions: int(perm.Permissions),
				CanRead:     perm.Permissions.CanRead(),
				CanCreate:   perm.Permissions.CanCreate(),
				CanUpdate:   perm.Permissions.CanUpdate(),
				CanDelete:   perm.Permissions.CanDelete(),
			}
		}
	}

	return resp
}

// Cache invalidation helpers
func (uc *roleUseCase) invalidateRoleCache(ctx context.Context, roleID uuid.UUID) {
	cacheKey := fmt.Sprintf(cacheKeyRole, roleID.String())
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		logger.Error("Failed to invalidate role cache", err, nil)
	}
}

func (uc *roleUseCase) invalidateRolesListCache(ctx context.Context) {
	if err := uc.cache.Delete(ctx, cacheKeyRolesList); err != nil {
		logger.Error("Failed to invalidate roles list cache", err, nil)
	}
}

func (uc *roleUseCase) invalidateUserCaches(ctx context.Context, userID uuid.UUID) {
	// Invalidate user roles cache
	userRolesKey := fmt.Sprintf(cacheKeyUserRoles, userID.String())
	if err := uc.cache.Delete(ctx, userRolesKey); err != nil {
		logger.Error("Failed to invalidate user roles cache", err, nil)
	}

	// Invalidate user permissions
	if err := uc.permissionUC.InvalidateUserPermissions(ctx, userID); err != nil {
		logger.Error("Failed to invalidate user permissions", err, nil)
	}
}

func (uc *roleUseCase) invalidateAllUserPermissionCaches(ctx context.Context) {
	pattern := "rbac:user:*:resource:*"
	if err := uc.cache.DeleteByPattern(ctx, pattern); err != nil {
		logger.Error("Failed to invalidate all user permission caches", err, nil)
	}
}

// Ensure interface compliance (re-export for backward compatibility)
var _ = errors.New // Silence unused import
