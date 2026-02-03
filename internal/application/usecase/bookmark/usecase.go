package bookmark

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// BookmarkUseCase defines the interface for bookmark business logic
type BookmarkUseCase interface {
	// BookmarkBlog adds a bookmark for a user
	BookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error

	// UnbookmarkBlog removes a bookmark for a user
	UnbookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error

	// GetUserBookmarks returns a paginated list of bookmarks for a user
	GetUserBookmarks(ctx context.Context, userID uuid.UUID, params *dto.BlogFilterParams) (*repository.PaginatedResult[dto.BlogListResponse], error)
}

type bookmarkUseCase struct {
	bookmarkRepo repository.BookmarkRepository
}

// NewBookmarkUseCase creates a new instance of BookmarkUseCase
func NewBookmarkUseCase(bookmarkRepo repository.BookmarkRepository) BookmarkUseCase {
	return &bookmarkUseCase{
		bookmarkRepo: bookmarkRepo,
	}
}

func (u *bookmarkUseCase) BookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error {
	// Check if already bookmarked? The repo.Add might handle it or return error.
	// We'll just call the repo.
	return u.bookmarkRepo.Add(ctx, userID, blogID)
}

func (u *bookmarkUseCase) UnbookmarkBlog(ctx context.Context, userID, blogID uuid.UUID) error {
	return u.bookmarkRepo.Remove(ctx, userID, blogID)
}

func (u *bookmarkUseCase) GetUserBookmarks(ctx context.Context, userID uuid.UUID, params *dto.BlogFilterParams) (*repository.PaginatedResult[dto.BlogListResponse], error) {
	pagination := repository.Pagination{
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	result, err := u.bookmarkRepo.FindByUser(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	// Map entity.Blog to dto.BlogListResponse
	dtos := make([]dto.BlogListResponse, len(result.Data))
	for i, blog := range result.Data {
		dtos[i] = u.mapBlogToDTO(&blog)
	}

	return &repository.PaginatedResult[dto.BlogListResponse]{
		Data:       dtos,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (u *bookmarkUseCase) mapBlogToDTO(blog *entity.Blog) dto.BlogListResponse {
	resp := dto.BlogListResponse{
		ID:            blog.ID,
		AuthorID:      blog.AuthorID,
		Title:         blog.Title,
		Slug:          blog.Slug,
		Excerpt:       blog.Excerpt,
		ThumbnailURL:  blog.ThumbnailURL,
		Status:        blog.Status,
		Visibility:    blog.Visibility,
		PublishedAt:   blog.PublishedAt,
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
		resp.CategoryID = &blog.Category.ID
		resp.Category = &dto.CategoryResponse{
			ID:   blog.Category.ID,
			Name: blog.Category.Name,
			Slug: blog.Category.Slug,
		}
	}

	if len(blog.Tags) > 0 {
		resp.Tags = make([]dto.TagResponse, len(blog.Tags))
		for i, tag := range blog.Tags {
			resp.Tags[i] = dto.TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
				Slug: tag.Slug,
			}
		}
	} else {
		resp.Tags = []dto.TagResponse{}
	}

	return resp
}
