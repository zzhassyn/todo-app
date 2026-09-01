package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

// GetTask returns a task by ID regardless of archival state — archived
// tasks remain individually retrievable (e.g. to view or unarchive them);
// only the *listing* endpoint excludes them by default.
func (r *TasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id, priority, due_at, archived_at, folder_id, position
		FROM todoapp.tasks
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel

	err := row.Scan(
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
		&taskModel.Position,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d': %w", id, core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan task row: %w", err)
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
