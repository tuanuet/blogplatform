package service_test

import (
	"context"
	"errors"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTagTierService_AssignTagToTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewTagTierService(mockTagTierRepo, mockTagRepo, mockBlogRepo)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	t.Run("successfully_assign_tag_to_tier", func(t *testing.T) {
		mockTagRepo.EXPECT().FindByID(ctx, tagID).Return(&entity.Tag{ID: tagID, Name: "premium"}, nil)
		mockBlogRepo.EXPECT().ExistsByAuthorAndTag(ctx, authorID, tagID).Return(true, nil)
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tagID).Return(int64(5), nil)
		mockTagTierRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Do(func(_ context.Context, m *entity.TagTierMapping) {
			assert.Equal(t, authorID, m.AuthorID)
			assert.Equal(t, tagID, m.TagID)
			assert.Equal(t, entity.TierSilver, m.RequiredTier)
		}).Return(nil)

		mapping, blogCount, err := svc.AssignTagToTier(ctx, authorID, tagID, entity.TierSilver)

		assert.NoError(t, err)
		assert.NotNil(t, mapping)
		assert.Equal(t, int64(5), blogCount)
		assert.Equal(t, entity.TierSilver, mapping.RequiredTier)
	})

	t.Run("tag_not_found", func(t *testing.T) {
		mockTagRepo.EXPECT().FindByID(ctx, tagID).Return(nil, errors.New("tag not found"))

		_, _, err := svc.AssignTagToTier(ctx, authorID, tagID, entity.TierSilver)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("tag_not_used_by_author", func(t *testing.T) {
		mockTagRepo.EXPECT().FindByID(ctx, tagID).Return(&entity.Tag{ID: tagID, Name: "premium"}, nil)
		mockBlogRepo.EXPECT().ExistsByAuthorAndTag(ctx, authorID, tagID).Return(false, nil)

		_, _, err := svc.AssignTagToTier(ctx, authorID, tagID, entity.TierSilver)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not used by author")
	})

	t.Run("cannot_assign_free_tier", func(t *testing.T) {
		_, _, err := svc.AssignTagToTier(ctx, authorID, tagID, entity.TierFree)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot assign FREE tier")
	})

	t.Run("upsert_failure", func(t *testing.T) {
		mockTagRepo.EXPECT().FindByID(ctx, tagID).Return(&entity.Tag{ID: tagID, Name: "premium"}, nil)
		mockBlogRepo.EXPECT().ExistsByAuthorAndTag(ctx, authorID, tagID).Return(true, nil)
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tagID).Return(int64(5), nil)
		mockTagTierRepo.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		_, _, err := svc.AssignTagToTier(ctx, authorID, tagID, entity.TierSilver)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to assign tag")
	})
}

func TestTagTierService_UnassignTagFromTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewTagTierService(mockTagTierRepo, mockTagRepo, mockBlogRepo)

	ctx := context.Background()
	authorID := uuid.New()
	tagID := uuid.New()

	t.Run("successfully_unassign_tag", func(t *testing.T) {
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tagID).Return(int64(0), nil)
		mockTagTierRepo.EXPECT().Delete(ctx, authorID, tagID).Return(nil)

		blogCount, err := svc.UnassignTagFromTier(ctx, authorID, tagID)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), blogCount)
	})

	t.Run("delete_failure", func(t *testing.T) {
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tagID).Return(int64(0), nil)
		mockTagTierRepo.EXPECT().Delete(ctx, authorID, tagID).Return(errors.New("database error"))

		_, err := svc.UnassignTagFromTier(ctx, authorID, tagID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unassign tag")
	})
}

