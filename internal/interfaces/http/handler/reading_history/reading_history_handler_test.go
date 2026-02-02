package reading_history_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/interfaces/http/handler/reading_history"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReadingHistoryUseCase
type MockReadingHistoryUseCase struct {
	mock.Mock
}

func (m *MockReadingHistoryUseCase) MarkAsRead(ctx context.Context, userID, blogID uuid.UUID) error {
	args := m.Called(ctx, userID, blogID)
	return args.Error(0)
}

func (m *MockReadingHistoryUseCase) GetHistory(ctx context.Context, userID uuid.UUID, limit int) (*dto.ReadingHistoryListResponse, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReadingHistoryListResponse), args.Error(1)
}

func TestReadingHistoryHandler_MarkAsRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := new(MockReadingHistoryUseCase)
		handler := reading_history.NewReadingHistoryHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		blogID := uuid.New()

		// Simulate authenticated user
		c.Set("userID", userID)
		c.Params = []gin.Param{{Key: "id", Value: blogID.String()}}

		c.Request, _ = http.NewRequest(http.MethodPost, "/blogs/"+blogID.String()+"/read", nil)

		mockUC.On("MarkAsRead", mock.Anything, userID, blogID).Return(nil)

		handler.MarkAsRead(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUC.AssertExpectations(t)
	})
}
