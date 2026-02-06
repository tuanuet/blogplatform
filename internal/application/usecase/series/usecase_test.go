package series_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/series"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateSeries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSeriesRepository(ctrl)
	uc := series.NewSeriesUseCase(mockRepo)
	userID := uuid.New()

	req := &dto.CreateSeriesRequest{
		Title:       "My Series",
		Slug:        "my-series",
		Description: "A description",
	}

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := uc.CreateSeries(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Slug, resp.Slug)
	assert.Equal(t, userID, resp.AuthorID)
}

func TestUpdateSeries_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSeriesRepository(ctrl)
	uc := series.NewSeriesUseCase(mockRepo)
	userID := uuid.New()
	seriesID := uuid.New()

	existingSeries := &entity.Series{
		ID:          seriesID,
		AuthorID:    userID,
		Title:       "Old Title",
		Slug:        "old-title",
		Description: "Old Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	req := &dto.UpdateSeriesRequest{
		Title:       "New Title",
		Description: "New Description",
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), seriesID).Return(existingSeries, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := uc.UpdateSeries(context.Background(), userID, seriesID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Title", resp.Title)
	assert.Equal(t, "New Description", resp.Description)
}

func TestUpdateSeries_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSeriesRepository(ctrl)
	uc := series.NewSeriesUseCase(mockRepo)
	userID := uuid.New()
	otherUserID := uuid.New()
	seriesID := uuid.New()

	existingSeries := &entity.Series{
		ID:       seriesID,
		AuthorID: otherUserID, // Different author
		Title:    "Old Title",
	}

	req := &dto.UpdateSeriesRequest{
		Title: "New Title",
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), seriesID).Return(existingSeries, nil)

	resp, err := uc.UpdateSeries(context.Background(), userID, seriesID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "unauthorized: you are not the author of this series", err.Error())
}

func TestDeleteSeries_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSeriesRepository(ctrl)
	uc := series.NewSeriesUseCase(mockRepo)
	userID := uuid.New()
	seriesID := uuid.New()

	existingSeries := &entity.Series{
		ID:       seriesID,
		AuthorID: userID,
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), seriesID).Return(existingSeries, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), seriesID).Return(nil)

	err := uc.DeleteSeries(context.Background(), userID, seriesID)

	assert.NoError(t, err)
}

func TestGetHighlightedSeries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSeriesRepository(ctrl)
	uc := series.NewSeriesUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		seriesID := uuid.New()
		authorID := uuid.New()
		now := time.Now()

		avatarURL := "https://example.com/avatar.jpg"
		repoResults := []repository.HighlightedSeriesResult{
			{
				Series: &entity.Series{
					ID:          seriesID,
					Title:       "Test Series",
					Slug:        "test-series",
					Description: "Test Description",
					AuthorID:    authorID,
					CreatedAt:   now,
				},
				Author: &entity.User{
					ID:        authorID,
					Name:      "Test Author",
					AvatarURL: &avatarURL,
				},
				SubscriberCount: 100,
				BlogCount:       5,
			},
			{
				Series: &entity.Series{
					ID:       uuid.New(),
					Title:    "Another Series",
					AuthorID: authorID,
				},
				Author: &entity.User{
					ID:   authorID,
					Name: "Test Author",
				}, // Avatar is nil/empty
				SubscriberCount: 50,
				BlogCount:       2,
			},
		}

		mockRepo.EXPECT().GetHighlighted(gomock.Any(), 10).Return(repoResults, nil)

		results, err := uc.GetHighlightedSeries(context.Background())

		assert.NoError(t, err)
		assert.Len(t, results, 2)

		// Check first result
		assert.Equal(t, seriesID, results[0].ID)
		assert.Equal(t, "Test Series", results[0].Title)
		assert.Equal(t, "Test Author", results[0].AuthorName)
		assert.Equal(t, "https://example.com/avatar.jpg", *results[0].AuthorAvatarURL)
		assert.Equal(t, 100, results[0].SubscriberCount)
		assert.Equal(t, 5, results[0].BlogCount)

		// Check second result (nil avatar)
		assert.Equal(t, "Another Series", results[1].Title)
		assert.Nil(t, results[1].AuthorAvatarURL)
	})

	t.Run("success_nil_author", func(t *testing.T) {
		// This test case simulates a scenario where the repository returns a result
		// with a nil Author (e.g., deleted user or failed join), but valid Series data.
		// The usecase should handle this gracefully without panicking.
		seriesID := uuid.New()
		authorID := uuid.New()

		repoResults := []repository.HighlightedSeriesResult{
			{
				Series: &entity.Series{
					ID:       seriesID,
					Title:    "Orphan Series",
					AuthorID: authorID,
				},
				Author:          nil, // Simulate missing author
				SubscriberCount: 10,
				BlogCount:       1,
			},
		}

		mockRepo.EXPECT().GetHighlighted(gomock.Any(), 10).Return(repoResults, nil)

		results, err := uc.GetHighlightedSeries(context.Background())

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, seriesID, results[0].ID)
		assert.Equal(t, authorID, results[0].AuthorID) // Should fallback to Series.AuthorID
		assert.Empty(t, results[0].AuthorName)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetHighlighted(gomock.Any(), 10).Return(nil, errors.New("db error"))

		results, err := uc.GetHighlightedSeries(context.Background())

		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, "db error", err.Error())
	})
}
