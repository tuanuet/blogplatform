package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookmarkRepository is a mock implementation of repository.BookmarkRepository
type MockBookmarkRepository struct {
	mock.Mock
}

func (m *MockBookmarkRepository) Add(ctx context.Context, userID, blogID uuid.UUID) error {
	args := m.Called(ctx, userID, blogID)
	return args.Error(0)
}

func (m *MockBookmarkRepository) Remove(ctx context.Context, userID, blogID uuid.UUID) error {
	args := m.Called(ctx, userID, blogID)
	return args.Error(0)
}

func (m *MockBookmarkRepository) IsBookmarked(ctx context.Context, userID, blogID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, blogID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookmarkRepository) CountByBlog(ctx context.Context, blogID uuid.UUID) (int64, error) {
	args := m.Called(ctx, blogID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookmarkRepository) FindByUser(ctx context.Context, userID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error) {
	args := m.Called(ctx, userID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[entity.Blog]), args.Error(1)
}

func TestBookmarkBlog(t *testing.T) {
	// Arrange
	mockRepo := new(MockBookmarkRepository)
	uc := usecase.NewBookmarkUseCase(mockRepo)
	userID := uuid.New()
	blogID := uuid.New()

	mockRepo.On("Add", context.Background(), userID, blogID).Return(nil)

	// Act
	err := uc.BookmarkBlog(context.Background(), userID, blogID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUnbookmarkBlog(t *testing.T) {
	// Arrange
	mockRepo := new(MockBookmarkRepository)
	uc := usecase.NewBookmarkUseCase(mockRepo)
	userID := uuid.New()
	blogID := uuid.New()

	mockRepo.On("Remove", context.Background(), userID, blogID).Return(nil)

	// Act
	err := uc.UnbookmarkBlog(context.Background(), userID, blogID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserBookmarks(t *testing.T) {
	// Arrange
	mockRepo := new(MockBookmarkRepository)
	uc := usecase.NewBookmarkUseCase(mockRepo)
	userID := uuid.New()

	params := &dto.BlogFilterParams{
		Page:     1,
		PageSize: 10,
	}

	pagination := repository.Pagination{
		Page:     1,
		PageSize: 10,
	}

	blogID := uuid.New()
	authorID := uuid.New()
	now := time.Now()

	blogs := []entity.Blog{
		{
			ID:          blogID,
			Title:       "Test Blog",
			Slug:        "test-blog",
			AuthorID:    authorID,
			Status:      entity.BlogStatusPublished,
			Visibility:  entity.BlogVisibilityPublic,
			PublishedAt: &now,
			CreatedAt:   now,
		},
	}

	paginatedResult := &repository.PaginatedResult[entity.Blog]{
		Data:       blogs,
		Total:      1,
		TotalPages: 1,
		Page:       1,
		PageSize:   10,
	}

	mockRepo.On("FindByUser", context.Background(), userID, pagination).Return(paginatedResult, nil)

	// Act
	result, err := uc.GetUserBookmarks(context.Background(), userID, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, blogID, result.Data[0].ID)
	assert.Equal(t, "Test Blog", result.Data[0].Title)
	mockRepo.AssertExpectations(t)
}
