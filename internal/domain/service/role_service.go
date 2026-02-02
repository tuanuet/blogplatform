package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// Domain errors
var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
)

// RoleService handles role-related domain logic
// This is a pure domain service with no infrastructure dependencies
type RoleService interface {
	// CreateRole creates a new role
	CreateRole(ctx context.Context, name, description string) (*entity.Role, error)

	// GetRole retrieves a role by ID
	GetRole(ctx context.Context, id uuid.UUID) (*entity.Role, error)

	// GetRoleByName retrieves a role by name
	GetRoleByName(ctx context.Context, name string) (*entity.Role, error)

	// ListRoles returns all roles
	ListRoles(ctx context.Context) ([]entity.Role, error)

	// UpdateRole updates a role
	UpdateRole(ctx context.Context, role *entity.Role) error

	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// SetPermission sets the permission for a role on a resource
	SetPermission(ctx context.Context, roleID uuid.UUID, resource string, permission entity.Permission) error

	// GetUserRoles returns all roles for a user
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)

	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error

	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
}

type roleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService creates a new role domain service
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (s *roleService) CreateRole(ctx context.Context, name, description string) (*entity.Role, error) {
	// Check if role exists
	existing, err := s.roleRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleAlreadyExists
	}

	role := &entity.Role{
		Name:        name,
		Description: description,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) GetRole(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *roleService) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	role, err := s.roleRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *roleService) ListRoles(ctx context.Context) ([]entity.Role, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *roleService) UpdateRole(ctx context.Context, role *entity.Role) error {
	existing, err := s.roleRepo.FindByID(ctx, role.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	role.UpdatedAt = time.Now()
	return s.roleRepo.Update(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	existing, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	return s.roleRepo.Delete(ctx, id)
}

func (s *roleService) SetPermission(ctx context.Context, roleID uuid.UUID, resource string, permission entity.Permission) error {
	existing, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	return s.roleRepo.SetPermission(ctx, roleID, resource, permission)
}

func (s *roleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

func (s *roleService) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	existing, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrRoleNotFound
	}

	return s.roleRepo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *roleService) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.RemoveRoleFromUser(ctx, userID, roleID)
}
