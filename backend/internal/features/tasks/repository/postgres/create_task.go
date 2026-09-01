package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
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

	query := `INSERT INTO todoapp.tasks (title, description, completed, author_user_id, priority, due_at, folder_id, position)
	VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE(NULLIF($8::float8, 0.0), EXTRACT(EPOCH FROM NOW())))
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id, priority, due_at, archived_at, folder_id, position;
	`

	row := tx.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.AuthorUserID,
		string(task.Priority),
		task.DueAt,
		task.FolderID,
		task.Position,
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
		&taskModel.Position,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrForeignKeyViolation) {
			return domain.Task{}, fmt.Errorf("folder does not exist: %w", core_errors.ErrInvalidArgument)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	tagIDs, err := upsertTagsTx(ctx, tx, tagNames)
	if err != nil {
		return domain.Task{}, fmt.Errorf("upsert tags: %w", err)
	}

	if err := replaceTaskTagsTx(ctx, tx, taskModel.ID, tagIDs); err != nil {
		return domain.Task{}, fmt.Errorf("link task tags: %w", err)
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
