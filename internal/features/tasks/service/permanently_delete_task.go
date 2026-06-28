package tasks_service

import (
	"context"
	"fmt"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// PermanentlyDeleteTask hard-deletes a task, but only if it is already
// archived. This is intentional: the regular delete/archive endpoint
// (ArchiveTask) is the safety net that prevents accidental data loss, and
// permanent deletion is reachable only as a second, explicit step from
// the archive view — never directly from an active task. A task that
// isn't archived yet gets ErrConflict, not silently archived-then-deleted,
// so the caller (and the UI) can't skip the confirmation step by mistake.
func (s *TasksService) PermanentlyDeleteTask(ctx context.Context, id int, requestingUserID int) error {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return err
	}

	if !task.IsArchived() {
		return fmt.Errorf(
			"task with id='%d' must be archived before it can be permanently deleted: %w",
			id, core_errors.ErrConflict,
		)
	}

	if err := s.tasksRepository.PermanentlyDeleteTask(ctx, id); err != nil {
		return fmt.Errorf("permanently delete task from repository: %w", err)
	}
	return nil
}
