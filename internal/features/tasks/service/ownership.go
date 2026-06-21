package tasks_service

import (
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// requireOwner returns ErrNotFound (not ErrConflict/forbidden) when the
// task does not belong to requestingUserID. Responding as if the resource
// does not exist — rather than revealing that it exists but belongs to
// someone else — avoids leaking which task IDs are in use by other users.
func requireOwner(task domain.Task, requestingUserID int) error {
	if task.AuthorUserID != requestingUserID {
		return fmt.Errorf("task with id='%d': %w", task.ID, core_errors.ErrNotFound)
	}
	return nil
}
