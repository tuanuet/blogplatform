package series_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/series"
	"github.com/aiagent/internal/domain/entity"
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