func TestTagTierService_GetRequiredTierForBlog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewTagTierService(mockTagTierRepo, mockTagRepo, mockBlogRepo)

	ctx := context.Background()
	blogID := uuid.New()
	authorID := uuid.New()
	tag1ID := uuid.New()
	tag2ID := uuid.New()
	tag3ID := uuid.New()

	t.Run("mixed_tiers_returns_highest", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: authorID,
			Title:    "Test Blog",
			Tags: []entity.Tag{
				{ID: tag1ID, Name: "free"},
				{ID: tag2ID, Name: "premium"},
				{ID: tag3ID, Name: "exclusive"},
			},
		}

		mappings := []entity.TagTierMapping{
			{AuthorID: authorID, TagID: tag1ID, RequiredTier: entity.TierFree},
			{AuthorID: authorID, TagID: tag2ID, RequiredTier: entity.TierSilver},
			{AuthorID: authorID, TagID: tag3ID, RequiredTier: entity.TierGold},
		}

		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)
		mockTagTierRepo.EXPECT().FindByTagIDs(ctx, authorID, []uuid.UUID{tag1ID, tag2ID, tag3ID}).Return(mappings, nil)

		tier, err := svc.GetRequiredTierForBlog(ctx, blogID)

		assert.NoError(t, err)
		assert.Equal(t, entity.TierGold, tier)
	})

	t.Run("blog_with_no_tags_returns_free", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: authorID,
			Title:    "Test Blog",
			Tags:     []entity.Tag{},
		}

		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)

		tier, err := svc.GetRequiredTierForBlog(ctx, blogID)

		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, tier)
	})

	t.Run("blog_with_unassigned_tags_returns_free", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: authorID,
			Title:    "Test Blog",
			Tags: []entity.Tag{
				{ID: tag1ID, Name: "unassigned"},
				{ID: tag2ID, Name: "also-unassigned"},
			},
		}

		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)
		mockTagTierRepo.EXPECT().FindByTagIDs(ctx, authorID, []uuid.UUID{tag1ID, tag2ID}).Return([]entity.TagTierMapping{}, nil)

		tier, err := svc.GetRequiredTierForBlog(ctx, blogID)

		assert.NoError(t, err)
		assert.Equal(t, entity.TierFree, tier)
	})

	t.Run("all_tiers_same_level", func(t *testing.T) {
		blog := &entity.Blog{
			ID:       blogID,
			AuthorID: authorID,
			Title:    "Test Blog",
			Tags: []entity.Tag{
				{ID: tag1ID, Name: "silver1"},
				{ID: tag2ID, Name: "silver2"},
			},
		}

		mappings := []entity.TagTierMapping{
			{AuthorID: authorID, TagID: tag1ID, RequiredTier: entity.TierSilver},
			{AuthorID: authorID, TagID: tag2ID, RequiredTier: entity.TierSilver},
		}

		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(blog, nil)
		mockTagTierRepo.EXPECT().FindByTagIDs(ctx, authorID, []uuid.UUID{tag1ID, tag2ID}).Return(mappings, nil)

		tier, err := svc.GetRequiredTierForBlog(ctx, blogID)

		assert.NoError(t, err)
		assert.Equal(t, entity.TierSilver, tier)
	})

	t.Run("blog_not_found", func(t *testing.T) {
		mockBlogRepo.EXPECT().FindByID(ctx, blogID).Return(nil, errors.New("blog not found"))

		_, err := svc.GetRequiredTierForBlog(ctx, blogID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "blog not found")
	})
}

func TestTagTierService_GetAuthorTagTiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagTierRepo := mocks.NewMockTagTierMappingRepository(ctrl)
	mockTagRepo := mocks.NewMockTagRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)

	svc := service.NewTagTierService(mockTagTierRepo, mockTagRepo, mockBlogRepo)

	ctx := context.Background()
	authorID := uuid.New()
	tag1ID := uuid.New()
	tag2ID := uuid.New()
	tag3ID := uuid.New()

	t.Run("returns_mappings_with_tag_names_and_counts", func(t *testing.T) {
		mappings := []entity.TagTierMapping{
			{ID: uuid.New(), AuthorID: authorID, TagID: tag1ID, RequiredTier: entity.TierBronze},
			{ID: uuid.New(), AuthorID: authorID, TagID: tag2ID, RequiredTier: entity.TierSilver},
			{ID: uuid.New(), AuthorID: authorID, TagID: tag3ID, RequiredTier: entity.TierGold},
		}

		tags := []entity.Tag{
			{ID: tag1ID, Name: "bronze-tag"},
			{ID: tag2ID, Name: "silver-tag"},
			{ID: tag3ID, Name: "gold-tag"},
		}

		mockTagTierRepo.EXPECT().FindByAuthor(ctx, authorID).Return(mappings, nil)
		mockTagRepo.EXPECT().FindByIDs(ctx, []uuid.UUID{tag1ID, tag2ID, tag3ID}).Return(tags, nil)
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tag1ID).Return(int64(5), nil)
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tag2ID).Return(int64(3), nil)
		mockTagTierRepo.EXPECT().CountBlogsByTagAndAuthor(ctx, authorID, tag3ID).Return(int64(1), nil)

		result, err := svc.GetAuthorTagTiers(ctx, authorID)

		assert.NoError(t, err)
		assert.Len(t, result, 3)

		// Find bronze mapping
		var bronzeTag *service.TagTierWithCount
		for i, r := range result {
			if r.TagName == "bronze-tag" {
				bronzeTag = &result[i]
			}
		}

		assert.NotNil(t, bronzeTag)
		assert.Equal(t, entity.TierBronze, bronzeTag.Mapping.RequiredTier)
		assert.Equal(t, "bronze-tag", bronzeTag.TagName)
		assert.Equal(t, int64(5), bronzeTag.BlogCount)
	})

	t.Run("returns_empty_list_on_error", func(t *testing.T) {
		mockTagTierRepo.EXPECT().FindByAuthor(ctx, authorID).Return(nil, errors.New("database error"))

		_, err := svc.GetAuthorTagTiers(ctx, authorID)

		assert.Error(t, err)
	})
}
