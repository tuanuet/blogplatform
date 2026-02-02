package entity_test

import (
	"testing"

	"github.com/aiagent/internal/domain/entity"
)

func TestPermission_Has(t *testing.T) {
	tests := []struct {
		name       string
		permission entity.Permission
		check      entity.Permission
		expected   bool
	}{
		{
			name:       "full permission has read",
			permission: entity.PermissionAll,
			check:      entity.PermissionRead,
			expected:   true,
		},
		{
			name:       "full permission has create",
			permission: entity.PermissionAll,
			check:      entity.PermissionCreate,
			expected:   true,
		},
		{
			name:       "full permission has update",
			permission: entity.PermissionAll,
			check:      entity.PermissionUpdate,
			expected:   true,
		},
		{
			name:       "full permission has delete",
			permission: entity.PermissionAll,
			check:      entity.PermissionDelete,
			expected:   true,
		},
		{
			name:       "read only does not have create",
			permission: entity.PermissionRead,
			check:      entity.PermissionCreate,
			expected:   false,
		},
		{
			name:       "read+create has both",
			permission: entity.PermissionRead | entity.PermissionCreate,
			check:      entity.PermissionRead,
			expected:   true,
		},
		{
			name:       "read+create has create",
			permission: entity.PermissionRead | entity.PermissionCreate,
			check:      entity.PermissionCreate,
			expected:   true,
		},
		{
			name:       "read+create does not have delete",
			permission: entity.PermissionRead | entity.PermissionCreate,
			check:      entity.PermissionDelete,
			expected:   false,
		},
		{
			name:       "no permission",
			permission: 0,
			check:      entity.PermissionRead,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.permission.Has(tt.check)
			if result != tt.expected {
				t.Errorf("Permission(%d).Has(%d) = %v, expected %v",
					tt.permission, tt.check, result, tt.expected)
			}
		})
	}
}

func TestPermission_CanMethods(t *testing.T) {
	tests := []struct {
		name      string
		perm      entity.Permission
		canRead   bool
		canCreate bool
		canUpdate bool
		canDelete bool
	}{
		{
			name:      "all permissions",
			perm:      entity.PermissionAll,
			canRead:   true,
			canCreate: true,
			canUpdate: true,
			canDelete: true,
		},
		{
			name:      "read only",
			perm:      entity.PermissionRead,
			canRead:   true,
			canCreate: false,
			canUpdate: false,
			canDelete: false,
		},
		{
			name:      "editor (read+create+update)",
			perm:      entity.PermissionRead | entity.PermissionCreate | entity.PermissionUpdate,
			canRead:   true,
			canCreate: true,
			canUpdate: true,
			canDelete: false,
		},
		{
			name:      "contributor (read+create)",
			perm:      entity.PermissionRead | entity.PermissionCreate,
			canRead:   true,
			canCreate: true,
			canUpdate: false,
			canDelete: false,
		},
		{
			name:      "no permission",
			perm:      0,
			canRead:   false,
			canCreate: false,
			canUpdate: false,
			canDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.perm.CanRead() != tt.canRead {
				t.Errorf("CanRead() = %v, expected %v", tt.perm.CanRead(), tt.canRead)
			}
			if tt.perm.CanCreate() != tt.canCreate {
				t.Errorf("CanCreate() = %v, expected %v", tt.perm.CanCreate(), tt.canCreate)
			}
			if tt.perm.CanUpdate() != tt.canUpdate {
				t.Errorf("CanUpdate() = %v, expected %v", tt.perm.CanUpdate(), tt.canUpdate)
			}
			if tt.perm.CanDelete() != tt.canDelete {
				t.Errorf("CanDelete() = %v, expected %v", tt.perm.CanDelete(), tt.canDelete)
			}
		})
	}
}

func TestPermission_String(t *testing.T) {
	tests := []struct {
		name     string
		perm     entity.Permission
		expected string
	}{
		{
			name:     "all permissions",
			perm:     entity.PermissionAll,
			expected: "READ | CREATE | UPDATE | DELETE",
		},
		{
			name:     "read only",
			perm:     entity.PermissionRead,
			expected: "READ",
		},
		{
			name:     "read and create",
			perm:     entity.PermissionRead | entity.PermissionCreate,
			expected: "READ | CREATE",
		},
		{
			name:     "no permission",
			perm:     0,
			expected: "NONE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.perm.String()
			if result != tt.expected {
				t.Errorf("String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestPermissionValues(t *testing.T) {
	// Verify bitmask values are correct
	if entity.PermissionRead != 1 {
		t.Errorf("PermissionRead = %d, expected 1", entity.PermissionRead)
	}
	if entity.PermissionCreate != 2 {
		t.Errorf("PermissionCreate = %d, expected 2", entity.PermissionCreate)
	}
	if entity.PermissionUpdate != 4 {
		t.Errorf("PermissionUpdate = %d, expected 4", entity.PermissionUpdate)
	}
	if entity.PermissionDelete != 8 {
		t.Errorf("PermissionDelete = %d, expected 8", entity.PermissionDelete)
	}
	if entity.PermissionAll != 15 {
		t.Errorf("PermissionAll = %d, expected 15", entity.PermissionAll)
	}
}
