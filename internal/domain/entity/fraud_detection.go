package entity

import (
	"time"

	"github.com/google/uuid"
)

// FollowerEvent tracks all follow actions with metadata for bot detection analysis
type FollowerEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_events_follower" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_events_following" json:"following_id"`
	Timestamp   time.Time `gorm:"not null;index:idx_follower_events_time" json:"timestamp"`
	IPAddress   string    `gorm:"type:varchar(45);index:idx_follower_events_ip" json:"ip_address"` // IPv6 max length
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	Referrer    string    `gorm:"type:varchar(500)" json:"referrer"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

func (FollowerEvent) TableName() string {
	return "follower_events"
}

// BotDetectionSignal stores detected bot patterns and suspicious activities
type BotDetectionSignal struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null;index:idx_bot_signals_user" json:"user_id"`                                             // The suspected bot account
	SignalType      string      `gorm:"type:varchar(50);not null;index:idx_bot_signals_type" json:"signal_type"`                                  // rapid_follows, ip_cluster, no_profile, suspicious_engagement
	ConfidenceScore float64     `gorm:"type:decimal(3,2);not null;check:confidence_score >= 0 and confidence_score <= 1" json:"confidence_score"` // 0.0 to 1.0
	DetectedAt      time.Time   `gorm:"not null;index:idx_bot_signals_time" json:"detected_at"`
	RelatedAccounts []uuid.UUID `gorm:"type:uuid[]" json:"related_accounts"` // Array of connected suspicious accounts
	Evidence        string      `gorm:"type:jsonb" json:"evidence"`          // JSON with detailed evidence
	Processed       bool        `gorm:"not null;default:false;index:idx_bot_signals_processed" json:"processed"`
	CreatedAt       time.Time   `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"not null;default:now()" json:"updated_at"`
}

func (BotDetectionSignal) TableName() string {
	return "bot_detection_signals"
}

// UserRiskScore caches calculated risk scores for users
type UserRiskScore struct {
	ID                        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID                    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_risk_user" json:"user_id"`
	OverallScore              int       `gorm:"not null;check:overall_score >= 0 and overall_score <= 100;index:idx_user_risk_score" json:"overall_score"`                 // 0-100, higher = more suspicious
	FollowerAuthenticityScore int       `gorm:"not null;check:follower_authenticity_score >= 0 and follower_authenticity_score <= 100" json:"follower_authenticity_score"` // 0-100
	EngagementQualityScore    int       `gorm:"not null;check:engagement_quality_score >= 0 and engagement_quality_score <= 100" json:"engagement_quality_score"`          // 0-100
	AccountAgeFactor          float64   `gorm:"type:decimal(3,2);not null;check:account_age_factor >= 0 and account_age_factor <= 1" json:"account_age_factor"`            // 0-1
	CalculationVersion        string    `gorm:"type:varchar(20);not null" json:"calculation_version"`                                                                      // Version of algorithm used
	LastCalculatedAt          time.Time `gorm:"not null;index:idx_user_risk_calculated" json:"last_calculated_at"`
	CreatedAt                 time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt                 time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (UserRiskScore) TableName() string {
	return "user_risk_scores"
}

// UserBadgeStatus tracks verified badge eligibility and status
type UserBadgeStatus struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_badge_status_user" json:"user_id"`
	BadgeType        string     `gorm:"type:varchar(50);not null" json:"badge_type"`                           // "verified", "authentic"
	Status           string     `gorm:"type:varchar(20);not null;index:idx_badge_status_status" json:"status"` // "eligible", "active", "revoked", "pending"
	EligibleSince    *time.Time `json:"eligible_since"`                                                        // When user first became eligible
	ActivatedAt      *time.Time `json:"activated_at"`                                                          // When badge was activated
	RevokedAt        *time.Time `json:"revoked_at"`                                                            // When badge was revoked
	RevocationReason string     `gorm:"type:varchar(255)" json:"revocation_reason"`
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

func (UserBadgeStatus) TableName() string {
	return "user_badge_status"
}

// AdminReview logs admin actions and reviews on suspicious accounts
type AdminReview struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AdminID           uuid.UUID `gorm:"type:uuid;not null;index:idx_admin_reviews_admin" json:"admin_id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index:idx_admin_reviews_user" json:"user_id"`         // The reviewed user
	Action            string    `gorm:"type:varchar(50);not null;index:idx_admin_reviews_action" json:"action"` // "reviewed", "banned", "warned", "cleared"
	RiskScoreAtReview int       `json:"risk_score_at_review"`                                                   // Risk score when action was taken
	Notes             string    `gorm:"type:text" json:"notes"`
	ReviewedAt        time.Time `gorm:"not null;index:idx_admin_reviews_time" json:"reviewed_at"`
	CreatedAt         time.Time `gorm:"not null;default:now()" json:"created_at"`
}

func (AdminReview) TableName() string {
	return "admin_reviews"
}

// BotFollowerNotification tracks notifications sent to users about flagged followers
type BotFollowerNotification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index:idx_bot_notif_user" json:"user_id"` // User who received notification
	BotFollowerID    uuid.UUID  `gorm:"type:uuid;not null" json:"bot_follower_id"`                  // The flagged bot account
	SignalID         uuid.UUID  `gorm:"type:uuid;not null" json:"signal_id"`                        // Reference to BotDetectionSignal
	NotificationType string     `gorm:"type:varchar(50);not null" json:"notification_type"`         // "email", "in_app", "both"
	SentAt           time.Time  `gorm:"not null" json:"sent_at"`
	ReadAt           *time.Time `json:"read_at"`
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"created_at"`
}

func (BotFollowerNotification) TableName() string {
	return "bot_follower_notifications"
}
