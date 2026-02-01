package repository

import (
	"context"
	"math"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new bookmark repository
func NewBookmarkRepository(db *gorm.DB) repository.BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) Add(ctx context.Context, userID, blogID uuid.UUID) error {
	// We use direct SQL execution or GORM association.
	// Using association might be safer if entities are setup correctly.
	// db.Model(&entity.User{ID: userID}).Association("BookmarkedBlogs").Append(&entity.Blog{ID: blogID})
	// But let's use the explicit many2many handling or just a raw insert if simple.

	// Since we defined the relationship in User entity, let's try to use it.
	// Note: Append might try to create the blog if not careful, but usually checking ID is enough.
	// Use Omit(".*") to avoid updating the blog itself?

	return r.db.WithContext(ctx).Model(&entity.User{ID: userID}).
		Association("BookmarkedBlogs").
		Append(&entity.Blog{ID: blogID})
}

func (r *bookmarkRepository) Remove(ctx context.Context, userID, blogID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.User{ID: userID}).
		Association("BookmarkedBlogs").
		Delete(&entity.Blog{ID: blogID})
}

func (r *bookmarkRepository) IsBookmarked(ctx context.Context, userID, blogID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("user_bookmarks").
		Where("user_id = ? AND blog_id = ?", userID, blogID).
		Count(&count).Error
	return count > 0, err
}

func (r *bookmarkRepository) CountByBlog(ctx context.Context, blogID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("user_bookmarks").
		Where("blog_id = ?", blogID).
		Count(&count).Error
	return count, err
}

func (r *bookmarkRepository) FindByUser(ctx context.Context, userID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error) {
	var blogs []entity.Blog
	var total int64

	// Base query joining blogs with user_bookmarks
	query := r.db.WithContext(ctx).Model(&entity.Blog{}).
		Joins("JOIN user_bookmarks ON user_bookmarks.blog_id = blogs.id").
		Where("user_bookmarks.user_id = ? AND blogs.deleted_at IS NULL", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize

	err := query.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("user_bookmarks.created_at DESC"). // Order by bookmark creation time usually
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&blogs).Error

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))
	if pagination.PageSize == 0 {
		totalPages = 0
	}

	return &repository.PaginatedResult[entity.Blog]{
		Data:       blogs,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}
