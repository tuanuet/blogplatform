package repository

import (
	"context"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userFollowerSnapshotRepository struct {
	db *gorm.DB
}

// NewUserFollowerSnapshotRepository creates a new user follower snapshot repository
func NewUserFollowerSnapshotRepository(db *gorm.DB) repository.UserFollowerSnapshotRepository {
	return &userFollowerSnapshotRepository{db: db}
}

func (r *userFollowerSnapshotRepository) Create(ctx context.Context, snapshot *entity.UserFollowerSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

func (r *userFollowerSnapshotRepository) FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.UserFollowerSnapshot, error) {
	var snapshot entity.UserFollowerSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND snapshot_date = ?", userID, date).
		First(&snapshot).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &snapshot, err
}

func (r *userFollowerSnapshotRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserFollowerSnapshot, error) {
	var snapshot entity.UserFollowerSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("snapshot_date DESC").
		First(&snapshot).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &snapshot, err
}

func (r *userFollowerSnapshotRepository) FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.UserFollowerSnapshot, error) {
	var snapshots []entity.UserFollowerSnapshot
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND snapshot_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("snapshot_date DESC").
		Find(&snapshots).Error
	return snapshots, err
}

func (r *userFollowerSnapshotRepository) CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("author_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *userFollowerSnapshotRepository) CountBlogs(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Blog{}).
		Where("author_id = ? AND status = ? AND published_at BETWEEN ? AND ?",
			userID, entity.BlogStatusPublished, startDate, endDate).
		Count(&count).Error
	return count, err
}
