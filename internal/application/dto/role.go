package dto

import "github.com/google/uuid"

// CreateRoleRequest represents the request to create a role
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=255"`
}

// UpdateRoleRequest represents the request to update a role
type UpdateRoleRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=50"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

// SetPermissionRequest represents the request to set a permission
type SetPermissionRequest struct {
	Resource    string `json:"resource" validate:"required,max=50"`
	Permissions int    `json:"permissions" validate:"min=0,max=15"`
}

// AssignRoleRequest represents the request to assign a role to a user
type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"roleId" validate:"required"`
}

// RoleResponse represents a role in responses
type RoleResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string               `json:"createdAt"`
}

// PermissionResponse represents a permission in responses
type PermissionResponse struct {
	Resource    string `json:"resource"`
	Permissions int    `json:"permissions"`
	CanRead     bool   `json:"canRead"`
	CanCreate   bool   `json:"canCreate"`
	CanUpdate   bool   `json:"canUpdate"`
	CanDelete   bool   `json:"canDelete"`
}

// UserRolesResponse represents user roles in responses
type UserRolesResponse struct {
	UserID uuid.UUID      `json:"userId"`
	Roles  []RoleResponse `json:"roles"`
}

// UserPermissionResponse represents a user's permission on a resource
type UserPermissionResponse struct {
	UserID      uuid.UUID `json:"userId"`
	Resource    string    `json:"resource"`
	Permissions int       `json:"permissions"`
	CanRead     bool      `json:"canRead"`
	CanCreate   bool      `json:"canCreate"`
	CanUpdate   bool      `json:"canUpdate"`
	CanDelete   bool      `json:"canDelete"`
}
