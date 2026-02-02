package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
)

var (
	ErrBlogNotFound         = domainService.ErrBlogNotFound
	ErrBlogAccessDenied     = domainService.ErrBlogAccessDenied
	ErrBlogAlreadyPublished = domainService.ErrBlogAlreadyPublished
	ErrSlugAlreadyExists    = domainService.ErrSlugAlreadyExists
)

type BlogUseCase interface {
	Create(ctx context.Context, authorID uuid.UUID, req *dto.CreateBlogRequest) (*dto.BlogResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*dto.BlogResponse, error)
	GetBySlug(ctx context.Context, authorID uuid.UUID, slug string, viewerID *uuid.UUID) (*dto.BlogResponse, error)
	List(ctx context.Context, params *dto.BlogFilterParams, viewerID *uuid.UUID) (*repository.PaginatedResult[dto.BlogListResponse], error)
	Update(ctx context.Context, id uuid.UUID, authorID uuid.UUID, req *dto.UpdateBlogRequest) (*dto.BlogResponse, error)
	Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, authorID uuid.UUID, req *dto.PublishBlogRequest) (*dto.BlogResponse, error)
	Unpublish(ctx context.Context, id uuid.UUID, authorID uuid.UUID) (*dto.BlogResponse, error)
	React(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.ReactionRequest) (*dto.ReactionResponse, error)
}

type blogUseCase struct {
	blogSvc domainService.BlogService
}

func NewBlogUseCase(blogSvc domainService.BlogService) BlogUseCase {
	return &blogUseCase{
		blogSvc: blogSvc,
	}
}

func (uc *blogUseCase) Create(ctx context.Context, authorID uuid.UUID, req *dto.CreateBlogRequest) (*dto.BlogResponse, error) {
	blog := &entity.Blog{
		AuthorID:     authorID,
		Title:        req.Title,
		Slug:         req.Slug,
		Content:      req.Content,
		Excerpt:      req.Excerpt,
		ThumbnailURL: req.ThumbnailURL,
		Status:       entity.BlogStatusDraft,
		Visibility:   entity.BlogVisibilityPublic,
		PublishedAt:  req.PublishedAt,
	}

	if req.CategoryID != nil {
		if id, err := uuid.Parse(*req.CategoryID); err == nil {
			blog.CategoryID = &id
		}
	}

	tagIDs := make([]uuid.UUID, 0, len(req.TagIDs))
	for _, idStr := range req.TagIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			tagIDs = append(tagIDs, id)
		}
	}

	if err := uc.blogSvc.Create(ctx, blog, tagIDs); err != nil {
		return nil, err
	}

	// Fetch fresh to get relationships
	return uc.GetByID(ctx, blog.ID, &authorID)
}

func (uc *blogUseCase) GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := uc.blogSvc.GetByID(ctx, id, viewerID)
	if err != nil {
		return nil, err
	}
	return uc.toBlogResponse(blog), nil
}

func (uc *blogUseCase) GetBySlug(ctx context.Context, authorID uuid.UUID, slug string, viewerID *uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := uc.blogSvc.GetBySlug(ctx, authorID, slug, viewerID)
	if err != nil {
		return nil, err
	}
	return uc.toBlogResponse(blog), nil
}

