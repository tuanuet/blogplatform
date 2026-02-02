package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

func (r *userRepository) GetInterests(ctx context.Context, userID uuid.UUID) ([]entity.Tag, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Interests").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return user.Interests, nil
}

func (r *userRepository) ReplaceInterests(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error {
	var user entity.User
	user.ID = userID

	var tags []entity.Tag
	if len(tagIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).Model(&user).Association("Interests").Replace(tags)
}

func (r *userRepository) CountByMonth(ctx context.Context, months int) ([]entity.MonthlyCount, error) {
	var results []entity.MonthlyCount
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at >= NOW() - make_interval(months => ?)", months).
		Group("month").
		Order("month DESC").
		Scan(&results).Error
	return results, err
}
