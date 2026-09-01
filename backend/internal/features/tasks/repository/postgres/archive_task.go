package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

// ArchiveTask persists a task's archived/unarchived state (and whatever
// other fields the caller has already mutated on the domain.Task, the
// same get-modify-put pattern used by PatchTask). It does not touch tags.
func (r *TasksRepository) ArchiveTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.tasks
		SET
			archived_at=$1,
			version=version+1
		WHERE id=$2 AND version=$3
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id, priority, due_at, archived_at, folder_id;
	`

	row := r.pool.QueryRow(ctx, query,
		task.ArchivedAt,
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

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
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
