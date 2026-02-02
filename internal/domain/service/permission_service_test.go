package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPermissionService_CheckPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	svc := service.NewPermissionService(mockRoleRepo)
	ctx := context.Background()
	userID := uuid.New()
	resource := entity.ResourceBlogs

	tests := []struct {
		name          string
		permission    entity.Permission
		mockSetup     func()
		expectedAllow bool
		expectedError error
	}{
		{
			name:       "Has permission",
			permission: entity.PermissionRead,
			mockSetup: func() {
				mockRoleRepo.EXPECT().
					GetUserPermission(ctx, userID, resource).
					Return(entity.PermissionRead|entity.PermissionCreate, nil)
			},
			expectedAllow: true,
			expectedError: nil,
		},
		{
			name:       "Does not have permission",
			permission: entity.PermissionDelete,
			mockSetup: func() {
				mockRoleRepo.EXPECT().
					GetUserPermission(ctx, userID, resource).
					Return(entity.PermissionRead|entity.PermissionCreate, nil)
			},
			expectedAllow: false,
			expectedError: nil,
		},
		{
			name:       "Repository error",
			permission: entity.PermissionRead,
			mockSetup: func() {
				mockRoleRepo.EXPECT().
					GetUserPermission(ctx, userID, resource).
					Return(entity.Permission(0), errors.New("repo error"))
			},
			expectedAllow: false,
			expectedError: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			allowed, err := svc.CheckPermission(ctx, userID, resource, tt.permission)
			assert.Equal(t, tt.expectedAllow, allowed)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPermissionService_GetUserPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	svc := service.NewPermissionService(mockRoleRepo)
	ctx := context.Background()
	userID := uuid.New()
	resource := entity.ResourceBlogs

	tests := []struct {
		name           string
		mockSetup      func()
		expectedResult entity.Permission
		expectedError  error
	}{
		{
			name: "Success",
			mockSetup: func() {
				mockRoleRepo.EXPECT().
					GetUserPermission(ctx, userID, resource).
					Return(entity.PermissionRead|entity.PermissionCreate, nil)
			},
			expectedResult: entity.PermissionRead | entity.PermissionCreate,
			expectedError:  nil,
		},
		{
			name: "Repository error",
			mockSetup: func() {
				mockRoleRepo.EXPECT().
					GetUserPermission(ctx, userID, resource).
					Return(entity.Permission(0), errors.New("repo error"))
			},
			expectedResult: entity.Permission(0),
			expectedError:  errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result, err := svc.GetUserPermission(ctx, userID, resource)
			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPermissionService_HasPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	svc := service.NewPermissionService(mockRoleRepo)

	tests := []struct {
		name       string
		permission entity.Permission
		required   entity.Permission
		expected   bool
	}{
		{
			name:       "Exact match",
			permission: entity.PermissionRead,
			required:   entity.PermissionRead,
			expected:   true,
		},
		{
			name:       "Superset match",
			permission: entity.PermissionRead | entity.PermissionCreate,
			required:   entity.PermissionRead,
			expected:   true,
		},
		{
			name:       "Subset mismatch",
			permission: entity.PermissionRead,
			required:   entity.PermissionRead | entity.PermissionCreate,
			expected:   false,
		},
		{
			name:       "No match",
			permission: entity.PermissionRead,
			required:   entity.PermissionDelete,
			expected:   false,
		},
		{
			name:       "Empty required",
			permission: entity.PermissionRead,
			required:   0,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.HasPermission(tt.permission, tt.required)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPermissionService_CombinePermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	svc := service.NewPermissionService(mockRoleRepo)

	tests := []struct {
		name        string
		permissions []entity.Permission
		expected    entity.Permission
	}{
		{
			name:        "Single permission",
			permissions: []entity.Permission{entity.PermissionRead},
			expected:    entity.PermissionRead,
		},
		{
			name:        "Multiple permissions",
			permissions: []entity.Permission{entity.PermissionRead, entity.PermissionCreate},
			expected:    entity.PermissionRead | entity.PermissionCreate,
		},
		{
			name:        "Overlapping permissions",
			permissions: []entity.Permission{entity.PermissionRead | entity.PermissionCreate, entity.PermissionCreate | entity.PermissionUpdate},
			expected:    entity.PermissionRead | entity.PermissionCreate | entity.PermissionUpdate,
		},
		{
			name:        "Empty list",
			permissions: []entity.Permission{},
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.CombinePermissions(tt.permissions...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
