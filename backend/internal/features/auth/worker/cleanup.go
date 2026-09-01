package auth_worker

import (
	"context"
	"time"

	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	"go.uber.org/zap"
)

// AuthRepository is the subset of the auth repository needed by the cleanup
// worker. Keeping the interface here avoids importing the full repository
// package and follows the same pattern as tasks_worker.TasksRepository.
type AuthRepository interface {
	DeleteExpiredRefreshTokens(ctx context.Context) (int64, error)
}

// StartRefreshTokenCleanupWorker launches a background goroutine that
// periodically deletes expired refresh tokens from the database. This
// prevents the refresh_tokens table from growing unboundedly as users
// log in and rotate tokens over time.
//
// The worker respects context cancellation for graceful shutdown.
func StartRefreshTokenCleanupWorker(ctx context.Context, repo AuthRepository, interval time.Duration) {
	log := core_logger.FromContext(ctx)
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		log.Info("Started refresh token cleanup worker", zap.String("interval", interval.String()))

		for {
			select {
			case <-ctx.Done():
				log.Info("Stopping refresh token cleanup worker")
				return
			case <-ticker.C:
				deleted, err := repo.DeleteExpiredRefreshTokens(ctx)
				if err != nil {
					log.Error("failed to clean up expired refresh tokens", zap.Error(err))
				} else if deleted > 0 {
					log.Info("cleaned up expired refresh tokens", zap.Int64("deleted", deleted))
				}
			}
		}
	}()
}
