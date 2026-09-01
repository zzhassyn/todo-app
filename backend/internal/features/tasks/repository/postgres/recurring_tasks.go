package tasks_postgres_repository

import (
	"context"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

// ProcessDueRecurringTasks finds all recurring tasks that are due, creates regular tasks for them,
// and updates their next_run_at based on their cron expressions. 
// A function parameter is passed to calculate the next run time from the cron expression.
func (r *TasksRepository) ProcessDueRecurringTasks(
	ctx context.Context, 
	now time.Time,
	calculateNextRun func(cronExpr string, from time.Time) (time.Time, error),
) error {
	tx, err := r.pool.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Select due recurring tasks FOR UPDATE to lock them and prevent concurrent workers from processing
	rows, err := tx.Query(ctx, `
		SELECT id, author_user_id, title, description, priority, folder_id, tags, cron_expression, last_run_at, next_run_at 
		FROM todoapp.recurring_tasks 
		WHERE next_run_at <= $1 FOR UPDATE SKIP LOCKED
	`, now)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dueTasks []domain.RecurringTask
	for rows.Next() {
		var t domain.RecurringTask
		if err := rows.Scan(
			&t.ID, &t.AuthorUserID, &t.Title, &t.Description, &t.Priority, &t.FolderID, &t.Tags, 
			&t.CronExpression, &t.LastRunAt, &t.NextRunAt,
		); err != nil {
			return err
		}
		dueTasks = append(dueTasks, t)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, rt := range dueTasks {
		// 1. Create a regular task from the recurring template
		insertTaskSQL := `
			INSERT INTO todoapp.tasks (
				author_user_id, title, description, priority, folder_id, position
			)
			VALUES ($1, $2, $3, $4, $5, COALESCE((SELECT MAX(position) + 1024 FROM todoapp.tasks WHERE author_user_id = $1), 1024))
			RETURNING id
		`
		var newTaskId int
		if err := tx.QueryRow(ctx, insertTaskSQL, rt.AuthorUserID, rt.Title, rt.Description, rt.Priority, rt.FolderID).Scan(&newTaskId); err != nil {
			return err
		}

		// Insert tags if any
		if len(rt.Tags) > 0 {
			// Get or create tags
			tagIDs := make([]int, 0, len(rt.Tags))
			for _, tagName := range rt.Tags {
				var tagID int
				err := tx.QueryRow(ctx, `
					INSERT INTO todoapp.tags (name) VALUES ($1)
					ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
					RETURNING id
				`, tagName).Scan(&tagID)
				if err != nil {
					return err
				}
				tagIDs = append(tagIDs, tagID)
			}

			// Associate tags
			for _, tagID := range tagIDs {
				_, err := tx.Exec(ctx, `
					INSERT INTO todoapp.task_tags (task_id, tag_id)
					VALUES ($1, $2)
				`, newTaskId, tagID)
				if err != nil {
					return err
				}
			}
		}

		// 2. Update next_run_at for the recurring task
		nextRun, err := calculateNextRun(rt.CronExpression, now)
		if err != nil {
			// If cron expression is invalid, we might want to skip or log.
			// Returning err will rollback the tx. 
			return err
		}

		updateSQL := `
			UPDATE todoapp.recurring_tasks
			SET last_run_at = $1, next_run_at = $2, updated_at = $1
			WHERE id = $3
		`
		if _, err := tx.Exec(ctx, updateSQL, now, nextRun, rt.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
