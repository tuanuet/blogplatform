package service

import (
	"context"
	"time"

	"github.com/aiagent/pkg/logger"
)

// RankingJob handles background jobs for ranking calculations
type RankingJob struct {
	rankingSvc RankingService
}

// NewRankingJob creates a new ranking job
func NewRankingJob(rankingSvc RankingService) *RankingJob {
	return &RankingJob{
		rankingSvc: rankingSvc,
	}
}

// DailyRecalculation performs the daily ranking recalculation
// This should be called by a scheduler (e.g., cron job)
func (j *RankingJob) DailyRecalculation(ctx context.Context) error {
	logger.Info("Starting daily ranking recalculation")
	startTime := time.Now()

	// Step 1: Archive current rankings
	if err := j.rankingSvc.ArchiveRankings(ctx); err != nil {
		logger.Error("Failed to archive rankings", err)
		return err
	}
	logger.Info("Rankings archived successfully")

	// Step 2: Calculate all velocity scores
	if err := j.rankingSvc.CalculateAllVelocityScores(ctx); err != nil {
		logger.Error("Failed to calculate velocity scores", err)
		return err
	}
	logger.Info("Velocity scores calculated successfully")

	// Step 3: Assign rank positions
	if err := j.rankingSvc.AssignRankPositions(ctx); err != nil {
		logger.Error("Failed to assign rank positions", err)
		return err
	}
	logger.Info("Rank positions assigned successfully")

	duration := time.Since(startTime)
	logger.Info("Daily ranking recalculation completed", map[string]interface{}{"duration": duration})

	return nil
}

// StartDailyScheduler starts a background scheduler for daily recalculation at 0 AM (midnight)
// Note: In production, you might want to use a proper job scheduler like cron
func (j *RankingJob) StartDailyScheduler(ctx context.Context) {
	logger.Info("Starting daily ranking scheduler")

	// Calculate time until next 0 AM (midnight)
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	timeUntilNextRun := nextRun.Sub(now)
	logger.Info("Next ranking calculation scheduled", map[string]interface{}{
		"at": nextRun,
		"in": timeUntilNextRun,
	})

	// Wait until next run time
	time.AfterFunc(timeUntilNextRun, func() {
		j.runDaily(ctx)
	})
}

func (j *RankingJob) runDaily(ctx context.Context) {
	if err := j.DailyRecalculation(ctx); err != nil {
		logger.Error("Daily recalculation failed", err)
	}

	// Schedule next run (24 hours from now)
	time.AfterFunc(24*time.Hour, func() {
		j.runDaily(ctx)
	})
}
