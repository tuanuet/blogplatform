package repository

import (
	"context"
	"math"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *entity.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *subscriptionRepository) Delete(ctx context.Context, subscriberID, authorID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("subscriber_id = ? AND author_id = ?", subscriberID, authorID).
		Delete(&entity.Subscription{}).Error
}

func (r *subscriptionRepository) Exists(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscriber_id = ? AND author_id = ?", subscriberID, authorID).
		Count(&count).Error
	return count > 0, err
}

func (r *subscriptionRepository) FindBySubscriber(ctx context.Context, subscriberID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Subscription], error) {
	var subscriptions []entity.Subscription
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscriber_id = ?", subscriberID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Preload("Author").
		Order("created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Subscription]{
		Data:       subscriptions,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *subscriptionRepository) FindByAuthor(ctx context.Context, authorID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Subscription], error) {
	var subscriptions []entity.Subscription
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("author_id = ?", authorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Preload("Subscriber").
		Order("created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Subscription]{
		Data:       subscriptions,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *subscriptionRepository) FindTopSubscribedAuthors(ctx context.Context, pagination repository.Pagination) (*repository.PaginatedResult[entity.TopSubscribedAuthor], error) {
	countQuery := `SELECT COUNT(*) FROM (
		SELECT s.author_id
		FROM subscriptions s
		JOIN users u ON u.id = s.author_id
		WHERE u.is_active = TRUE
		  AND u.deleted_at IS NULL
		GROUP BY s.author_id
	) ranked_authors`

	var total int64
	if err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize

	type topSubscribedAuthorRow struct {
		AuthorID        uuid.UUID `gorm:"column:author_id"`
		Username        string    `gorm:"column:username"`
		DisplayName     string    `gorm:"column:display_name"`
		AvatarURL       string    `gorm:"column:avatar_url"`
		SubscriberCount int64     `gorm:"column:subscriber_count"`
	}

	query := `SELECT
		s.author_id,
		u.name AS username,
		COALESCE(u.display_name, u.name) AS display_name,
		COALESCE(u.avatar_url, '') AS avatar_url,
		COUNT(*) AS subscriber_count
	FROM subscriptions s
	JOIN users u ON u.id = s.author_id
	WHERE u.is_active = TRUE
	  AND u.deleted_at IS NULL
	GROUP BY s.author_id, u.name, u.display_name, u.avatar_url
	ORDER BY subscriber_count DESC, s.author_id ASC
	LIMIT ? OFFSET ?`

	var rows []topSubscribedAuthorRow
	if err := r.db.WithContext(ctx).Raw(query, pagination.PageSize, offset).Scan(&rows).Error; err != nil {
		return nil, err
	}

	authors := make([]entity.TopSubscribedAuthor, 0, len(rows))
	for _, row := range rows {
		authors = append(authors, entity.TopSubscribedAuthor{
			AuthorID:        row.AuthorID,
			Username:        row.Username,
			DisplayName:     row.DisplayName,
			AvatarURL:       row.AvatarURL,
			SubscriberCount: row.SubscriberCount,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.TopSubscribedAuthor]{
		Data:       authors,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *subscriptionRepository) CountSubscribers(ctx context.Context, authorID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("author_id = ?", authorID).
		Count(&count).Error
	return count, err
}

func (r *subscriptionRepository) CountBySubscriber(ctx context.Context, subscriberID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscriber_id = ?", subscriberID).
		Count(&count).Error
	return count, err
}

func (r *subscriptionRepository) UpdateExpiry(ctx context.Context, userID, authorID uuid.UUID, expiresAt time.Time, tier string) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscriber_id = ? AND author_id = ?", userID, authorID).
		Updates(map[string]interface{}{
			"expires_at": expiresAt,
			"tier":       tier,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *subscriptionRepository) FindActiveSubscription(ctx context.Context, userID, authorID uuid.UUID) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.WithContext(ctx).
		Where("subscriber_id = ? AND author_id = ?", userID, authorID).
		Where("expires_at > ? OR expires_at IS NULL", time.Now()).
		First(&subscription).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &subscription, nil
}

// WithTx returns a new repository with the given transaction
func (r *subscriptionRepository) WithTx(tx interface{}) repository.SubscriptionRepository {
	if gormDB, ok := tx.(*gorm.DB); ok {
		return &subscriptionRepository{db: gormDB}
	}
	return r
}
