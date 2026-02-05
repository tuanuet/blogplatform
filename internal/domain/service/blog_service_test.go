package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBlogService_Create_AutoSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := repoMocks.NewMockBlogRepository(ctrl)
	mockSubRepo := repoMocks.NewMockSubscriptionRepository(ctrl)
	mockTagRepo := repoMocks.NewMockTagRepository(ctrl)
	mockVersionService := serviceMocks.NewMockVersionService(ctrl)

	blog := &entity.Blog{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Test Blog",
		Slug:     "test-blog",
		Content:  "Content",
	}

	ctx := context.Background()

	// Expect Create
	mockBlogRepo.EXPECT().FindBySlug(ctx, blog.AuthorID, blog.Slug).Return(nil, nil)
	mockBlogRepo.EXPECT().Create(ctx, blog).Return(nil)

	// Expect Version Creation
	mockVersionService.EXPECT().CreateVersion(ctx, blog, blog.AuthorID, service.VersionInitial).Return(nil, nil)

	s := service.NewBlogService(mockBlogRepo, mockSubRepo, mockTagRepo, nil, mockVersionService)

	err := s.Create(ctx, blog, nil)
	assert.NoError(t, err)
}

func TestBlogService_Update_AutoSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := repoMocks.NewMockBlogRepository(ctrl)
	mockSubRepo := repoMocks.NewMockSubscriptionRepository(ctrl)
	mockTagRepo := repoMocks.NewMockTagRepository(ctrl)
	mockVersionService := serviceMocks.NewMockVersionService(ctrl)

	blog := &entity.Blog{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Updated Blog",
		Slug:     "updated-blog",
		Content:  "Updated Content",
	}

	ctx := context.Background()

	// Expect Update
	mockBlogRepo.EXPECT().FindBySlug(ctx, blog.AuthorID, blog.Slug).Return(nil, nil)
	mockBlogRepo.EXPECT().Update(ctx, blog).Return(nil)

	// Expect Version Creation
	mockVersionService.EXPECT().CreateVersion(ctx, blog, blog.AuthorID, service.VersionAutoSave).Return(nil, nil)

	s := service.NewBlogService(mockBlogRepo, mockSubRepo, mockTagRepo, nil, mockVersionService)

	err := s.Update(ctx, blog, nil)
	assert.NoError(t, err)
}

func TestBlogService_Update_VersionFailureDoesNotBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlogRepo := repoMocks.NewMockBlogRepository(ctrl)
	mockSubRepo := repoMocks.NewMockSubscriptionRepository(ctrl)
	mockTagRepo := repoMocks.NewMockTagRepository(ctrl)
	mockVersionService := serviceMocks.NewMockVersionService(ctrl)

	blog := &entity.Blog{
		ID:       uuid.New(),
		AuthorID: uuid.New(),
		Title:    "Updated Blog",
		Slug:     "updated-blog",
		Content:  "Updated Content",
	}

	ctx := context.Background()

	// Expect Update to succeed
	mockBlogRepo.EXPECT().FindBySlug(ctx, blog.AuthorID, blog.Slug).Return(nil, nil)
	mockBlogRepo.EXPECT().Update(ctx, blog).Return(nil)

	// Expect Version Creation to fail
	mockVersionService.EXPECT().
		CreateVersion(ctx, blog, blog.AuthorID, service.VersionAutoSave).
		Return(nil, errors.New("version creation failed"))

	s := service.NewBlogService(mockBlogRepo, mockSubRepo, mockTagRepo, nil, mockVersionService)

	// Should still return no error
	err := s.Update(ctx, blog, nil)
	assert.NoError(t, err)
}
