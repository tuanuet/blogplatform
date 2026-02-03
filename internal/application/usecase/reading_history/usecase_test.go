package reading_history_test

import (
	"context"
	"testing"
	"time"

	readinghistory "github.com/aiagent/internal/application/usecase/reading_history"
	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReadingHistoryUseCase_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockReadingHistoryRepository(ctrl)
	uc := readinghistory.NewReadingHistoryUseCase(mockRepo)

	userID := uuid.New()
	blogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)

		err := uc.MarkAsRead(context.Background(), userID, blogID)
		assert.NoError(t, err)
	})
}

func TestReadingHistoryUseCase_GetHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockReadingHistoryRepository(ctrl)
	uc := readinghistory.NewReadingHistoryUseCase(mockRepo)

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

		mockRepo.EXPECT().GetRecentByUserID(gomock.Any(), userID, 20).Return(histories, nil)

		res, err := uc.GetHistory(context.Background(), userID, 20)
		assert.NoError(t, err)
		assert.Len(t, res.History, 1)
	})
}
