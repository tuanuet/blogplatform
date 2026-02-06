package series_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/series/mocks"
	"github.com/aiagent/internal/interfaces/http/handler/series"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetHighlightedSeries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockSeriesUseCase(ctrl)
		handler := series.NewSeriesHandler(mockUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/series/highlighted", nil)

		expectedSeries := []*dto.HighlightedSeriesResponse{
			{
				ID:              uuid.New(),
				Title:           "Test Series",
				Slug:            "test-series",
				Description:     "Test Description",
				AuthorID:        uuid.New(),
				AuthorName:      "Test Author",
				SubscriberCount: 10,
				BlogCount:       5,
				CreatedAt:       time.Now(),
			},
		}

		mockUseCase.EXPECT().GetHighlightedSeries(gomock.Any()).Return(expectedSeries, nil)

		handler.GetHighlightedSeries(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.Len(t, data, 1)
		seriesData := data[0].(map[string]interface{})
		assert.Equal(t, "Test Series", seriesData["title"])
	})

	t.Run("internal error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockSeriesUseCase(ctrl)
		handler := series.NewSeriesHandler(mockUseCase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/series/highlighted", nil)

		mockUseCase.EXPECT().GetHighlightedSeries(gomock.Any()).Return(nil, errors.New("db error"))

		handler.GetHighlightedSeries(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
	})
}
