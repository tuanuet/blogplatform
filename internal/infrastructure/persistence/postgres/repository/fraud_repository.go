package repository

import (
	"context"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/service"
	"github.com/aiagent/boilerplate/internal/interfaces/http/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fraudDetectionRepository implements the FraudDetectionRepository interface using GORM
type fraudDetectionRepository struct {
	db *gorm.DB
}

// NewFraudDetectionRepository creates a new fraud detection repository instance
func NewFraudDetectionRepository(db *gorm.DB) service.FraudDetectionRepository {
	return &fraudDetectionRepository{
		db: db,
	}
}

// CreateFollowerEvent creates a new follower event record
func (r *fraudDetectionRepository) CreateFollowerEvent(ctx context.Context, event *entity.FollowerEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetFollowerEventsByUser retrieves follower events for a specific user
func (r *fraudDetectionRepository) GetFollowerEventsByUser(ctx context.Context, userID uuid.UUID, from, to *time.Time) ([]entity.FollowerEvent, error) {
	var events []entity.FollowerEvent
	query := r.db.WithContext(ctx).Where("following_id = ?", userID)

	if from != nil {
		query = query.Where("timestamp >= ?", *from)
	}
	if to != nil {
		query = query.Where("timestamp <= ?", *to)
	}

	err := query.Order("timestamp desc").Find(&events).Error
	return events, err
}

// GetFollowerEventsByIP retrieves follower events from a specific IP address
func (r *fraudDetectionRepository) GetFollowerEventsByIP(ctx context.Context, ipAddress string, from, to *time.Time) ([]entity.FollowerEvent, error) {
	var events []entity.FollowerEvent
	query := r.db.WithContext(ctx).Where("ip_address = ?", ipAddress)

	if from != nil {
		query = query.Where("timestamp >= ?", *from)
	}
	if to != nil {
		query = query.Where("timestamp <= ?", *to)
	}

	err := query.Order("timestamp desc").Find(&events).Error
	return events, err
}

// CreateBotSignal creates a new bot detection signal
func (r *fraudDetectionRepository) CreateBotSignal(ctx context.Context, signal *entity.BotDetectionSignal) error {
	return r.db.WithContext(ctx).Create(signal).Error
}

// GetBotSignalsByUser retrieves bot detection signals for a user
func (r *fraudDetectionRepository) GetBotSignalsByUser(ctx context.Context, userID uuid.UUID, processed bool) ([]entity.BotDetectionSignal, error) {
	var signals []entity.BotDetectionSignal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND processed = ?", userID, processed).
		Order("detected_at desc").
		Find(&signals).Error
	return signals, err
}

// GetUnprocessedBotSignals retrieves unprocessed bot signals with a limit
func (r *fraudDetectionRepository) GetUnprocessedBotSignals(ctx context.Context, limit int) ([]entity.BotDetectionSignal, error) {
	var signals []entity.BotDetectionSignal
	err := r.db.WithContext(ctx).
		Where("processed = ?", false).
		Order("detected_at asc").
		Limit(limit).
		Find(&signals).Error
	return signals, err
}

// MarkBotSignalAsProcessed marks a bot signal as processed
func (r *fraudDetectionRepository) MarkBotSignalAsProcessed(ctx context.Context, signalID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.BotDetectionSignal{}).
		Where("id = ?", signalID).
		Update("processed", true).
		Error
}

// CreateOrUpdateRiskScore creates or updates a user's risk score
func (r *fraudDetectionRepository) CreateOrUpdateRiskScore(ctx context.Context, score *entity.UserRiskScore) error {
	return r.db.WithContext(ctx).Save(score).Error
}

// GetRiskScoreByUser retrieves the risk score for a specific user
func (r *fraudDetectionRepository) GetRiskScoreByUser(ctx context.Context, userID uuid.UUID) (*entity.UserRiskScore, error) {
	var score entity.UserRiskScore
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&score).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &score, err
}

