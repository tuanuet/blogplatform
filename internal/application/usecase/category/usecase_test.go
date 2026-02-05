package category_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/category"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestCategoryUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategorySvc := mocks.NewMockCategoryService(ctrl)
	uc := category.NewCategoryUseCase(mockCategorySvc)

	tests := []struct {
		name      string
		req       *dto.CreateCategoryRequest
		setupMock func()
		wantErr   bool
		errType   error
	}{
		{
			name: "success - create category",
			req: &dto.CreateCategoryRequest{
				Name:        "Technology",
				Slug:        "technology",
				Description: strPtr("Tech related content"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, cat *entity.Category) error {
						cat.ID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
						cat.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
						cat.UpdatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "error - slug already exists",
			req: &dto.CreateCategoryRequest{
				Name: "Technology",
				Slug: "technology",
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(domainService.ErrCategorySlugExists).
					Times(1)
			},
			wantErr: true,
			errType: domainService.ErrCategorySlugExists,
		},
		{
			name: "error - repository failure",
			req: &dto.CreateCategoryRequest{
				Name: "Technology",
				Slug: "technology",
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := uc.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Create() expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Create() error = %v, want error type %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Error("Create() result should not be nil")
				return
			}

			if result.Name != tt.req.Name {
				t.Errorf("Create() Name = %v, want %v", result.Name, tt.req.Name)
			}
			if result.Slug != tt.req.Slug {
				t.Errorf("Create() Slug = %v, want %v", result.Slug, tt.req.Slug)
			}
		})
	}
}

func TestCategoryUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategorySvc := mocks.NewMockCategoryService(ctrl)
	uc := category.NewCategoryUseCase(mockCategorySvc)

	categoryID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func()
		want      *dto.CategoryResponse
		wantErr   bool
		errType   error
	}{
		{
			name: "success - category found",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:          categoryID,
						Name:        "Technology",
						Slug:        "technology",
						Description: strPtr("Tech content"),
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).
					Times(1)
			},
			want: &dto.CategoryResponse{
				ID:          categoryID,
				Name:        "Technology",
				Slug:        "technology",
				Description: strPtr("Tech content"),
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: false,
		},
		{
			name: "error - category not found",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(nil, domainService.ErrCategoryNotFound).
					Times(1)
			},
			wantErr: true,
			errType: domainService.ErrCategoryNotFound,
		},
		{
			name: "error - repository failure",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := uc.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetByID() expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("GetByID() error = %v, want error type %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("GetByID() unexpected error = %v", err)
				return
			}

			if result.ID != tt.want.ID {
				t.Errorf("GetByID() ID = %v, want %v", result.ID, tt.want.ID)
			}
			if result.Name != tt.want.Name {
				t.Errorf("GetByID() Name = %v, want %v", result.Name, tt.want.Name)
			}
		})
	}
}

func TestCategoryUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategorySvc := mocks.NewMockCategoryService(ctrl)
	uc := category.NewCategoryUseCase(mockCategorySvc)

	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		page      int
		pageSize  int
		setupMock func()
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:     "success - list with pagination",
			page:     1,
			pageSize: 10,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					List(gomock.Any(), 1, 10).
					Return(&repository.PaginatedResult[entity.Category]{
						Data: []entity.Category{
							{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), Name: "Tech", Slug: "tech", CreatedAt: now, UpdatedAt: now},
							{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), Name: "Lifestyle", Slug: "lifestyle", CreatedAt: now, UpdatedAt: now},
						},
						Total:      2,
						Page:       1,
						PageSize:   10,
						TotalPages: 1,
					}, nil).
					Times(1)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:     "success - empty list",
			page:     1,
			pageSize: 10,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					List(gomock.Any(), 1, 10).
					Return(&repository.PaginatedResult[entity.Category]{
						Data:       []entity.Category{},
						Total:      0,
						Page:       1,
						PageSize:   10,
						TotalPages: 0,
					}, nil).
					Times(1)
			},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name:     "error - repository failure",
			page:     1,
			pageSize: 10,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					List(gomock.Any(), 1, 10).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := uc.List(context.Background(), tt.page, tt.pageSize)

			if tt.wantErr {
				if err == nil {
					t.Errorf("List() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("List() unexpected error = %v", err)
				return
			}

			if len(result.Data) != tt.wantCount {
				t.Errorf("List() returned %d items, want %d", len(result.Data), tt.wantCount)
			}
			if result.Total != tt.wantTotal {
				t.Errorf("List() Total = %d, want %d", result.Total, tt.wantTotal)
			}
			if result.Page != tt.page {
				t.Errorf("List() Page = %d, want %d", result.Page, tt.page)
			}
		})
	}
}

func TestCategoryUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategorySvc := mocks.NewMockCategoryService(ctrl)
	uc := category.NewCategoryUseCase(mockCategorySvc)

	categoryID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        uuid.UUID
		req       *dto.UpdateCategoryRequest
		setupMock func()
		wantErr   bool
		errType   error
	}{
		{
			name: "success - update name only",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Name: strPtr("Updated Technology"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:        categoryID,
						Name:      "Technology",
						Slug:      "technology",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, cat *entity.Category) error {
						cat.UpdatedAt = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "success - update slug",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Slug: strPtr("new-slug"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:        categoryID,
						Name:      "Technology",
						Slug:      "technology",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "success - update description",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Description: strPtr("New description"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:          categoryID,
						Name:        "Technology",
						Slug:        "technology",
						Description: strPtr("Old description"),
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "success - update all fields",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Name:        strPtr("Updated Name"),
				Slug:        strPtr("updated-slug"),
				Description: strPtr("Updated description"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:        categoryID,
						Name:      "Technology",
						Slug:      "technology",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "success - no fields to update (nil request)",
			id:   categoryID,
			req:  &dto.UpdateCategoryRequest{},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:        categoryID,
						Name:      "Technology",
						Slug:      "technology",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "error - category not found",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Name: strPtr("Updated Name"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(nil, domainService.ErrCategoryNotFound).
					Times(1)
			},
			wantErr: true,
			errType: domainService.ErrCategoryNotFound,
		},
		{
			name: "error - slug already exists",
			id:   categoryID,
			req: &dto.UpdateCategoryRequest{
				Slug: strPtr("existing-slug"),
			},
			setupMock: func() {
				mockCategorySvc.EXPECT().
					GetByID(gomock.Any(), categoryID).
					Return(&entity.Category{
						ID:        categoryID,
						Name:      "Technology",
						Slug:      "technology",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)

				mockCategorySvc.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(domainService.ErrCategorySlugExists).
					Times(1)
			},
			wantErr: true,
			errType: domainService.ErrCategorySlugExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := uc.Update(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Update() expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Update() error = %v, want error type %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("Update() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Error("Update() result should not be nil")
				return
			}

			if tt.req.Name != nil && result.Name != *tt.req.Name {
				t.Errorf("Update() Name = %v, want %v", result.Name, *tt.req.Name)
			}
			if tt.req.Slug != nil && result.Slug != *tt.req.Slug {
				t.Errorf("Update() Slug = %v, want %v", result.Slug, *tt.req.Slug)
			}
		})
	}
}

func TestCategoryUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategorySvc := mocks.NewMockCategoryService(ctrl)
	uc := category.NewCategoryUseCase(mockCategorySvc)

	categoryID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func()
		wantErr   bool
		errType   error
	}{
		{
			name: "success - delete category",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Delete(gomock.Any(), categoryID).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "error - category not found",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Delete(gomock.Any(), categoryID).
					Return(domainService.ErrCategoryNotFound).
					Times(1)
			},
			wantErr: true,
			errType: domainService.ErrCategoryNotFound,
		},
		{
			name: "error - repository failure",
			id:   categoryID,
			setupMock: func() {
				mockCategorySvc.EXPECT().
					Delete(gomock.Any(), categoryID).
					Return(errors.New("database error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := uc.Delete(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Delete() expected error but got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Delete() error = %v, want error type %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
			}
		})
	}
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}
