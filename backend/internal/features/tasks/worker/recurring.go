package tasks_worker

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	"go.uber.org/zap"
)

type TasksRepository interface {
	ProcessDueRecurringTasks(
		ctx context.Context,
		now time.Time,
		calculateNextRun func(cronExpr string, from time.Time) (time.Time, error),
	) error
}

func StartRecurringTasksWorker(ctx context.Context, repo TasksRepository, interval time.Duration) {
	log := core_logger.FromContext(ctx)
	ticker := time.NewTicker(interval)
	
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	
	calculateNextRun := func(cronExpr string, from time.Time) (time.Time, error) {
		schedule, err := parser.Parse(cronExpr)
		if err != nil {
			return time.Time{}, err
		}
		return schedule.Next(from), nil
	}

	go func() {
		defer ticker.Stop()
		log.Info("Started recurring tasks worker", zap.String("interval", interval.String()))
		for {
			select {
			case <-ctx.Done():
				log.Info("Stopping recurring tasks worker")
				return
			case t := <-ticker.C:
				if err := repo.ProcessDueRecurringTasks(ctx, t, calculateNextRun); err != nil {
					log.Error("failed to process recurring tasks", zap.String("error", err.Error()))
				}
			}
		}
	}()
}
