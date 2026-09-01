package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// PermanentlyDeleteTask hard-deletes a task row. The repository itself
// does not check archival state — that's a service-layer policy decision
// (see tasks_service.PermanentlyDeleteTask) — but the only code path that
// calls this requires the task to already be archived first, so a task
// can never be hard-deleted without having gone through the archive step.
func (r *TasksRepository) PermanentlyDeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.tasks WHERE id = $1;`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with id='%d': %w", id, core_errors.ErrNotFound)
	}
	return nil
}
