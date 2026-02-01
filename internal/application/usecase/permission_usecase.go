package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/domain/entity"
	domainService "github.com/aiagent/boilerplate/internal/domain/service"
	"github.com/aiagent/boilerplate/internal/infrastructure/cache"
	"github.com/aiagent/boilerplate/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Cache key prefixes and TTL
const (
	cacheKeyUserPermission = "rbac:user:%s:resource:%s"
	permissionCacheTTL     = 5 * time.Minute
)

// PermissionUseCase handles permission-related application logic
// Orchestrates domain service + caching
type PermissionUseCase interface {
	// CheckPermission checks if a user has a specific permission (with caching)
	CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permissionName string) (bool, error)

	// GetUserPermission returns user's permission on a resource (with caching)
	GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (*dto.UserPermissionResponse, error)

	// InvalidateUserPermissions invalidates all cached permissions for a user
	InvalidateUserPermissions(ctx context.Context, userID uuid.UUID) error

	// InvalidateResourcePermissions invalidates all cached permissions for a resource
	InvalidateResourcePermissions(ctx context.Context, resource string) error
}

type permissionUseCase struct {
	permissionSvc domainService.PermissionService
	cache         *cache.RedisClient
}

// NewPermissionUseCase creates a new permission use case
func NewPermissionUseCase(
	permissionSvc domainService.PermissionService,
	cache *cache.RedisClient,
) PermissionUseCase {
	return &permissionUseCase{
		permissionSvc: permissionSvc,
		cache:         cache,
	}
}

func (uc *permissionUseCase) CheckPermission(ctx context.Context, userID uuid.UUID, resource string, permissionName string) (bool, error) {
	permInt, err := uc.getUserPermissionCached(ctx, userID, resource)
	if err != nil {
		return false, err
	}
	perm := entity.Permission(permInt)

	// Map permission name to constant
	requiredInt := mapPermissionName(permissionName)
	required := entity.Permission(requiredInt)
	return uc.permissionSvc.HasPermission(perm, required), nil
}

func (uc *permissionUseCase) GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (*dto.UserPermissionResponse, error) {
	permInt, err := uc.getUserPermissionCached(ctx, userID, resource)
	if err != nil {
		return nil, err
	}
	perm := entity.Permission(permInt)

	return &dto.UserPermissionResponse{
		UserID:      userID,
		Resource:    resource,
		Permissions: int(perm),
		CanRead:     perm.CanRead(),
		CanCreate:   perm.CanCreate(),
		CanUpdate:   perm.CanUpdate(),
		CanDelete:   perm.CanDelete(),
	}, nil
}

func (uc *permissionUseCase) InvalidateUserPermissions(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("rbac:user:%s:resource:*", userID.String())
	return uc.cache.DeleteByPattern(ctx, pattern)
}

func (uc *permissionUseCase) InvalidateResourcePermissions(ctx context.Context, resource string) error {
	pattern := fmt.Sprintf("rbac:user:*:resource:%s", resource)
	return uc.cache.DeleteByPattern(ctx, pattern)
}

// getUserPermissionCached gets user permission from cache or domain service
func (uc *permissionUseCase) getUserPermissionCached(ctx context.Context, userID uuid.UUID, resource string) (int, error) {
	cacheKey := fmt.Sprintf(cacheKeyUserPermission, userID.String(), resource)

	// Try cache first
	var cached int
	if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	} else if err != redis.Nil {
		logger.Error("Failed to get permission from cache", err, nil)
	}

	// Cache miss - get from domain service
	perm, err := uc.permissionSvc.GetUserPermission(ctx, userID, resource)
	if err != nil {
		return 0, err
	}

	// Store in cache
	if err := uc.cache.Set(ctx, cacheKey, int(perm), permissionCacheTTL); err != nil {
		logger.Error("Failed to cache permission", err, nil)
	}

	return int(perm), nil
}

// Helper to map permission names to constants
func mapPermissionName(name string) int {
	switch name {
	case "read":
		return 1
	case "create":
		return 2
	case "update":
		return 4
	case "delete":
		return 8
	default:
		return 0
	}
}
