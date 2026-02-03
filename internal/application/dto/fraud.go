package dto

import (
	"time"

	"github.com/google/uuid"
)

// RiskScoreResponse represents a user's risk assessment
type RiskScoreResponse struct {
	ID                        uuid.UUID `json:"id"`
	UserID                    uuid.UUID `json:"user_id"`
	OverallScore              int       `json:"overall_score"`               // 0-100
	FollowerAuthenticityScore int       `json:"follower_authenticity_score"` // 0-100
	EngagementQualityScore    int       `json:"engagement_quality_score"`    // 0-100
	AccountAgeFactor          float64   `json:"account_age_factor"`          // 0-1
	BadgeStatus               string    `json:"badge_status"`                // "eligible", "active", "revoked", "none"
	LastCalculatedAt          time.Time `json:"last_calculated_at"`
}

// FraudDashboardRequest represents filters for the admin dashboard
type FraudDashboardRequest struct {
	MinRiskScore int        `form:"min_risk_score" binding:"min=0,max=100"`                                  // Filter by minimum risk score
	MaxRiskScore int        `form:"max_risk_score" binding:"min=0,max=100"`                                  // Filter by maximum risk score
	SignalTypes  []string   `form:"signal_types"`                                                            // Filter by specific signal types
	ReviewStatus string     `form:"review_status" binding:"omitempty,oneof=pending reviewed banned cleared"` // Filter by review status
	FromDate     *time.Time `form:"from_date" time_format:"2006-01-02"`                                      // Filter from date
	ToDate       *time.Time `form:"to_date" time_format:"2006-01-02"`                                        // Filter to date
	Page         int        `form:"page,default=1" binding:"min=1"`
	PageSize     int        `form:"page_size,default=20" binding:"min=1,max=100"`
}

// FraudDashboardUser represents a user entry in the admin dashboard
type FraudDashboardUser struct {
	UserID                uuid.UUID          `json:"user_id"`
	Username              string             `json:"username"`
	Email                 string             `json:"email"`
	OverallScore          int                `json:"overall_score"`
	FollowerCount         int                `json:"follower_count"`
	BotFollowerCount      int                `json:"bot_follower_count"`
	ActiveSignals         []BotSignalSummary `json:"active_signals"`
	LastReviewAction      string             `json:"last_review_action"` // "reviewed", "banned", "warned", "cleared", ""
	LastReviewedAt        *time.Time         `json:"last_reviewed_at"`
	RiskScoreCalculatedAt time.Time          `json:"risk_score_calculated_at"`
}

// BotSignalSummary represents a summary of a bot detection signal
type BotSignalSummary struct {
	SignalType      string    `json:"signal_type"`
	ConfidenceScore float64   `json:"confidence_score"`
	DetectedAt      time.Time `json:"detected_at"`
}

// FraudDashboardResponse represents the admin dashboard data
type FraudDashboardResponse struct {
	Users      []FraudDashboardUser `json:"users"`
	TotalCount int                  `json:"total_count"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// ReviewUserRequest represents a request to mark a user as reviewed
type ReviewUserRequest struct {
	Notes string `json:"notes" binding:"max=1000"`
}

// ReviewUserResponse represents the result of reviewing a user
type ReviewUserResponse struct {
	ReviewID   uuid.UUID `json:"review_id"`
	UserID     uuid.UUID `json:"user_id"`
	AdminID    uuid.UUID `json:"admin_id"`
	Action     string    `json:"action"` // "reviewed"
	Notes      string    `json:"notes"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

// BanUserRequest represents a request to ban a user
type BanUserRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
	Notes  string `json:"notes" binding:"max=1000"`
}

// BanUserResponse represents the result of banning a user
type BanUserResponse struct {
	ReviewID uuid.UUID `json:"review_id"`
	UserID   uuid.UUID `json:"user_id"`
	AdminID  uuid.UUID `json:"admin_id"`
	Action   string    `json:"action"` // "banned"
	Reason   string    `json:"reason"`
	Notes    string    `json:"notes"`
	BannedAt time.Time `json:"banned_at"`
}

// FraudTrendsRequest represents filters for fraud analytics
type FraudTrendsRequest struct {
	Period   string     `form:"period,default=7d" binding:"oneof=24h 7d 30d 90d"` // Time period
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02"`
}

// FraudTrendsResponse represents fraud analytics data
type FraudTrendsResponse struct {
	Period                string           `json:"period"`
	FromDate              time.Time        `json:"from_date"`
	ToDate                time.Time        `json:"to_date"`
	TotalBotSignals       int              `json:"total_bot_signals"`
	SignalsByType         map[string]int   `json:"signals_by_type"`
	NewSuspiciousAccounts int              `json:"new_suspicious_accounts"`
	BannedAccounts        int              `json:"banned_accounts"`
	ReviewedAccounts      int              `json:"reviewed_accounts"`
	AverageRiskScore      float64          `json:"average_risk_score"`
	RiskScoreDistribution map[string]int   `json:"risk_score_distribution"` // "0-20", "21-40", etc.
	DailyStats            []DailyFraudStat `json:"daily_stats"`
}

// DailyFraudStat represents fraud statistics for a single day
type DailyFraudStat struct {
	Date                  string `json:"date"`
	NewSignals            int    `json:"new_signals"`
	NewSuspiciousAccounts int    `json:"new_suspicious_accounts"`
	BannedAccounts        int    `json:"banned_accounts"`
}

// BatchAnalyzeRequest represents a request to trigger batch analysis
type BatchAnalyzeRequest struct {
	DateFrom *time.Time `json:"date_from" time_format:"2006-01-02"` // Analyze from date (default: last run)
	DateTo   *time.Time `json:"date_to" time_format:"2006-01-02"`   // Analyze to date (default: now)
}

// BatchAnalyzeResponse represents the result of batch analysis
type BatchAnalyzeResponse struct {
	JobID              uuid.UUID  `json:"job_id"`
	Status             string     `json:"status"` // "started", "completed", "failed"
	StartedAt          time.Time  `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	ProcessedFollowers int        `json:"processed_followers"`
	NewSignalsDetected int        `json:"new_signals_detected"`
	UsersScored        int        `json:"users_scored"`
	Message            string     `json:"message"`
}

// UserBadgeResponse represents a user's badge status
type UserBadgeResponse struct {
	UserID        uuid.UUID  `json:"user_id"`
	BadgeType     string     `json:"badge_type"`
	Status        string     `json:"status"`
	EligibleSince *time.Time `json:"eligible_since"`
	ActivatedAt   *time.Time `json:"activated_at"`
}

// BotFollowerNotificationResponse represents a notification about flagged followers
type BotFollowerNotificationResponse struct {
	ID              uuid.UUID  `json:"id"`
	BotFollowerID   uuid.UUID  `json:"bot_follower_id"`
	BotFollowerName string     `json:"bot_follower_name"`
	SignalType      string     `json:"signal_type"`
	ConfidenceScore float64    `json:"confidence_score"`
	SentAt          time.Time  `json:"sent_at"`
	ReadAt          *time.Time `json:"read_at"`
}
