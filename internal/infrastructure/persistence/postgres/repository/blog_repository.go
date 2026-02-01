package repository

import (
	"context"
	"math"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type blogRepository struct {
	db *gorm.DB
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *gorm.DB) repository.BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(ctx context.Context, blog *entity.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *blogRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Blog, error) {
	var blog entity.Blog
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&blog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &blog, err
}

func (r *blogRepository) FindBySlug(ctx context.Context, authorID uuid.UUID, slug string) (*entity.Blog, error) {
	var blog entity.Blog
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("author_id = ? AND slug = ? AND deleted_at IS NULL", authorID, slug).
		First(&blog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &blog, err
}

func (r *blogRepository) FindAll(ctx context.Context, filter repository.BlogFilter, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error) {
	var blogs []entity.Blog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Blog{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Visibility != nil {
		query = query.Where("visibility = ?", *filter.Visibility)
	}
	if filter.Search != nil && *filter.Search != "" {
		searchPattern := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern)
	}
	if len(filter.TagIDs) > 0 {
		query = query.Joins("JOIN blog_tags ON blog_tags.blog_id = blogs.id").
			Where("blog_tags.tag_id IN ?", filter.TagIDs)
	}
	if filter.PublishedBefore != nil {
		query = query.Where("published_at <= ?", filter.PublishedBefore)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination and fetch
	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Blog]{
		Data:       blogs,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *blogRepository) Update(ctx context.Context, blog *entity.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *blogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Blog{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *blogRepository) AddTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		err := r.db.WithContext(ctx).Exec(
			"INSERT INTO blog_tags (blog_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			blogID, tagID,
		).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *blogRepository) RemoveTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Exec("DELETE FROM blog_tags WHERE blog_id = ? AND tag_id IN ?", blogID, tagIDs).
		Error
}

func (r *blogRepository) ReplaceTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	// Remove all existing tags
	if err := r.db.WithContext(ctx).Exec("DELETE FROM blog_tags WHERE blog_id = ?", blogID).Error; err != nil {
		return err
	}
	// Add new tags
	return r.AddTags(ctx, blogID, tagIDs)
}

func (r *blogRepository) React(ctx context.Context, blogID, userID uuid.UUID, reactionType entity.ReactionType) (int, int, error) {
	var upDelta, downDelta int

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Get existing reaction
		var existingReaction entity.BlogReaction
		err := tx.Where("blog_id = ? AND user_id = ?", blogID, userID).First(&existingReaction).Error

		isNewReaction := err == gorm.ErrRecordNotFound
		if err != nil && !isNewReaction {
			return err
		}

		// 2. Determine action
		if isNewReaction {
			if reactionType == entity.ReactionTypeNone {
				// Nothing to do
			} else {
				// Create new reaction
				newReaction := entity.BlogReaction{
					BlogID: blogID,
					UserID: userID,
					Type:   reactionType,
				}
				if err := tx.Create(&newReaction).Error; err != nil {
					return err
				}

				// Calculate delta
				if reactionType == entity.ReactionTypeDownvote {
					downDelta = 1
				} else {
					upDelta = 1
				}
			}
		} else {
			// Reaction exists
			if reactionType == entity.ReactionTypeNone {
				// Remove reaction
				if err := tx.Delete(&existingReaction).Error; err != nil {
					return err
				}
				// Calculate delta
				if existingReaction.Type == entity.ReactionTypeDownvote {
					downDelta = -1
				} else {
					upDelta = -1
				}
			} else if reactionType != existingReaction.Type {
				// Change reaction (swap)
				existingReaction.Type = reactionType
				if err := tx.Save(&existingReaction).Error; err != nil {
					return err
				}

				// Calculate delta
				if reactionType == entity.ReactionTypeUpvote {
					upDelta = 1
					downDelta = -1
				} else {
					downDelta = 1
					upDelta = -1
				}
			}
			// If reactionType == existingReaction.Type, do nothing (deltas stay 0)
		}

		return nil
	})

	return upDelta, downDelta, err
}

func (r *blogRepository) UpdateCounts(ctx context.Context, blogID uuid.UUID, upDelta, downDelta int) error {
	updates := map[string]interface{}{}
	if upDelta != 0 {
		updates["upvote_count"] = gorm.Expr("upvote_count + ?", upDelta)
	}
	if downDelta != 0 {
		updates["downvote_count"] = gorm.Expr("downvote_count + ?", downDelta)
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Model(&entity.Blog{}).Where("id = ?", blogID).Updates(updates).Error
}

func (r *blogRepository) FindRelated(ctx context.Context, blogID uuid.UUID, limit int) ([]entity.Blog, error) {
	var blogs []entity.Blog

	// Subquery to find tags of the current blog
	subQuery := r.db.Table("blog_tags").Select("tag_id").Where("blog_id = ?", blogID)

	err := r.db.WithContext(ctx).
		Model(&entity.Blog{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Joins("JOIN blog_tags ON blog_tags.blog_id = blogs.id").
		Where("blog_tags.tag_id IN (?)", subQuery).
		Where("blogs.id != ?", blogID).
		Where("blogs.status = ?", entity.BlogStatusPublished).
		Group("blogs.id").
		Order("COUNT(blog_tags.tag_id) DESC, blogs.published_at DESC").
		Limit(limit).
		Find(&blogs).Error

	return blogs, err
}
