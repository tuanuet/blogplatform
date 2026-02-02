package modules

import (
	"context"

	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/infrastructure/config"
	"github.com/aiagent/pkg/logger"
	"go.uber.org/fx"
)

// SchedulerModule provides background job scheduling with lifecycle management
var SchedulerModule = fx.Module("scheduler",
	fx.Provide(newRankingJob),
	fx.Invoke(startScheduler),
)

// newRankingJob creates the ranking job instance
func newRankingJob(rankingSvc service.RankingService) *service.RankingJob {
	return service.NewRankingJob(rankingSvc)
}

// startScheduler starts the background scheduler with lifecycle hooks
func startScheduler(lc fx.Lifecycle, job *service.RankingJob, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !cfg.Scheduler.Enabled {
				logger.Info("Background scheduler is disabled")
				return nil
			}

			logger.Info("Starting background job scheduler", map[string]interface{}{
				"daily_recalculation_hour": cfg.Scheduler.DailyRecalculationHour,
				"timezone":                 cfg.Scheduler.Timezone,
			})

			// Start the daily ranking recalculation scheduler
			// This will run at the configured hour (default: 0 AM / midnight) every day
			go job.StartDailyScheduler(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !cfg.Scheduler.Enabled {
				return nil
			}

			logger.Info("Stopping background job scheduler")
			// Note: The scheduler uses time.AfterFunc which cannot be easily cancelled.
			// In production, consider using a proper cron library with context support.
			return nil
		},
	})
}
