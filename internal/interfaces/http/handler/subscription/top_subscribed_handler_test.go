package subscription

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	usecasemocks "github.com/aiagent/internal/application/usecase/subscription/mocks"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubscriptionHandler_GetTopSubscribedAuthors(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("defaults_and_response_shape", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := usecasemocks.NewMockSubscriptionUseCase(ctrl)
		h := NewSubscriptionHandler(mockUC)

		authorID := uuid.New()
		mockUC.EXPECT().
			GetTopSubscribedAuthors(gomock.Any(), 1, 20).
			DoAndReturn(func(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.TopSubscribedAuthorResponse], error) {
				assert.NotNil(t, ctx)
				return &repository.PaginatedResult[dto.TopSubscribedAuthorResponse]{
					Data: []dto.TopSubscribedAuthorResponse{
						{
							AuthorID:        authorID,
							Username:        "author_one",
							DisplayName:     "Author One",
							AvatarURL:       "https://cdn.example.com/avatar.jpg",
							SubscriberCount: 42,
						},
					},
					Page:       1,
					PageSize:   20,
					Total:      100,
					TotalPages: 5,
				}, nil
			})

		r := gin.New()
		r.GET("/api/v1/authors/top-subscribed", h.GetTopSubscribedAuthors)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		if assert.NotNil(t, resp.Meta) {
			assert.Equal(t, 1, resp.Meta.Page)
			assert.Equal(t, 20, resp.Meta.PageSize)
			assert.Equal(t, int64(100), resp.Meta.Total)
			assert.Equal(t, 5, resp.Meta.TotalPages)
		}

		var rawResp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)
		metaMap, ok := rawResp["meta"].(map[string]interface{})
		if assert.True(t, ok) {
			assert.Equal(t, float64(100), metaMap["totalItems"])
		}

		data, ok := resp.Data.([]interface{})
		if assert.True(t, ok) && assert.Len(t, data, 1) {
			item, ok := data[0].(map[string]interface{})
			if assert.True(t, ok) {
				assert.Equal(t, authorID.String(), item["authorId"])
				assert.Equal(t, "author_one", item["username"])
				assert.Equal(t, "Author One", item["displayName"])
				assert.Equal(t, "https://cdn.example.com/avatar.jpg", item["avatarUrl"])
				assert.Equal(t, float64(42), item["subscriberCount"])
			}
		}
	})

	t.Run("invalid_page_returns_bad_request", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := usecasemocks.NewMockSubscriptionUseCase(ctrl)
		h := NewSubscriptionHandler(mockUC)

		r := gin.New()
		r.GET("/api/v1/authors/top-subscribed", h.GetTopSubscribedAuthors)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed?page=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Success)
		if assert.NotNil(t, resp.Error) {
			assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
		}
	})

	t.Run("invalid_page_size_returns_bad_request", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := usecasemocks.NewMockSubscriptionUseCase(ctrl)
		h := NewSubscriptionHandler(mockUC)

		r := gin.New()
		r.GET("/api/v1/authors/top-subscribed", h.GetTopSubscribedAuthors)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed?pageSize=101", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Success)
		if assert.NotNil(t, resp.Error) {
			assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
		}
	})

	t.Run("usecase_error_returns_internal_error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := usecasemocks.NewMockSubscriptionUseCase(ctrl)
		h := NewSubscriptionHandler(mockUC)

		mockUC.EXPECT().
			GetTopSubscribedAuthors(gomock.Any(), 2, 10).
			Return(nil, errors.New("db failure"))

		r := gin.New()
		r.GET("/api/v1/authors/top-subscribed", h.GetTopSubscribedAuthors)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed?page=2&pageSize=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Success)
		if assert.NotNil(t, resp.Error) {
			assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
			assert.Equal(t, "db failure", resp.Error.Message)
		}
	})
}
