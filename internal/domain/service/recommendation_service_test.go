package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetPersonalizedFeed_WithInterests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()

	// Expectations
	mockUserRepo.EXPECT().GetInterests(gomock.Any(), userID).Return([]entity.Tag{
		{ID: tagID, Name: "Go", Slug: "go"},
	}, nil)

	// Construct expected filter
	expectedFilter := repository.BlogFilter{
		TagIDs: []uuid.UUID{tagID},
	}
	// We might need to ensure other fields match what the service produces (zero values).
	// The service creates `filter := repository.BlogFilter{}` and sets TagIDs.

	mockBlogRepo.EXPECT().FindAll(gomock.Any(), expectedFilter, gomock.Any()).Return(&repository.PaginatedResult[entity.Blog]{
		Data:  []entity.Blog{{Title: "Go Blog"}},
		Total: 1,
	}, nil)

	// Act
	result, err := svc.GetPersonalizedFeed(context.Background(), userID, repository.Pagination{Page: 1, PageSize: 10})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Go Blog", result.Data[0].Title)
}

func TestGetPersonalizedFeed_NoInterests_Fallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()

	// Expectations
	mockUserRepo.EXPECT().GetInterests(gomock.Any(), userID).Return([]entity.Tag{}, nil)

	// Expected filter is empty (zero value)
	expectedFilter := repository.BlogFilter{}

	mockBlogRepo.EXPECT().FindAll(gomock.Any(), expectedFilter, gomock.Any()).Return(&repository.PaginatedResult[entity.Blog]{
		Data:  []entity.Blog{{Title: "Recent Blog"}},
		Total: 1,
	}, nil)

	// Act
	result, err := svc.GetPersonalizedFeed(context.Background(), userID, repository.Pagination{Page: 1, PageSize: 10})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Recent Blog", result.Data[0].Title)
}

func TestGetRelatedBlogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	blogID := uuid.New()

	mockBlogRepo.EXPECT().FindRelated(gomock.Any(), blogID, 3).Return([]entity.Blog{
		{Title: "Related 1"},
	}, nil)

	blogs, err := svc.GetRelatedBlogs(context.Background(), blogID, 3)

	assert.NoError(t, err)
	assert.Len(t, blogs, 1)
	assert.Equal(t, "Related 1", blogs[0].Title)
}

func TestUpdateInterests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()
	tagIDs := []uuid.UUID{tagID}

	mockTagRepo.EXPECT().FindByIDs(gomock.Any(), tagIDs).Return([]entity.Tag{{ID: tagID}}, nil)
	mockUserRepo.EXPECT().ReplaceInterests(gomock.Any(), userID, tagIDs).Return(nil)

	err := svc.UpdateInterests(context.Background(), userID, tagIDs)

	assert.NoError(t, err)
}

func TestUpdateInterests_TagNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()
	tagIDs := []uuid.UUID{tagID}

	// Mock returning empty list
	mockTagRepo.EXPECT().FindByIDs(gomock.Any(), tagIDs).Return([]entity.Tag{}, nil)

	err := svc.UpdateInterests(context.Background(), userID, tagIDs)

	assert.Error(t, err)
	assert.Equal(t, "one or more tags not found", err.Error())
}

func TestUpdateInterests_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()
	tagIDs := []uuid.UUID{tagID}

	mockTagRepo.EXPECT().FindByIDs(gomock.Any(), tagIDs).Return(nil, errors.New("db error"))

	err := svc.UpdateInterests(context.Background(), userID, tagIDs)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
