package blog_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/blog"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlogService is a mock implementation of domainService.BlogService
type MockBlogService struct {
	mock.Mock
}

func (m *MockBlogService) Create(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error {
	args := m.Called(ctx, blog, tagIDs)
	return args.Error(0)
}

func (m *MockBlogService) GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*entity.Blog, error) {
	args := m.Called(ctx, id, viewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}

func (m *MockBlogService) GetBySlug(ctx context.Context, authorID uuid.UUID, slug string, viewerID *uuid.UUID) (*entity.Blog, error) {
	args := m.Called(ctx, authorID, slug, viewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}

func (m *MockBlogService) List(ctx context.Context, filter repository.BlogFilter, pagination repository.Pagination, viewerID *uuid.UUID) (*repository.PaginatedResult[entity.Blog], error) {
	args := m.Called(ctx, filter, pagination, viewerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[entity.Blog]), args.Error(1)
}

func (m *MockBlogService) Update(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error {
	args := m.Called(ctx, blog, tagIDs)
	return args.Error(0)
}

func (m *MockBlogService) Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error {
	args := m.Called(ctx, id, authorID)
	return args.Error(0)
}

func (m *MockBlogService) Publish(ctx context.Context, id uuid.UUID, authorID uuid.UUID, visibility entity.BlogVisibility, publishedAt *time.Time) (*entity.Blog, error) {
	args := m.Called(ctx, id, authorID, visibility, publishedAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}

func (m *MockBlogService) Unpublish(ctx context.Context, id uuid.UUID, authorID uuid.UUID) (*entity.Blog, error) {
	args := m.Called(ctx, id, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}

func (m *MockBlogService) React(ctx context.Context, id uuid.UUID, userID uuid.UUID, reactionType entity.ReactionType) (int, int, error) {
	args := m.Called(ctx, id, userID, reactionType)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockBlogService) CheckAccess(ctx context.Context, blog *entity.Blog, viewerID *uuid.UUID) error {
	args := m.Called(ctx, blog, viewerID)
	return args.Error(0)
}

func TestCreateBlog_WithPublishedAt(t *testing.T) {
	mockService := new(MockBlogService)
	uc := blog.NewBlogUseCase(mockService)

	authorID := uuid.New()
	futureTime := time.Now().Add(24 * time.Hour)

	req := &dto.CreateBlogRequest{
		Title:       "Future Post",
		Slug:        "future-post",
		Content:     "Content",
		PublishedAt: &futureTime, // This field doesn't exist yet in DTO, so this test will fail to compile first!
	}

	// Expectation: Create calls service with blog containing PublishedAt
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(b *entity.Blog) bool {
		return b.PublishedAt != nil && b.PublishedAt.Equal(futureTime)
	}), mock.Anything).Return(nil)

	// Mock GetByID for the return
	mockService.On("GetByID", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Blog{
		ID:          uuid.New(),
		AuthorID:    authorID,
		Title:       req.Title,
		PublishedAt: &futureTime,
	}, nil)

	_, err := uc.Create(context.Background(), authorID, req)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestListBlog_PublicFiltering(t *testing.T) {
	mockService := new(MockBlogService)
	uc := blog.NewBlogUseCase(mockService)

	// Scenario: Public listing (no viewer, or viewer is not author)
	// Should set PublishedBefore to Now

	params := &dto.BlogFilterParams{
		Page:     1,
		PageSize: 10,
	}

	mockService.On("List", mock.Anything, mock.MatchedBy(func(f repository.BlogFilter) bool {
		// We expect PublishedBefore to be set (approx Now)
		if f.PublishedBefore == nil {
			return false
		}
		// Check if it's close to Now (within 1 sec)
		return time.Since(*f.PublishedBefore) < time.Second
	}), mock.Anything, mock.Anything).Return(&repository.PaginatedResult[entity.Blog]{}, nil)

	_, err := uc.List(context.Background(), params, nil)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
