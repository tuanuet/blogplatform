package service

import (
	"context"
	"sync"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/valueobject"
	"github.com/aiagent/pkg/logger"
	"github.com/google/uuid"
)

// Constants for batch job processing
const (
	// Risk score thresholds
	lowRiskThreshold      = 20
	highRiskThreshold     = 70
	notificationThreshold = 50

	// Badge statuses
	badgeStatusEligible = "eligible"
	badgeStatusRevoked  = "revoked"
	badgeTypeVerified   = "verified"

	// Notification types
	notificationTypeInApp = "in_app"

	// Revocation reasons
	revocationReasonHighRisk = "Risk score exceeded threshold"

	// Batch processing
	defaultSignalsBatchSize = 1000
	defaultAnalysisDays     = 1
)

// batchJobService implements the BatchJobService interface
type batchJobService struct {
	repo      FraudDetectionRepository
	algorithm BotDetectionAlgorithm
	notifier  NotificationService
	jobs      map[uuid.UUID]*valueobject.BatchAnalyzeResult
	mu        sync.RWMutex
}

// NewBatchJobService creates a new batch job service instance
func NewBatchJobService(repo FraudDetectionRepository, algorithm BotDetectionAlgorithm, notifier NotificationService) BatchJobService {
	return &batchJobService{
		repo:      repo,
		algorithm: algorithm,
		notifier:  notifier,
		jobs:      make(map[uuid.UUID]*valueobject.BatchAnalyzeResult),
	}
}

// StartBatchAnalysis starts the batch analysis job
func (s *batchJobService) StartBatchAnalysis(ctx context.Context, dateFrom, dateTo *time.Time) (uuid.UUID, error) {
	jobID := uuid.New()

	// Determine date range
	now := time.Now()
	var from, to time.Time

	if dateFrom != nil {
		from = *dateFrom
	} else {
		// Default: analyze last N days
		from = now.AddDate(0, 0, -defaultAnalysisDays)
	}

	if dateTo != nil {
		to = *dateTo
	} else {
		to = now
	}

	// Initialize job status
	s.mu.Lock()
	s.jobs[jobID] = &valueobject.BatchAnalyzeResult{
		JobID:     jobID,
		Status:    "running", // "started" was defined as const but string is used in struct. Using "running" or const.
		StartedAt: time.Now(),
		Message:   "Analysis in progress",
	}
	s.mu.Unlock()

	// Run analysis in background (non-blocking)
	go s.runAnalysis(context.Background(), jobID, from, to)

	return jobID, nil
}

// runAnalysis performs the actual batch analysis
func (s *batchJobService) runAnalysis(ctx context.Context, jobID uuid.UUID, from, to time.Time) {
	// This is a simplified implementation
	// In production, you'd want to:
	// 1. Track job progress in a database
	// 2. Process in smaller batches to avoid memory issues
	// 3. Handle errors gracefully
	// 4. Send notifications when complete

	// Step 1: Get all follower events in the time range
	// (This is simplified - you'd need to implement GetFollowerEventsInRange in the repository)

	// Step 2: Analyze each unique follower
	processedFollowers := 0
	newSignals := 0
	usersScored := make(map[uuid.UUID]bool)

	// Get unprocessed bot signals to analyze
	signals, err := s.repo.GetUnprocessedBotSignals(ctx, defaultSignalsBatchSize)
	if err != nil {
		logger.Error("Error getting unprocessed signals", err, map[string]interface{}{"job_id": jobID})
		return
	}

	// Group signals by user and calculate risk scores
	userSignals := make(map[uuid.UUID][]entity.BotDetectionSignal)
	for _, signal := range signals {
		userSignals[signal.UserID] = append(userSignals[signal.UserID], signal)
		processedFollowers++
		newSignals++

		// Mark signal as processed
		if err := s.repo.MarkBotSignalAsProcessed(ctx, signal.ID); err != nil {
			logger.Error("Error marking signal as processed", err, map[string]interface{}{"job_id": jobID})
		}
	}

	// Calculate risk scores for each user
	for userID, signals := range userSignals {
		riskScore, err := s.algorithm.CalculateRiskScore(ctx, userID, signals, 0)
		if err != nil {
			logger.Error("Error calculating risk score", err, map[string]interface{}{"job_id": jobID, "user_id": userID})
			continue
		}

		// Save risk score
		if err := s.repo.CreateOrUpdateRiskScore(ctx, riskScore); err != nil {
			logger.Error("Error saving risk score", err, map[string]interface{}{"job_id": jobID, "user_id": userID})
			continue
		}

		usersScored[userID] = true

		// Update badge status based on risk score
		if err := s.updateBadgeStatus(ctx, userID, riskScore); err != nil {
			logger.Error("Error updating badge status", err, map[string]interface{}{"job_id": jobID, "user_id": userID})
		}

		// Send notifications for high-risk users
		if riskScore.OverallScore > notificationThreshold {
			if err := s.sendBotNotifications(ctx, userID, signals); err != nil {
				logger.Error("Error sending notifications", err, map[string]interface{}{"job_id": jobID, "user_id": userID})
			}
		}
	}

	// Detect coordinated bot networks
	networks, err := s.algorithm.DetectCoordinatedBots(ctx, signals)
	if err != nil {
		logger.Error("Error detecting coordinated bots", err, map[string]interface{}{"job_id": jobID})
	} else if len(networks) > 0 {
		logger.Info("Detected coordinated bot networks", map[string]interface{}{"job_id": jobID, "network_count": len(networks)})
	}

	logger.Info("Batch job completed", map[string]interface{}{
		"job_id":              jobID,
		"processed_followers": processedFollowers,
		"new_signals":         newSignals,
		"users_scored":        len(usersScored),
	})

	// Update job status
	s.mu.Lock()
	if job, exists := s.jobs[jobID]; exists {
		now := time.Now()
		job.Status = "completed"
		job.CompletedAt = &now
		job.ProcessedFollowers = processedFollowers
		job.NewSignalsDetected = newSignals
		job.UsersScored = len(usersScored)
		job.Message = "Analysis completed successfully"
	}
	s.mu.Unlock()
}

