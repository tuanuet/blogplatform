package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

type RecommendationService interface {
	GetPopularTags(ctx context.Context, limit int) ([]entity.Tag, error)
	GetRelatedBlogs(ctx context.Context, blogID uuid.UUID, limit int) ([]entity.Blog, error)
	GetPersonalizedFeed(ctx context.Context, userID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error)
	UpdateInterests(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error
}

type recommendationService struct {
	blogRepo repository.BlogRepository
	tagRepo  repository.TagRepository
	userRepo repository.UserRepository
}

func NewRecommendationService(
	blogRepo repository.BlogRepository,
	tagRepo repository.TagRepository,
	userRepo repository.UserRepository,
) RecommendationService {
	return &recommendationService{
		blogRepo: blogRepo,
		tagRepo:  tagRepo,
		userRepo: userRepo,
	}
}

func (s *recommendationService) GetPopularTags(ctx context.Context, limit int) ([]entity.Tag, error) {
	return s.tagRepo.FindPopular(ctx, limit)
}

func (s *recommendationService) GetRelatedBlogs(ctx context.Context, blogID uuid.UUID, limit int) ([]entity.Blog, error) {
	return s.blogRepo.FindRelated(ctx, blogID, limit)
}

func (s *recommendationService) UpdateInterests(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error {
	// Check if all tagIDs exist.
	if len(tagIDs) > 0 {
		existingTags, err := s.tagRepo.FindByIDs(ctx, tagIDs)
		if err != nil {
			return err
		}
		if len(existingTags) != len(tagIDs) {
			return errors.New("one or more tags not found")
		}
	}

	return s.userRepo.ReplaceInterests(ctx, userID, tagIDs)
}

func (s *recommendationService) GetPersonalizedFeed(ctx context.Context, userID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error) {
	// 1. Get user interests
	interests, err := s.userRepo.GetInterests(ctx, userID)
	if err != nil {
		return nil, err
	}

	filter := repository.BlogFilter{}

	// 2. If user has interests, filter by them
	if len(interests) > 0 {
		tagIDs := make([]uuid.UUID, len(interests))
		for i, tag := range interests {
			tagIDs[i] = tag.ID
		}
		filter.TagIDs = tagIDs
	}

	// 3. Fallback/Standard logic
	return s.blogRepo.FindAll(ctx, filter, pagination)
}