func (uc *blogUseCase) List(ctx context.Context, params *dto.BlogFilterParams, viewerID *uuid.UUID) (*repository.PaginatedResult[dto.BlogListResponse], error) {
	filter := repository.BlogFilter{}
	if params.AuthorID != nil {
		if id, err := uuid.Parse(*params.AuthorID); err == nil {
			filter.AuthorID = &id
		}
	}
	if params.CategoryID != nil {
		if id, err := uuid.Parse(*params.CategoryID); err == nil {
			filter.CategoryID = &id
		}
	}
	if params.Status != nil {
		status := entity.BlogStatus(*params.Status)
		filter.Status = &status
	}
	if params.Visibility != nil {
		visibility := entity.BlogVisibility(*params.Visibility)
		filter.Visibility = &visibility
	}
	if params.Search != nil {
		filter.Search = params.Search
	}

	// Filter scheduled posts if not viewing own posts
	shouldFilterScheduled := true
	if viewerID != nil && filter.AuthorID != nil && *filter.AuthorID == *viewerID {
		shouldFilterScheduled = false
	}

	if shouldFilterScheduled {
		now := time.Now()
		filter.PublishedBefore = &now
	}

	result, err := uc.blogSvc.List(ctx, filter, repository.Pagination{Page: params.Page, PageSize: params.PageSize}, viewerID)
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

func (uc *blogUseCase) Update(ctx context.Context, id uuid.UUID, authorID uuid.UUID, req *dto.UpdateBlogRequest) (*dto.BlogResponse, error) {
	// First get the blog to ensure existence and ownership logic is handled by service,
	// but service Update expects a populated blog object.
	// Actually, service.Update logic expects us to pass the updated entity info.
	// But to construct updated entity safely, we need current state.

	// Better: UseCase fetches current, updates fields, calls Service Update.
	// Service Update will re-check ownership?
	// My service Update implementation DOES check ownership if we pass the original ID?
	// Wait, my Service Update:
	// func (s *blogService) Update(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error
	// It finds existing by Slug, but DOES NOT check author ownership inside Update explicitly?
	// It checks `if existing != nil && existing.ID != blog.ID`

	// Ah, I need to fetch it first to verify author ID matches.
	// Service Update takes *entity.Blog. It assumes *entity.Blog has correct ID.
	// Does it check if blog.AuthorID matches the one in DB?
	// The Service implementation of Update calls `s.blogRepo.Update(ctx, blog)`.
	// The REPO update might inherently be unsafe if we change AuthorID.

	// Safest: UseCase fetches, verifies author, then updates fields.

	// Wait, Service Update implementation I wrote:
	// `existing, _ := s.blogRepo.FindBySlug(...)` checks slug collision.
	// It doesn't check AuthorID matches caller.

	// So UseCase MUST fetch first.
	// However, `uc.blogSvc.GetByID` checks ACCESS, which might fail if draft and wrong viewer.
	// But here `authorID` is viewerID.

	// I'll call `uc.blogSvc.GetByID`.

	blog, err := uc.blogSvc.GetByID(ctx, id, &authorID)
	if err != nil {
		return nil, err
	}

	// Verify ownership explicitly just in case GetByID allows read but not write?
	// (CheckAccess logic: Author can always access. So if err==nil, we are good to read).
	// But to Update, we must be Author.
	if blog.AuthorID != authorID {
		return nil, ErrBlogAccessDenied
	}

	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Slug != nil {
		blog.Slug = *req.Slug
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.Excerpt != nil {
		blog.Excerpt = req.Excerpt
	}
	if req.ThumbnailURL != nil {
		blog.ThumbnailURL = req.ThumbnailURL
	}
	if req.PublishedAt != nil {
		blog.PublishedAt = req.PublishedAt
	}
	if req.CategoryID != nil {
		if id, err := uuid.Parse(*req.CategoryID); err == nil {
			blog.CategoryID = &id
		}
	}

	var tagIDs []uuid.UUID
	if req.TagIDs != nil {
		tagIDs = make([]uuid.UUID, 0, len(req.TagIDs))
		for _, idStr := range req.TagIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				tagIDs = append(tagIDs, id)
			}
		}
	}

	if err := uc.blogSvc.Update(ctx, blog, tagIDs); err != nil {
		return nil, err
	}

	// Fetch fresh
	return uc.GetByID(ctx, id, &authorID)
}

func (uc *blogUseCase) Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error {
	return uc.blogSvc.Delete(ctx, id, authorID)
}

func (uc *blogUseCase) Publish(ctx context.Context, id uuid.UUID, authorID uuid.UUID, req *dto.PublishBlogRequest) (*dto.BlogResponse, error) {
	blog, err := uc.blogSvc.Publish(ctx, id, authorID, entity.BlogVisibility(req.Visibility), req.PublishedAt)
	if err != nil {
		return nil, err
	}
	return uc.toBlogResponse(blog), nil
}

func (uc *blogUseCase) Unpublish(ctx context.Context, id uuid.UUID, authorID uuid.UUID) (*dto.BlogResponse, error) {
	blog, err := uc.blogSvc.Unpublish(ctx, id, authorID)
	if err != nil {
		return nil, err
	}
	return uc.toBlogResponse(blog), nil
}

func (uc *blogUseCase) React(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.ReactionRequest) (*dto.ReactionResponse, error) {
	upvotes, downvotes, err := uc.blogSvc.React(ctx, id, userID, entity.ReactionType(req.Reaction))
	if err != nil {
		return nil, err
	}

	return &dto.ReactionResponse{
		BlogID:        id,
		UpvoteCount:   upvotes,
		DownvoteCount: downvotes,
		UserReaction:  entity.ReactionType(req.Reaction),
	}, nil
}

func (uc *blogUseCase) toBlogResponse(blog *entity.Blog) *dto.BlogResponse {
	resp := &dto.BlogResponse{
		ID:            blog.ID,
		AuthorID:      blog.AuthorID,
		CategoryID:    blog.CategoryID,
		Title:         blog.Title,
		Slug:          blog.Slug,
		Excerpt:       blog.Excerpt,
		Content:       blog.Content,
		ThumbnailURL:  blog.ThumbnailURL,
		Status:        blog.Status,
		Visibility:    blog.Visibility,
		PublishedAt:   blog.PublishedAt,
		Tags:          make([]dto.TagResponse, 0),
		UpvoteCount:   blog.UpvoteCount,
		DownvoteCount: blog.DownvoteCount,
		CreatedAt:     blog.CreatedAt,
		UpdatedAt:     blog.UpdatedAt,
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

func (uc *blogUseCase) toBlogListResponse(blog *entity.Blog) dto.BlogListResponse {
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
