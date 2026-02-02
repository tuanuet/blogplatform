package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReadingHistoryRepository
type MockReadingHistoryRepository struct {
	mock.Mock
}

func (m *MockReadingHistoryRepository) Upsert(ctx context.Context, history *entity.UserReadingHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockReadingHistoryRepository) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.UserReadingHistory, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.UserReadingHistory), args.Error(1)
}

func TestReadingHistoryUseCase_MarkAsRead(t *testing.T) {
	mockRepo := new(MockReadingHistoryRepository)
	uc := usecase.NewReadingHistoryUseCase(mockRepo)

	userID := uuid.New()
	blogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Upsert", mock.Anything, mock.MatchedBy(func(h *entity.UserReadingHistory) bool {
			return h.UserID == userID && h.BlogID == blogID
		})).Return(nil).Once()

		err := uc.MarkAsRead(context.Background(), userID, blogID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestReadingHistoryUseCase_GetHistory(t *testing.T) {
	mockRepo := new(MockReadingHistoryRepository)
	uc := usecase.NewReadingHistoryUseCase(mockRepo)

	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		histories := []*entity.UserReadingHistory{
			{
				UserID:     userID,
				BlogID:     uuid.New(),
				LastReadAt: time.Now(),
				Blog: &entity.Blog{
					Title: "Test Blog",
				},
			},
		}

		mockRepo.On("GetRecentByUserID", mock.Anything, userID, 20).Return(histories, nil).Once()

		res, err := uc.GetHistory(context.Background(), userID, 20)
		assert.NoError(t, err)
		assert.Len(t, res.History, 1)
		mockRepo.AssertExpectations(t)
	})
}