// updateBadgeStatus updates the badge status based on risk score
func (s *batchJobService) updateBadgeStatus(ctx context.Context, userID uuid.UUID, riskScore *entity.UserRiskScore) error {
	// Get existing badge status
	existingStatus, err := s.repo.GetBadgeStatusByUser(ctx, userID)
	if err != nil {
		return err
	}

	// Determine new status
	var newStatus string
	if riskScore.OverallScore < lowRiskThreshold {
		newStatus = badgeStatusEligible
	} else if riskScore.OverallScore > highRiskThreshold && existingStatus != nil && existingStatus.Status == "active" {
		newStatus = badgeStatusRevoked
	} else {
		return nil // No change needed
	}

	now := time.Now()
	status := &entity.UserBadgeStatus{
		ID:        uuid.New(),
		UserID:    userID,
		BadgeType: badgeTypeVerified,
		Status:    newStatus,
		UpdatedAt: now,
	}

	if existingStatus != nil {
		status.ID = existingStatus.ID
		if existingStatus.EligibleSince != nil {
			status.EligibleSince = existingStatus.EligibleSince
		}
		if existingStatus.ActivatedAt != nil {
			status.ActivatedAt = existingStatus.ActivatedAt
		}
	}

	if newStatus == badgeStatusEligible && status.EligibleSince == nil {
		status.EligibleSince = &now
	}

	if newStatus == badgeStatusRevoked {
		status.RevokedAt = &now
		status.RevocationReason = revocationReasonHighRisk
	}

	return s.repo.CreateOrUpdateBadgeStatus(ctx, status)
}

// sendBotNotifications sends notifications to users about flagged bot followers
func (s *batchJobService) sendBotNotifications(ctx context.Context, userID uuid.UUID, signals []entity.BotDetectionSignal) error {
	// Create notification records
	notifications := make([]valueobject.BotFollowerNotificationResult, 0, len(signals))

	for _, signal := range signals {
		notification := &entity.BotFollowerNotification{
			ID:               uuid.New(),
			UserID:           userID,
			BotFollowerID:    signal.UserID,
			SignalID:         signal.ID,
			NotificationType: notificationTypeInApp,
			SentAt:           time.Now(),
		}

		if err := s.repo.CreateBotNotification(ctx, notification); err != nil {
			return err
		}

		notifications = append(notifications, valueobject.BotFollowerNotificationResult{
			ID:              notification.ID,
			BotFollowerID:   signal.UserID,
			SignalType:      signal.SignalType,
			ConfidenceScore: signal.ConfidenceScore,
			SentAt:          notification.SentAt,
		})
	}

	// Send actual notification
	if len(notifications) > 0 {
		return s.notifier.SendBotFollowerNotification(ctx, userID, notifications)
	}

	return nil
}

// GetBatchJobStatus retrieves the status of a batch job
func (s *batchJobService) GetBatchJobStatus(ctx context.Context, jobID uuid.UUID) (*valueobject.BatchAnalyzeResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return &valueobject.BatchAnalyzeResult{
			JobID:   jobID,
			Status:  "unknown",
			Message: "Job not found",
		}, nil
	}

	// Return a copy to avoid race conditions if caller modifies it (though returning pointer... technically safe if we don't modify the struct fields in place after returning, but struct has value semantics mostly)
	// Be safe and return a copy or new struct
	result := *job
	return &result, nil
}
