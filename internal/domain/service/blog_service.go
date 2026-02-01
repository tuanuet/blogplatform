package service

import (
	"context"
	"errors"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/aiagent/boilerplate/internal/infrastructure/cache"
	"github.com/google/uuid"
)

var (
	ErrBlogNotFound         = errors.New("blog not found")
	ErrBlogAccessDenied     = errors.New("access denied to this blog")
	ErrBlogAlreadyPublished = errors.New("blog is already published")
	ErrSlugAlreadyExists    = errors.New("slug already exists for this author")
)

type BlogService interface {
	Create(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*entity.Blog, error)
	GetBySlug(ctx context.Context, authorID uuid.UUID, slug string, viewerID *uuid.UUID) (*entity.Blog, error)
	List(ctx context.Context, filter repository.BlogFilter, pagination repository.Pagination, viewerID *uuid.UUID) (*repository.PaginatedResult[entity.Blog], error)
	Update(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, authorID uuid.UUID, visibility entity.BlogVisibility) (*entity.Blog, error)
	Unpublish(ctx context.Context, id uuid.UUID, authorID uuid.UUID) (*entity.Blog, error)
	React(ctx context.Context, id uuid.UUID, userID uuid.UUID, reactionType entity.ReactionType) (upvotes, downvotes int, err error)
	CheckAccess(ctx context.Context, blog *entity.Blog, viewerID *uuid.UUID) error
}

type blogService struct {
	blogRepo         repository.BlogRepository
	subscriptionRepo repository.SubscriptionRepository
	tagRepo          repository.TagRepository
	batcher          *ReactionBatcher
}

func NewBlogService(
	blogRepo repository.BlogRepository,
	subscriptionRepo repository.SubscriptionRepository,
	tagRepo repository.TagRepository,
	redis *cache.RedisClient,
) BlogService {
	// Initialize batcher with 5 second flush interval, now using Redis
	batcher := NewReactionBatcher(blogRepo, redis, 5*time.Second)
	batcher.Start()

	return &blogService{
		blogRepo:         blogRepo,
		subscriptionRepo: subscriptionRepo,
		tagRepo:          tagRepo,
		batcher:          batcher,
	}
}

func (s *blogService) Create(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error {
	existing, _ := s.blogRepo.FindBySlug(ctx, blog.AuthorID, blog.Slug)
	if existing != nil {
		return ErrSlugAlreadyExists
	}

	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return err
	}

	if len(tagIDs) > 0 {
		_ = s.blogRepo.AddTags(ctx, blog.ID, tagIDs)
	}

	return nil
}

func (s *blogService) GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}

	if err := s.CheckAccess(ctx, blog, viewerID); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogService) GetBySlug(ctx context.Context, authorID uuid.UUID, slug string, viewerID *uuid.UUID) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindBySlug(ctx, authorID, slug)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}

	if err := s.CheckAccess(ctx, blog, viewerID); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogService) List(ctx context.Context, filter repository.BlogFilter, pagination repository.Pagination, viewerID *uuid.UUID) (*repository.PaginatedResult[entity.Blog], error) {
	result, err := s.blogRepo.FindAll(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	// Filter by access
	accessibleBlogs := make([]entity.Blog, 0, len(result.Data))
	for _, blog := range result.Data {
		if s.CheckAccess(ctx, &blog, viewerID) == nil {
			accessibleBlogs = append(accessibleBlogs, blog)
		}
	}

	// Note: Total count might be inaccurate if filtered, but efficient count with permission is complex.
	// For now keeping simple: return filtered list.
	// Accurate count would requires repo support for permission check query.

	return &repository.PaginatedResult[entity.Blog]{
		Data:       accessibleBlogs,
		Total:      result.Total, // Should be count of accessible, strictly speaking.
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *blogService) Update(ctx context.Context, blog *entity.Blog, tagIDs []uuid.UUID) error {
	existing, _ := s.blogRepo.FindBySlug(ctx, blog.AuthorID, blog.Slug)
	if existing != nil && existing.ID != blog.ID {
		return ErrSlugAlreadyExists
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return err
	}

	if tagIDs != nil {
		_ = s.blogRepo.ReplaceTags(ctx, blog.ID, tagIDs)
	}
	return nil
}

func (s *blogService) Delete(ctx context.Context, id uuid.UUID, authorID uuid.UUID) error {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if blog == nil {
		return ErrBlogNotFound
	}

	if blog.AuthorID != authorID {
		return ErrBlogAccessDenied
	}

	return s.blogRepo.Delete(ctx, id)
}

func (s *blogService) Publish(ctx context.Context, id uuid.UUID, authorID uuid.UUID, visibility entity.BlogVisibility) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}

	if blog.AuthorID != authorID {
		return nil, ErrBlogAccessDenied
	}

	blog.Visibility = visibility
	blog.Publish()

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *blogService) Unpublish(ctx context.Context, id uuid.UUID, authorID uuid.UUID) (*entity.Blog, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blog == nil {
		return nil, ErrBlogNotFound
	}

	if blog.AuthorID != authorID {
		return nil, ErrBlogAccessDenied
	}

	blog.Unpublish()

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *blogService) React(ctx context.Context, id uuid.UUID, userID uuid.UUID, reactionType entity.ReactionType) (int, int, error) {
	// 1. Check if blog exists
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		return 0, 0, err
	}
	if blog == nil {
		return 0, 0, ErrBlogNotFound
	}

	// 2. Check access
	if err := s.CheckAccess(ctx, blog, &userID); err != nil {
		return 0, 0, err
	}

	// 3. Process reaction
	upDelta, downDelta, err := s.blogRepo.React(ctx, id, userID, reactionType)
	if err != nil {
		return 0, 0, err
	}

	// 4. Queue DB count update using Redis-based batcher
	s.batcher.Add(id, upDelta, downDelta)

	// 5. Return optimistic counts
	newUp := blog.UpvoteCount + upDelta
	newDown := blog.DownvoteCount + downDelta

	if newUp < 0 {
		newUp = 0
	}
	if newDown < 0 {
		newDown = 0
	}

	return newUp, newDown, nil
}

func (s *blogService) CheckAccess(ctx context.Context, blog *entity.Blog, viewerID *uuid.UUID) error {
	if blog.IsDraft() {
		if viewerID == nil || *viewerID != blog.AuthorID {
			return ErrBlogAccessDenied
		}
		return nil
	}

	if blog.IsPublic() {
		return nil
	}

	if blog.IsSubscribersOnly() {
		if viewerID != nil && *viewerID == blog.AuthorID {
			return nil
		}
		if viewerID != nil {
			isSubscribed, _ := s.subscriptionRepo.Exists(ctx, *viewerID, blog.AuthorID)
			if isSubscribed {
				return nil
			}
		}
		return ErrBlogAccessDenied
	}
	return nil
}
