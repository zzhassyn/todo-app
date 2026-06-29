package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

// PatchTask updates all mutable task fields (title, description,
// completion state, priority, due date, folder) and, when tagNames is
// non-nil, replaces the task's tag set entirely. Passing a nil tagNames
// leaves tags untouched — this lets callers that only change completion
// state (e.g. CompleteTask/UncompleteTask in the service) reuse this
// method without needing to know or re-supply the task's current tags.
// Passing a non-nil but empty slice clears all tags, which is the correct
// way for an actual PATCH request to remove every tag.
func (r *TasksRepository) PatchTask(
	ctx context.Context,
	id int,
	task domain.Task,
	tagNames []string,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a documented no-op

	query := `
		UPDATE todoapp.tasks
		SET
			title=$1,
			description=$2,
			completed=$3,
			completed_at=$4,
			priority=$5,
			due_at=$6,
			folder_id=$7,
			version=version+1
		WHERE id=$8 AND version=$9
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id, priority, due_at, archived_at, folder_id;
	`

	row := tx.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		string(task.Priority),
		task.DueAt,
		task.FolderID,
		id,
		task.Version,
	)

	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
		&taskModel.Priority,
		&taskModel.DueAt,
		&taskModel.ArchivedAt,
		&taskModel.FolderID,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d' concurrently accessed: %w", id, core_errors.ErrConflict)
		}

		if errors.Is(err, core_postgres_pool.ErrForeignKeyViolation) {
			return domain.Task{}, fmt.Errorf("folder does not exist: %w", core_errors.ErrInvalidArgument)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	if tagNames != nil {
		tagIDs, err := upsertTagsTx(ctx, tx, tagNames)
		if err != nil {
			return domain.Task{}, fmt.Errorf("upsert tags: %w", err)
		}

		if err := replaceTaskTagsTx(ctx, tx, taskModel.ID, tagIDs); err != nil {
			return domain.Task{}, fmt.Errorf("link task tags: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Task{}, fmt.Errorf("commit tx: %w", err)
	}

	tagsByTaskID, err := r.loadTagsByTaskID(ctx, []int{taskModel.ID})
	if err != nil {
		return domain.Task{}, fmt.Errorf("load tags: %w", err)
	}

	subtasksByTaskID, err := r.loadSubtasksByTaskID(ctx, []int{taskModel.ID})
	if err != nil {
		return domain.Task{}, fmt.Errorf("load subtasks: %w", err)
	}

	return taskDomainFromModel(taskModel, tagsByTaskID[taskModel.ID], subtasksByTaskID[taskModel.ID]), nil
}
