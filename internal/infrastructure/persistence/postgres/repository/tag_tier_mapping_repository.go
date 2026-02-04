package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tagTierMappingRepository struct {
	db *gorm.DB
}

// NewTagTierMappingRepository creates a new tag tier mapping repository
func NewTagTierMappingRepository(db *gorm.DB) repository.TagTierMappingRepository {
	return &tagTierMappingRepository{db: db}
}

// Create creates a new tag-tier mapping
func (r *tagTierMappingRepository) Create(ctx context.Context, mapping *entity.TagTierMapping) error {
	return r.db.WithContext(ctx).Create(mapping).Error
}

// FindByID retrieves a mapping by ID
func (r *tagTierMappingRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.TagTierMapping, error) {
	var mapping entity.TagTierMapping
	err := r.db.WithContext(ctx).First(&mapping, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// FindByAuthorAndTag retrieves a mapping by author ID and tag ID
func (r *tagTierMappingRepository) FindByAuthorAndTag(ctx context.Context, authorID, tagID uuid.UUID) (*entity.TagTierMapping, error) {
	var mapping entity.TagTierMapping
	err := r.db.WithContext(ctx).
		Where("author_id = ? AND tag_id = ?", authorID, tagID).
		First(&mapping).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// FindByAuthor retrieves all mappings for an author
func (r *tagTierMappingRepository) FindByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.TagTierMapping, error) {
	var mappings []entity.TagTierMapping
	err := r.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Find(&mappings).Error
	return mappings, err
}

// FindByTagIDs retrieves mappings for specific tags and author
func (r *tagTierMappingRepository) FindByTagIDs(ctx context.Context, authorID uuid.UUID, tagIDs []uuid.UUID) ([]entity.TagTierMapping, error) {
	var mappings []entity.TagTierMapping
	query := r.db.WithContext(ctx).Where("author_id = ?", authorID)

	if len(tagIDs) > 0 {
		query = query.Where("tag_id IN ?", tagIDs)
	}

	err := query.Find(&mappings).Error
	return mappings, err
}

// Update updates an existing mapping
func (r *tagTierMappingRepository) Update(ctx context.Context, mapping *entity.TagTierMapping) error {
	return r.db.WithContext(ctx).Save(mapping).Error
}

// Upsert creates or updates a mapping based on (author_id, tag_id) uniqueness
func (r *tagTierMappingRepository) Upsert(ctx context.Context, mapping *entity.TagTierMapping) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "author_id"}, {Name: "tag_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"required_tier", "updated_at"}),
		}).
		Create(mapping).Error
}

// Delete deletes a tag-tier mapping
func (r *tagTierMappingRepository) Delete(ctx context.Context, authorID, tagID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("author_id = ? AND tag_id = ?", authorID, tagID).
		Delete(&entity.TagTierMapping{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByID deletes a mapping by ID
func (r *tagTierMappingRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.TagTierMapping{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CountBlogsByTagAndAuthor counts blogs that have the specified tag from the author
func (r *tagTierMappingRepository) CountBlogsByTagAndAuthor(ctx context.Context, authorID, tagID uuid.UUID) (int64, error) {
	var count int64
	// Count blogs from the author that have the specified tag
	err := r.db.WithContext(ctx).
		Model(&entity.Blog{}).
		Joins("JOIN blog_tags ON blogs.id = blog_tags.blog_id").
		Where("blogs.author_id = ? AND blog_tags.tag_id = ?", authorID, tagID).
		Count(&count).Error
	return count, err
}

// WithTx returns a new repository with the given transaction
func (r *tagTierMappingRepository) WithTx(tx interface{}) repository.TagTierMappingRepository {
	if gormDB, ok := tx.(*gorm.DB); ok {
		return &tagTierMappingRepository{db: gormDB}
	}
	return r
}