// GetUsersByRiskScoreRange retrieves users within a risk score range
func (r *fraudDetectionRepository) GetUsersByRiskScoreRange(ctx context.Context, minScore, maxScore int, page, pageSize int) ([]entity.UserRiskScore, int, error) {
	var scores []entity.UserRiskScore
	var totalCount int64

	query := r.db.WithContext(ctx).Where("overall_score >= ? AND overall_score <= ?", minScore, maxScore)

	// Get total count
	if err := query.Model(&entity.UserRiskScore{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Order("overall_score desc").
		Offset(offset).
		Limit(pageSize).
		Find(&scores).Error

	return scores, int(totalCount), err
}

// CreateOrUpdateBadgeStatus creates or updates a user's badge status
func (r *fraudDetectionRepository) CreateOrUpdateBadgeStatus(ctx context.Context, status *entity.UserBadgeStatus) error {
	return r.db.WithContext(ctx).Save(status).Error
}

// GetBadgeStatusByUser retrieves the badge status for a specific user
func (r *fraudDetectionRepository) GetBadgeStatusByUser(ctx context.Context, userID uuid.UUID) (*entity.UserBadgeStatus, error) {
	var status entity.UserBadgeStatus
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&status).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &status, err
}

// CreateAdminReview creates a new admin review record
func (r *fraudDetectionRepository) CreateAdminReview(ctx context.Context, review *entity.AdminReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// GetAdminReviewsByUser retrieves admin reviews for a specific user
func (r *fraudDetectionRepository) GetAdminReviewsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.AdminReview, error) {
	var reviews []entity.AdminReview
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("reviewed_at desc").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

// GetLastReviewByUser retrieves the most recent review for a user
func (r *fraudDetectionRepository) GetLastReviewByUser(ctx context.Context, userID uuid.UUID) (*entity.AdminReview, error) {
	var review entity.AdminReview
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("reviewed_at desc").
		First(&review).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &review, err
}

// CreateBotNotification creates a bot follower notification
func (r *fraudDetectionRepository) CreateBotNotification(ctx context.Context, notification *entity.BotFollowerNotification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetBotNotificationsByUser retrieves bot notifications for a user
func (r *fraudDetectionRepository) GetBotNotificationsByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]entity.BotFollowerNotification, error) {
	var notifications []entity.BotFollowerNotification
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("read_at IS NULL")
	}

	err := query.Order("sent_at desc").Find(&notifications).Error
	return notifications, err
}

// MarkNotificationAsRead marks a notification as read
func (r *fraudDetectionRepository) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.BotFollowerNotification{}).
		Where("id = ?", notificationID).
		Update("read_at", &now).
		Error
}

// GetSignalsCountByType retrieves count of signals by type within a date range
func (r *fraudDetectionRepository) GetSignalsCountByType(ctx context.Context, from, to time.Time) (map[string]int, error) {
	var results []struct {
		SignalType string
		Count      int
	}

	err := r.db.WithContext(ctx).
		Model(&entity.BotDetectionSignal{}).
		Select("signal_type, COUNT(*) as count").
		Where("detected_at >= ? AND detected_at <= ?", from, to).
		Group("signal_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	signalsByType := make(map[string]int)
	for _, result := range results {
		signalsByType[result.SignalType] = result.Count
	}

	return signalsByType, nil
}

// GetNewSuspiciousAccountsCount retrieves count of new suspicious accounts
func (r *fraudDetectionRepository) GetNewSuspiciousAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserRiskScore{}).
		Where("overall_score > ? AND created_at >= ? AND created_at <= ?", 50, from, to).
		Count(&count).Error
	return int(count), err
}

// GetBannedAccountsCount retrieves count of banned accounts
func (r *fraudDetectionRepository) GetBannedAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.AdminReview{}).
		Where("action = ? AND reviewed_at >= ? AND reviewed_at <= ?", "banned", from, to).
		Count(&count).Error
	return int(count), err
}

// GetReviewedAccountsCount retrieves count of reviewed accounts
func (r *fraudDetectionRepository) GetReviewedAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.AdminReview{}).
		Where("reviewed_at >= ? AND reviewed_at <= ?", from, to).
		Count(&count).Error
	return int(count), err
}

// GetAverageRiskScore retrieves the average risk score
func (r *fraudDetectionRepository) GetAverageRiskScore(ctx context.Context, from, to time.Time) (float64, error) {
	var avg float64
	err := r.db.WithContext(ctx).
		Model(&entity.UserRiskScore{}).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Select("COALESCE(AVG(overall_score), 0)").
		Scan(&avg).Error
	return avg, err
}

// GetRiskScoreDistribution retrieves the distribution of risk scores
func (r *fraudDetectionRepository) GetRiskScoreDistribution(ctx context.Context, from, to time.Time) (map[string]int, error) {
	var results []struct {
		Range string
		Count int
	}

	// Use CASE to create score ranges
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			CASE 
				WHEN overall_score BETWEEN 0 AND 20 THEN '0-20'
				WHEN overall_score BETWEEN 21 AND 40 THEN '21-40'
				WHEN overall_score BETWEEN 41 AND 60 THEN '41-60'
				WHEN overall_score BETWEEN 61 AND 80 THEN '61-80'
				ELSE '81-100'
			END as range,
			COUNT(*) as count
		FROM user_risk_scores
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY range
		ORDER BY range
	`, from, to).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int)
	for _, result := range results {
		distribution[result.Range] = result.Count
	}

	return distribution, nil
}

// GetDailyFraudStats retrieves daily fraud statistics
func (r *fraudDetectionRepository) GetDailyFraudStats(ctx context.Context, from, to time.Time) ([]dto.DailyFraudStat, error) {
	var results []dto.DailyFraudStat

	// Query for daily stats - this is a simplified version
	// In production, you might want to pre-aggregate this data
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			DATE(detected_at) as date,
			COUNT(*) as new_signals,
			COUNT(DISTINCT user_id) as new_suspicious_accounts
		FROM bot_detection_signals
		WHERE detected_at >= ? AND detected_at <= ?
		GROUP BY DATE(detected_at)
		ORDER BY date
	`, from, to).Scan(&results).Error

	return results, err
}
