package entity

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a bitmask permission value
type Permission int

// Permission bitmask values
const (
	PermissionRead   Permission = 1 << iota // 1
	PermissionCreate                        // 2
	PermissionUpdate                        // 4
	PermissionDelete                        // 8
)

// PermissionAll represents full access (READ + CREATE + UPDATE + DELETE = 15)
const PermissionAll Permission = PermissionRead | PermissionCreate | PermissionUpdate | PermissionDelete

// Has checks if the permission includes the given permission
func (p Permission) Has(perm Permission) bool {
	return p&perm == perm
}

// CanRead checks if the permission includes READ
func (p Permission) CanRead() bool {
	return p.Has(PermissionRead)
}

// CanCreate checks if the permission includes CREATE
func (p Permission) CanCreate() bool {
	return p.Has(PermissionCreate)
}

// CanUpdate checks if the permission includes UPDATE
func (p Permission) CanUpdate() bool {
	return p.Has(PermissionUpdate)
}

// CanDelete checks if the permission includes DELETE
func (p Permission) CanDelete() bool {
	return p.Has(PermissionDelete)
}

// String returns the human-readable representation of permissions
func (p Permission) String() string {
	var perms []string
	if p.CanRead() {
		perms = append(perms, "READ")
	}
	if p.CanCreate() {
		perms = append(perms, "CREATE")
	}
	if p.CanUpdate() {
		perms = append(perms, "UPDATE")
	}
	if p.CanDelete() {
		perms = append(perms, "DELETE")
	}
	if len(perms) == 0 {
		return "NONE"
	}
	result := perms[0]
	for i := 1; i < len(perms); i++ {
		result += " | " + perms[i]
	}
	return result
}

// Role represents a user role
type Role struct {
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string           `gorm:"size:50;not null;unique" json:"name"`
	Description string           `gorm:"size:255" json:"description"`
	CreatedAt   time.Time        `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time        `gorm:"not null;default:now()" json:"updatedAt"`
	Permissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
	Users       []User           `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// RolePermission represents resource-based permissions for a role
type RolePermission struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"roleId"`
	Resource    string     `gorm:"size:50;not null" json:"resource"`
	Permissions Permission `gorm:"not null;default:0" json:"permissions"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	Role        *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName returns the table name for RolePermission
func (RolePermission) TableName() string {
	return "role_permissions"
}

// HasPermission checks if the role permission includes the given permission
func (rp *RolePermission) HasPermission(perm Permission) bool {
	return rp.Permissions.Has(perm)
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null;index" json:"roleId"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName returns the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}

// Resource constants for permission checking
const (
	ResourceBlogs      = "blogs"
	ResourceCategories = "categories"
	ResourceTags       = "tags"
	ResourceComments   = "comments"
	ResourceUsers      = "users"
	ResourceRoles      = "roles"
)
