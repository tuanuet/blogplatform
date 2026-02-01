package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeriesRepository is a mock implementation of repository.SeriesRepository
type MockSeriesRepository struct {
	mock.Mock
}

func (m *MockSeriesRepository) Create(ctx context.Context, series *entity.Series) error {
	args := m.Called(ctx, series)
	return args.Error(0)
}

func (m *MockSeriesRepository) Update(ctx context.Context, series *entity.Series) error {
	args := m.Called(ctx, series)
	return args.Error(0)
}

func (m *MockSeriesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSeriesRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Series, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Series), args.Error(1)
}

func (m *MockSeriesRepository) GetBySlug(ctx context.Context, slug string) (*entity.Series, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Series), args.Error(1)
}

func (m *MockSeriesRepository) List(ctx context.Context, params map[string]interface{}) ([]entity.Series, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]entity.Series), args.Get(1).(int64), args.Error(2)
}

func (m *MockSeriesRepository) AddBlog(ctx context.Context, seriesID, blogID uuid.UUID) error {
	args := m.Called(ctx, seriesID, blogID)
	return args.Error(0)
}

func (m *MockSeriesRepository) RemoveBlog(ctx context.Context, seriesID, blogID uuid.UUID) error {
	args := m.Called(ctx, seriesID, blogID)
	return args.Error(0)
}

func TestCreateSeries(t *testing.T) {
	mockRepo := new(MockSeriesRepository)
	uc := usecase.NewSeriesUseCase(mockRepo)
	userID := uuid.New()

	req := &dto.CreateSeriesRequest{
		Title:       "My Series",
		Slug:        "my-series",
		Description: "A description",
	}

	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*entity.Series")).Return(nil)

	resp, err := uc.CreateSeries(context.Background(), userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Slug, resp.Slug)
	assert.Equal(t, userID, resp.AuthorID)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSeries_Success(t *testing.T) {
	mockRepo := new(MockSeriesRepository)
	uc := usecase.NewSeriesUseCase(mockRepo)
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

	mockRepo.On("GetByID", context.Background(), seriesID).Return(existingSeries, nil)
	mockRepo.On("Update", context.Background(), mock.MatchedBy(func(s *entity.Series) bool {
		return s.ID == seriesID && s.Title == "New Title" && s.Description == "New Description"
	})).Return(nil)

	resp, err := uc.UpdateSeries(context.Background(), userID, seriesID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Title", resp.Title)
	assert.Equal(t, "New Description", resp.Description)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSeries_Unauthorized(t *testing.T) {
	mockRepo := new(MockSeriesRepository)
	uc := usecase.NewSeriesUseCase(mockRepo)
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

	mockRepo.On("GetByID", context.Background(), seriesID).Return(existingSeries, nil)

	resp, err := uc.UpdateSeries(context.Background(), userID, seriesID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "unauthorized: you are not the author of this series", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDeleteSeries_Success(t *testing.T) {
	mockRepo := new(MockSeriesRepository)
	uc := usecase.NewSeriesUseCase(mockRepo)
	userID := uuid.New()
	seriesID := uuid.New()

	existingSeries := &entity.Series{
		ID:       seriesID,
		AuthorID: userID,
	}

	mockRepo.On("GetByID", context.Background(), seriesID).Return(existingSeries, nil)
	mockRepo.On("Delete", context.Background(), seriesID).Return(nil)

	err := uc.DeleteSeries(context.Background(), userID, seriesID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
