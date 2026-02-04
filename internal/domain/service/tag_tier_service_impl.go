package service

import (
	"context"
	"fmt"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

type tagTierService struct {
	tagTierRepo repository.TagTierMappingRepository
	tagRepo     repository.TagRepository
	blogRepo    repository.BlogRepository
}

// NewTagTierService creates a new TagTierService instance
func NewTagTierService(
	tagTierRepo repository.TagTierMappingRepository,
	tagRepo repository.TagRepository,
	blogRepo repository.BlogRepository,
) TagTierService {
	return &tagTierService{
		tagTierRepo: tagTierRepo,
		tagRepo:     tagRepo,
		blogRepo:    blogRepo,
	}
}

// AssignTagToTier assigns a tag to a subscription tier
func (s *tagTierService) AssignTagToTier(
	ctx context.Context,
	authorID, tagID uuid.UUID,
	tier entity.SubscriptionTier,
) (*entity.TagTierMapping, int64, error) {
	// Validate tier
	if tier == entity.TierFree {
		return nil, 0, fmt.Errorf("cannot assign FREE tier to tag")
	}

	// Check if tag exists
	tag, err := s.tagRepo.FindByID(ctx, tagID)
	if err != nil {
		return nil, 0, fmt.Errorf("tag not found: %w", err)
	}
	if tag == nil {
		return nil, 0, fmt.Errorf("tag not found")
	}

	// Check if tag is used by the author
	exists, err := s.blogRepo.ExistsByAuthorAndTag(ctx, authorID, tagID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check tag usage: %w", err)
	}
	if !exists {
		return nil, 0, fmt.Errorf("tag '%s' is not used by author in any blog", tag.Name)
	}

	// Count affected blogs
	blogCount, err := s.tagTierRepo.CountBlogsByTagAndAuthor(ctx, authorID, tagID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
	}

	// Create or update mapping
	mapping := &entity.TagTierMapping{
		AuthorID:     authorID,
		TagID:        tagID,
		RequiredTier: tier,
	}

	if err := s.tagTierRepo.Upsert(ctx, mapping); err != nil {
		return nil, 0, fmt.Errorf("failed to assign tag to tier: %w", err)
	}

	return mapping, blogCount, nil
}

// UnassignTagFromTier removes tier requirement from a tag
func (s *tagTierService) UnassignTagFromTier(
	ctx context.Context,
	authorID, tagID uuid.UUID,
) (int64, error) {
	// Count affected blogs before deletion
	blogCount, err := s.tagTierRepo.CountBlogsByTagAndAuthor(ctx, authorID, tagID)
	if err != nil {
		return 0, fmt.Errorf("failed to count blogs: %w", err)
	}

	// Delete mapping
	if err := s.tagTierRepo.Delete(ctx, authorID, tagID); err != nil {
		return 0, fmt.Errorf("failed to unassign tag: %w", err)
	}

	return blogCount, nil
}

// GetRequiredTierForBlog determines the highest required tier for a blog based on its tags
// Returns: required tier, blog author ID, error
func (s *tagTierService) GetRequiredTierForBlog(
	ctx context.Context,
	blogID uuid.UUID,
) (entity.SubscriptionTier, uuid.UUID, error) {
	// Find blog
	blog, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return entity.TierFree, uuid.Nil, fmt.Errorf("blog not found: %w", err)
	}
	if blog == nil {
		return entity.TierFree, uuid.Nil, fmt.Errorf("blog not found")
	}

	// If no tags, return FREE with blog's author ID and nil error
	if len(blog.Tags) == 0 {
		return entity.TierFree, blog.AuthorID, nil
	}

	// Extract tag IDs
	tagIDs := make([]uuid.UUID, len(blog.Tags))
	for i, tag := range blog.Tags {
		tagIDs[i] = tag.ID
	}

	// Find tag-tier mappings
	mappings, err := s.tagTierRepo.FindByTagIDs(ctx, blog.AuthorID, tagIDs)
	if err != nil {
		return entity.TierFree, blog.AuthorID, fmt.Errorf("failed to find tag mappings: %w", err)
	}

	// Find highest tier
	highestTier := entity.TierFree
	for _, m := range mappings {
		if m.RequiredTier.Level() > highestTier.Level() {
			highestTier = m.RequiredTier
		}
	}

	return highestTier, blog.AuthorID, nil
}

// GetAuthorTagTiers retrieves all tag-tier mappings for an author
func (s *tagTierService) GetAuthorTagTiers(
	ctx context.Context,
	authorID uuid.UUID,
) ([]TagTierWithCount, error) {
	// Find all mappings
	mappings, err := s.tagTierRepo.FindByAuthor(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag tiers: %w", err)
	}

	// If no mappings, return empty
	if len(mappings) == 0 {
		return []TagTierWithCount{}, nil
	}

	// Get tag details
	tagIDs := make([]uuid.UUID, len(mappings))
	for i, m := range mappings {
		tagIDs[i] = m.TagID
	}

	tags, err := s.tagRepo.FindByIDs(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find tags: %w", err)
	}

	// Build tag name lookup
	tagNameMap := make(map[uuid.UUID]string)
	for _, tag := range tags {
		tagNameMap[tag.ID] = tag.Name
	}

	// Build result with blog counts
	result := make([]TagTierWithCount, len(mappings))
	for i, m := range mappings {
		blogCount, err := s.tagTierRepo.CountBlogsByTagAndAuthor(ctx, authorID, m.TagID)
		if err != nil {
			return nil, fmt.Errorf("failed to count blogs for tag %s: %w", m.TagID, err)
		}

		result[i] = TagTierWithCount{
			Mapping:   m,
			TagName:   tagNameMap[m.TagID],
			BlogCount: blogCount,
		}
	}

	return result, nil
}
