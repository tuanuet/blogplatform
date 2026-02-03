package bookmark_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/bookmark"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBookmarkBlog(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookmarkRepository(ctrl)
	uc := bookmark.NewBookmarkUseCase(mockRepo)
	userID := uuid.New()
	blogID := uuid.New()

	mockRepo.EXPECT().Add(gomock.Any(), userID, blogID).Return(nil)

	// Act
	err := uc.BookmarkBlog(context.Background(), userID, blogID)

	// Assert
	assert.NoError(t, err)
}

func TestUnbookmarkBlog(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookmarkRepository(ctrl)
	uc := bookmark.NewBookmarkUseCase(mockRepo)
	userID := uuid.New()
	blogID := uuid.New()

	mockRepo.EXPECT().Remove(gomock.Any(), userID, blogID).Return(nil)

	// Act
	err := uc.UnbookmarkBlog(context.Background(), userID, blogID)

	// Assert
	assert.NoError(t, err)
}

func TestGetUserBookmarks(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookmarkRepository(ctrl)
	uc := bookmark.NewBookmarkUseCase(mockRepo)
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

	mockRepo.EXPECT().FindByUser(gomock.Any(), userID, pagination).Return(paginatedResult, nil)

	// Act
	result, err := uc.GetUserBookmarks(context.Background(), userID, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, blogID, result.Data[0].ID)
	assert.Equal(t, "Test Blog", result.Data[0].Title)
}
