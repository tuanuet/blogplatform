package version_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/service/mocks"
	"github.com/aiagent/internal/interfaces/http/handler/version"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupRouter() (*gin.Engine, *mocks.MockVersionService, *mocks.MockBlogService, *version.VersionHandler, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(nil)
	mockVersionService := mocks.NewMockVersionService(ctrl)
	mockBlogService := mocks.NewMockBlogService(ctrl)
	handler := version.NewVersionHandler(mockVersionService, mockBlogService)

	userID := uuid.New()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	return r, mockVersionService, mockBlogService, handler, userID
}

func TestVersionHandler_List(t *testing.T) {
	r, mockService, mockBlogService, handler, userID := setupRouter()

	blogID := uuid.New()
	url := "/blogs/" + blogID.String() + "/versions"
	r.GET("/blogs/:id/versions", handler.List)

	t.Run("Success", func(t *testing.T) {
		versions := []entity.BlogVersion{
			{
				ID:            uuid.New(),
				BlogID:        blogID,
				VersionNumber: 1,
				Title:         "Version 1",
				CreatedAt:     time.Now(),
			},
		}

		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: userID,
		}

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		mockService.EXPECT().
			ListVersions(gomock.Any(), blogID, gomock.Any()).
			Return(&repository.PaginatedResult[entity.BlogVersion]{
				Data:       versions,
				Total:      1,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			}, nil)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data []dto.VersionResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Data))
		assert.Equal(t, 1, resp.Data[0].VersionNumber)
	})

	t.Run("Forbidden", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: uuid.New(), // Different author
		}

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/blogs/invalid-uuid/versions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: userID,
		}

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		mockService.EXPECT().
			ListVersions(gomock.Any(), blogID, gomock.Any()).
			Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestVersionHandler_Get(t *testing.T) {
	r, mockService, mockBlogService, handler, userID := setupRouter()

	blogID := uuid.New()
	versionID := uuid.New()
	url := "/blogs/" + blogID.String() + "/versions/" + versionID.String()
	r.GET("/blogs/:id/versions/:versionId", handler.Get)

	t.Run("Success", func(t *testing.T) {
		ver := &entity.BlogVersion{
			ID:            versionID,
			BlogID:        blogID,
			VersionNumber: 2,
			Title:         "Version 2",
			Content:       "Content 2",
			CreatedAt:     time.Now(),
		}

		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: userID,
		}

		mockService.EXPECT().
			GetVersion(gomock.Any(), versionID).
			Return(ver, nil)

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.VersionDetailResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, resp.Data.VersionNumber)
		assert.Equal(t, "Content 2", resp.Data.Content)
	})

	t.Run("Forbidden", func(t *testing.T) {
		ver := &entity.BlogVersion{
			ID:     versionID,
			BlogID: blogID,
		}

		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: uuid.New(), // Different author
		}

		mockService.EXPECT().
			GetVersion(gomock.Any(), versionID).
			Return(ver, nil)

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().
			GetVersion(gomock.Any(), versionID).
			Return(nil, service.ErrVersionNotFound)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestVersionHandler_Create(t *testing.T) {
	r, mockVersionService, mockBlogService, handler, userID := setupRouter()

	blogID := uuid.New()
	url := "/blogs/" + blogID.String() + "/versions"
	r.POST("/blogs/:id/versions", handler.Create)

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.CreateVersionRequest{
			ChangeSummary: "Updated title",
		}
		body, _ := json.Marshal(reqBody)

		blog := &entity.Blog{
			ID:       blogID,
			Title:    "Original Blog",
			AuthorID: userID,
		}

		ver := &entity.BlogVersion{
			ID:            uuid.New(),
			BlogID:        blogID,
			VersionNumber: 3,
			ChangeSummary: &reqBody.ChangeSummary,
		}

		mockBlogService.EXPECT().
			GetByID(gomock.Any(), blogID, gomock.Any()).
			Return(blog, nil)

		mockVersionService.EXPECT().
			CreateVersion(gomock.Any(), blog, gomock.Any(), "Updated title").
			Return(ver, nil)

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVersionHandler_Restore(t *testing.T) {
	r, mockVersionService, _, handler, _ := setupRouter()

	blogID := uuid.New()
	versionID := uuid.New()
	url := "/blogs/" + blogID.String() + "/versions/" + versionID.String() + "/restore"
	r.POST("/blogs/:id/versions/:versionId/restore", handler.Restore)

	t.Run("Success", func(t *testing.T) {
		updatedBlog := &entity.Blog{
			ID:    blogID,
			Title: "Restored Blog",
		}

		mockVersionService.EXPECT().
			RestoreVersion(gomock.Any(), blogID, versionID, gomock.Any()).
			Return(updatedBlog, nil)

		req, _ := http.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data dto.BlogResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Restored Blog", resp.Data.Title)
	})
}

func TestVersionHandler_Delete(t *testing.T) {
	r, mockVersionService, _, handler, _ := setupRouter()

	blogID := uuid.New()
	versionID := uuid.New()
	url := "/blogs/" + blogID.String() + "/versions/" + versionID.String()
	r.DELETE("/blogs/:id/versions/:versionId", handler.Delete)

	t.Run("Success", func(t *testing.T) {
		mockVersionService.EXPECT().
			DeleteVersion(gomock.Any(), versionID, gomock.Any()).
			Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
