package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", id).
		First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &role, err
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&role).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &role, err
}

func (r *roleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Order("name ASC").
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, "id = ?", id).Error
}

func (r *roleRepository) SetPermission(ctx context.Context, roleID uuid.UUID, resource string, permission entity.Permission) error {
	// Upsert permission
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO role_permissions (id, role_id, resource, permissions, created_at, updated_at)
		VALUES (gen_random_uuid(), ?, ?, ?, NOW(), NOW())
		ON CONFLICT (role_id, resource) 
		DO UPDATE SET permissions = ?, updated_at = NOW()
	`, roleID, resource, permission, permission).Error
}

func (r *roleRepository) GetPermission(ctx context.Context, roleID uuid.UUID, resource string) (entity.Permission, error) {
	var perm entity.RolePermission
	err := r.db.WithContext(ctx).
		Where("role_id = ? AND resource = ?", roleID, resource).
		First(&perm).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return perm.Permissions, nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	userRole := entity.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(&userRole).Error
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entity.UserRole{}).Error
}

func (r *roleRepository) GetUserPermission(ctx context.Context, userID uuid.UUID, resource string) (entity.Permission, error) {
	// Combine permissions from all user's roles using bitwise OR
	var result struct {
		Permissions int
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(BIT_OR(rp.permissions), 0) as permissions
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND rp.resource = ?
	`, userID, resource).Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return entity.Permission(result.Permissions), nil
}
