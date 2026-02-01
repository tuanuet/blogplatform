package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *entity.Role) error

	// FindByID finds a role by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*entity.Role, error)

	// FindAll returns all roles
	FindAll(ctx context.Context) ([]entity.Role, error)

	// Update updates a role
	Update(ctx context.Context, role *entity.Role) error

	// Delete deletes a role
	Delete(ctx context.Context, id uuid.UUID) error

	// SetPermission sets the permission for a role on a resource
	SetPermission(ctx context.Context, roleID uuid.UUID, resource string, permission entity.Permission) error

	// GetPermission gets the permission for a role on a resource
	GetPermission(ctx context.Context, roleID uuid.UUID, resource string) (entity.Permission, error)

	// GetUserRoles returns all roles for a user
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)

	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error

	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error

	// GetUserPermission returns the combined permission for a user on a resource (from all roles)
	GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (entity.Permission, error)
}
