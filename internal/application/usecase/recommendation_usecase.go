package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
)

type RecommendationUseCase interface {
	GetPopularTags(ctx context.Context, limit int) ([]dto.TagResponse, error)
	GetRelatedBlogs(ctx context.Context, blogID uuid.UUID, limit int) ([]dto.BlogListResponse, error)
	GetPersonalizedFeed(ctx context.Context, userID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.BlogListResponse], error)
	UpdateInterests(ctx context.Context, userID uuid.UUID, tagIDs []string) error
}

type recommendationUseCase struct {
	recService domainService.RecommendationService
}

func NewRecommendationUseCase(recService domainService.RecommendationService) RecommendationUseCase {
	return &recommendationUseCase{
		recService: recService,
	}
}

func (uc *recommendationUseCase) GetPopularTags(ctx context.Context, limit int) ([]dto.TagResponse, error) {
	tags, err := uc.recService.GetPopularTags(ctx, limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		dtos[i] = dto.TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
		}
	}
	return dtos, nil
}

func (uc *recommendationUseCase) GetRelatedBlogs(ctx context.Context, blogID uuid.UUID, limit int) ([]dto.BlogListResponse, error) {
	blogs, err := uc.recService.GetRelatedBlogs(ctx, blogID, limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.BlogListResponse, len(blogs))
	for i, blog := range blogs {
		dtos[i] = uc.toBlogListResponse(&blog)
	}
	return dtos, nil
}

func (uc *recommendationUseCase) GetPersonalizedFeed(ctx context.Context, userID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.BlogListResponse], error) {
	pagination := repository.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	result, err := uc.recService.GetPersonalizedFeed(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BlogListResponse, len(result.Data))
	for i, blog := range result.Data {
		items[i] = uc.toBlogListResponse(&blog)
	}

	return &repository.PaginatedResult[dto.BlogListResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *recommendationUseCase) UpdateInterests(ctx context.Context, userID uuid.UUID, tagIDStrs []string) error {
	tagIDs := make([]uuid.UUID, 0, len(tagIDStrs))
	for _, idStr := range tagIDStrs {
		if id, err := uuid.Parse(idStr); err == nil {
			tagIDs = append(tagIDs, id)
		}
	}
	return uc.recService.UpdateInterests(ctx, userID, tagIDs)
}

func (uc *recommendationUseCase) toBlogListResponse(blog *entity.Blog) dto.BlogListResponse {
	resp := dto.BlogListResponse{
		ID:            blog.ID,
		AuthorID:      blog.AuthorID,
		CategoryID:    blog.CategoryID,
		Title:         blog.Title,
		Slug:          blog.Slug,
		Excerpt:       blog.Excerpt,
		ThumbnailURL:  blog.ThumbnailURL,
		Status:        blog.Status,
		Visibility:    blog.Visibility,
		PublishedAt:   blog.PublishedAt,
		Tags:          make([]dto.TagResponse, 0),
		UpvoteCount:   blog.UpvoteCount,
		DownvoteCount: blog.DownvoteCount,
		CreatedAt:     blog.CreatedAt,
	}

	if blog.Author != nil {
		resp.Author = &dto.UserBriefResponse{
			ID:    blog.Author.ID,
			Name:  blog.Author.Name,
			Email: blog.Author.Email,
		}
	}

	if blog.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:          blog.Category.ID,
			Name:        blog.Category.Name,
			Slug:        blog.Category.Slug,
			Description: blog.Category.Description,
			CreatedAt:   blog.Category.CreatedAt,
			UpdatedAt:   blog.Category.UpdatedAt,
		}
	}

	for _, tag := range blog.Tags {
		resp.Tags = append(resp.Tags, dto.TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
		})
	}
	return resp
}
