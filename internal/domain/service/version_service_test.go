package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVersionService_CreateVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mocks.NewMockBlogVersionRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewVersionService(mockVersionRepo, mockBlogRepo)

	ctx := context.Background()
	blogID := uuid.New()
	editorID := uuid.New()
	blog := &entity.Blog{
		ID:       blogID,
		Title:    "Test Blog",
		Content:  "Content",
		AuthorID: editorID,
	}

	t.Run("success", func(t *testing.T) {
		nextVersion := 1
		mockVersionRepo.EXPECT().GetNextVersionNumber(ctx, blogID).Return(nextVersion, nil)
		mockVersionRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, v *entity.BlogVersion) error {
			assert.Equal(t, blogID, v.BlogID)
			assert.Equal(t, nextVersion, v.VersionNumber)
			assert.Equal(t, "Content", v.Content)
			assert.Equal(t, "Test Blog", v.Title)
			return nil
		})
		mockVersionRepo.EXPECT().DeleteOldest(ctx, blogID, 50).Return(nil)

		v, err := svc.CreateVersion(ctx, blog, editorID, "update")
		assert.NoError(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, nextVersion, v.VersionNumber)
	})

	t.Run("error_get_next_version", func(t *testing.T) {
		mockVersionRepo.EXPECT().GetNextVersionNumber(ctx, blogID).Return(0, errors.New("db error"))

		v, err := svc.CreateVersion(ctx, blog, editorID, "update")
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func TestVersionService_ListVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mocks.NewMockBlogVersionRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewVersionService(mockVersionRepo, mockBlogRepo)

	ctx := context.Background()
	blogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		versions := []entity.BlogVersion{
			{ID: uuid.New(), BlogID: blogID, VersionNumber: 2},
			{ID: uuid.New(), BlogID: blogID, VersionNumber: 1},
		}
		result := &repository.PaginatedResult[entity.BlogVersion]{
			Data: versions,
		}
		mockVersionRepo.EXPECT().FindByBlogID(ctx, blogID, repository.Pagination{Page: 1, PageSize: 10}).Return(result, nil)

		res, err := svc.ListVersions(ctx, blogID, repository.Pagination{Page: 1, PageSize: 10})
		assert.NoError(t, err)
		assert.Equal(t, result, res)
	})
}

func TestVersionService_GetVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mocks.NewMockBlogVersionRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewVersionService(mockVersionRepo, mockBlogRepo)

	ctx := context.Background()
	versionID := uuid.New()

	t.Run("success", func(t *testing.T) {
		version := &entity.BlogVersion{ID: versionID}
		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)

		res, err := svc.GetVersion(ctx, versionID)
		assert.NoError(t, err)
		assert.Equal(t, version, res)
	})

	t.Run("not_found", func(t *testing.T) {
		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(nil, nil)

		res, err := svc.GetVersion(ctx, versionID)
		assert.ErrorIs(t, err, service.ErrVersionNotFound)
		assert.Nil(t, res)
	})
}

func TestVersionService_RestoreVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mocks.NewMockBlogVersionRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewVersionService(mockVersionRepo, mockBlogRepo)

	ctx := context.Background()
	blogID := uuid.New()
	versionID := uuid.New()
	editorID := uuid.New()

	t.Run("success", func(t *testing.T) {
		version := &entity.BlogVersion{
			ID:            versionID,
			BlogID:        blogID,
			Title:         "Old Title",
			Content:       "Old Content",
			VersionNumber: 5,
		}
		blog := &entity.Blog{
			ID:       blogID,
			Title:    "Current Title",
			Content:  "Current Content",
			AuthorID: editorID,
		}

		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)

		// Update blog
		mockBlogRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, b *entity.Blog) error {
			assert.Equal(t, "Old Title", b.Title)
			assert.Equal(t, "Old Content", b.Content)
			return nil
		})

		// Restore tags
		mockBlogRepo.EXPECT().ReplaceTags(ctx, blogID, gomock.Any()).Return(nil)

		// Create new version
		mockVersionRepo.EXPECT().GetNextVersionNumber(ctx, blogID).Return(10, nil)
		mockVersionRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		mockVersionRepo.EXPECT().DeleteOldest(ctx, blogID, 50).Return(nil)

		res, err := svc.RestoreVersion(ctx, blogID, versionID, editorID)
		assert.NoError(t, err)
		assert.Equal(t, "Old Title", res.Title)
		assert.Equal(t, "Old Content", res.Content)
	})

	t.Run("version_mismatch", func(t *testing.T) {
		version := &entity.BlogVersion{
			ID:     versionID,
			BlogID: uuid.New(), // Different blog ID
		}
		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)

		_, err := svc.RestoreVersion(ctx, blogID, versionID, editorID)
		assert.ErrorIs(t, err, service.ErrVersionMismatch)
	})

	t.Run("access_denied", func(t *testing.T) {
		version := &entity.BlogVersion{
			ID:     versionID,
			BlogID: blogID,
		}
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: uuid.New(), // Different author
		}

		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)

		_, err := svc.RestoreVersion(ctx, blogID, versionID, editorID)
		assert.ErrorIs(t, err, service.ErrBlogAccessDenied)
	})
}

func TestVersionService_DeleteVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVersionRepo := mocks.NewMockBlogVersionRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewVersionService(mockVersionRepo, mockBlogRepo)

	ctx := context.Background()
	blogID := uuid.New()
	versionID := uuid.New()
	requesterID := uuid.New()
	authorID := uuid.New()

	t.Run("success_as_author", func(t *testing.T) {
		version := &entity.BlogVersion{ID: versionID, BlogID: blogID}
		blog := &entity.Blog{ID: blogID, AuthorID: requesterID}

		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)
		mockVersionRepo.EXPECT().Delete(ctx, versionID).Return(nil)

		err := svc.DeleteVersion(ctx, versionID, requesterID)
		assert.NoError(t, err)
	})

	t.Run("success_as_editor_in_version", func(t *testing.T) {
		// If requester is not author but created the version?
		// Spec says: "requester == editor OR requester == blog.Author"
		// Assuming version.EditorID is stored.
		version := &entity.BlogVersion{ID: versionID, BlogID: blogID, EditorID: requesterID}
		blog := &entity.Blog{ID: blogID, AuthorID: authorID}

		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)
		mockVersionRepo.EXPECT().Delete(ctx, versionID).Return(nil)

		err := svc.DeleteVersion(ctx, versionID, requesterID)
		assert.NoError(t, err)
	})

	t.Run("access_denied", func(t *testing.T) {
		version := &entity.BlogVersion{ID: versionID, BlogID: blogID, EditorID: uuid.New()}
		blog := &entity.Blog{ID: blogID, AuthorID: authorID}

		mockVersionRepo.EXPECT().FindByID(ctx, versionID).Return(version, nil)
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)

		err := svc.DeleteVersion(ctx, versionID, requesterID)
		assert.Error(t, err)
		// Assuming generic access denied error or specific
	})
}
