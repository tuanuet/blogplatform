package blog

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	blogUsecase "github.com/aiagent/internal/application/usecase/blog"
	blogmocks "github.com/aiagent/internal/application/usecase/blog/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBlogHandler_GetBySlug_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := blogmocks.NewMockBlogUseCase(ctrl)
	handler := NewBlogHandler(mockUC)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	authorID := uuid.New()
	slug := "test-slug"
	viewerID := uuid.New()
	c.Params = []gin.Param{{Key: "authorId", Value: authorID.String()}, {Key: "slug", Value: slug}}
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authors/%s/blogs/slug/%s", authorID, slug), nil)
	c.Set("userID", viewerID)

	resp := &dto.BlogResponse{ID: uuid.New(), AuthorID: authorID, Slug: slug, Title: "Test"}
	mockUC.EXPECT().GetBySlug(gomock.Any(), authorID, slug, gomock.AssignableToTypeOf(&uuid.UUID{})).DoAndReturn(
		func(ctx context.Context, a uuid.UUID, s string, viewer *uuid.UUID) (*dto.BlogResponse, error) {
			require.NotNil(t, viewer)
			require.Equal(t, viewerID, *viewer)
			return resp, nil
		},
	)

	handler.GetBySlug(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBlogHandler_GetBySlug_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := blogmocks.NewMockBlogUseCase(ctrl)
	handler := NewBlogHandler(mockUC)

	t.Run("invalid author id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "authorId", Value: "not-a-uuid"}, {Key: "slug", Value: "slug"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/authors/not-a-uuid/blogs/slug/slug", nil)

		handler.GetBySlug(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing slug", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		authorID := uuid.New()
		c.Params = []gin.Param{{Key: "authorId", Value: authorID.String()}, {Key: "slug", Value: "   "}}
		c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authors/%s/blogs/slug/", authorID), nil)

		handler.GetBySlug(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBlogHandler_GetBySlug_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := blogmocks.NewMockBlogUseCase(ctrl)
	handler := NewBlogHandler(mockUC)

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		authorID := uuid.New()
		slug := "missing"
		c.Params = []gin.Param{{Key: "authorId", Value: authorID.String()}, {Key: "slug", Value: slug}}
		c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authors/%s/blogs/slug/%s", authorID, slug), nil)

		mockUC.EXPECT().GetBySlug(gomock.Any(), authorID, slug, (*uuid.UUID)(nil)).Return(nil, blogUsecase.ErrBlogNotFound)

		handler.GetBySlug(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		authorID := uuid.New()
		slug := "premium"
		c.Params = []gin.Param{{Key: "authorId", Value: authorID.String()}, {Key: "slug", Value: slug}}
		c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authors/%s/blogs/slug/%s", authorID, slug), nil)

		mockUC.EXPECT().GetBySlug(gomock.Any(), authorID, slug, (*uuid.UUID)(nil)).Return(nil, blogUsecase.ErrBlogAccessDenied)

		handler.GetBySlug(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
