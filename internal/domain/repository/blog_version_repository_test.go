package repository

import (
	"context"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// mockBlogVersionRepository is a minimal implementation to verify the interface compiles
type mockBlogVersionRepository struct{}

func (m *mockBlogVersionRepository) Create(ctx context.Context, version *entity.BlogVersion) error {
	return nil
}

func (m *mockBlogVersionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.BlogVersion, error) {
	return nil, nil
}

func (m *mockBlogVersionRepository) FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.BlogVersion], error) {
	return nil, nil
}

func (m *mockBlogVersionRepository) GetNextVersionNumber(ctx context.Context, blogID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *mockBlogVersionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockBlogVersionRepository) CountByBlogID(ctx context.Context, blogID uuid.UUID) (int64, error) {
	return 0, nil
}

func (m *mockBlogVersionRepository) DeleteOldest(ctx context.Context, blogID uuid.UUID, keep int) error {
	return nil
}

// compileTimeCheck ensures mock implements BlogVersionRepository interface
var _ BlogVersionRepository = (*mockBlogVersionRepository)(nil)

func TestBlogVersionRepository_Interface(t *testing.T) {
	// This test verifies the BlogVersionRepository interface compiles correctly
	// and can be implemented by a concrete type

	ctx := context.Background()
	repo := &mockBlogVersionRepository{}
	blogID := uuid.New()

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Create method exists with correct signature",
			test: func(t *testing.T) {
				version := &entity.BlogVersion{}
				err := repo.Create(ctx, version)
				_ = err // We only care that it compiles and can be called
			},
		},
		{
			name: "FindByID method exists with correct signature",
			test: func(t *testing.T) {
				_, err := repo.FindByID(ctx, uuid.New())
				_ = err
			},
		},
		{
			name: "FindByBlogID method exists with correct signature",
			test: func(t *testing.T) {
				_, err := repo.FindByBlogID(ctx, blogID, Pagination{Page: 1, PageSize: 10})
				_ = err
			},
		},
		{
			name: "GetNextVersionNumber method exists with correct signature",
			test: func(t *testing.T) {
				_, err := repo.GetNextVersionNumber(ctx, blogID)
				_ = err
			},
		},
		{
			name: "Delete method exists with correct signature",
			test: func(t *testing.T) {
				err := repo.Delete(ctx, uuid.New())
				_ = err
			},
		},
		{
			name: "CountByBlogID method exists with correct signature",
			test: func(t *testing.T) {
				_, err := repo.CountByBlogID(ctx, blogID)
				_ = err
			},
		},
		{
			name: "DeleteOldest method exists with correct signature",
			test: func(t *testing.T) {
				err := repo.DeleteOldest(ctx, blogID, 5)
				_ = err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestBlogVersionRepository_PaginatedResultType(t *testing.T) {
	// Verify that FindByBlogID returns the correct paginated result type
	repo := &mockBlogVersionRepository{}

	result, err := repo.FindByBlogID(context.Background(), uuid.New(), Pagination{Page: 1, PageSize: 10})
	_ = result
	_ = err

	// The fact that this compiles means PaginatedResult[entity.BlogVersion] is valid
	var _ *PaginatedResult[entity.BlogVersion] = result
}
