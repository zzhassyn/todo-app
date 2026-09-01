package subtasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

func (r *SubtasksRepository) CreateSubtask(ctx context.Context, subtask domain.Subtask) (domain.Subtask, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.subtasks (task_id, title, position, completed_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, task_id, title, completed_at, position, created_at, updated_at;
	`

	var model SubtaskModel
	err := r.pool.QueryRow(ctx, query,
		subtask.TaskID,
		subtask.Title,
		subtask.Position,
		subtask.CompletedAt,
	).Scan(
		&model.ID,
		&model.TaskID,
		&model.Title,
		&model.CompletedAt,
		&model.Position,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrForeignKeyViolation) {
			return domain.Subtask{}, fmt.Errorf("task does not exist: %w", core_errors.ErrInvalidArgument)
		}
		return domain.Subtask{}, fmt.Errorf("scan error: %w", err)
	}

	return subtaskDomainFromModel(model), nil
}

func (r *SubtasksRepository) PatchSubtask(ctx context.Context, id int, subtask domain.Subtask) (domain.Subtask, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.subtasks
	SET title = $1, position = $2, completed_at = $3, updated_at = NOW()
	WHERE id = $4
	RETURNING id, task_id, title, completed_at, position, created_at, updated_at;
	`

	var model SubtaskModel
	err := r.pool.QueryRow(ctx, query,
		subtask.Title,
		subtask.Position,
		subtask.CompletedAt,
		id,
	).Scan(
		&model.ID,
		&model.TaskID,
		&model.Title,
		&model.CompletedAt,
		&model.Position,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subtask{}, fmt.Errorf("subtask not found: %w", core_errors.ErrNotFound)
		}
		return domain.Subtask{}, fmt.Errorf("scan error: %w", err)
	}

	return subtaskDomainFromModel(model), nil
}

func (r *SubtasksRepository) GetSubtask(ctx context.Context, id int) (domain.Subtask, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, task_id, title, completed_at, position, created_at, updated_at
	FROM todoapp.subtasks
	WHERE id = $1;
	`

	var model SubtaskModel
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.TaskID,
		&model.Title,
		&model.CompletedAt,
		&model.Position,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subtask{}, fmt.Errorf("subtask not found: %w", core_errors.ErrNotFound)
		}
		return domain.Subtask{}, fmt.Errorf("scan error: %w", err)
	}

	return subtaskDomainFromModel(model), nil
}

func (r *SubtasksRepository) DeleteSubtask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.subtasks
	WHERE id = $1;
	`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("subtask not found: %w", core_errors.ErrNotFound)
	}

	return nil
}

func (r *SubtasksRepository) ReorderSubtasks(ctx context.Context, taskID int, subtaskIDs []int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
	UPDATE todoapp.subtasks
	SET position = $1, updated_at = NOW()
	WHERE id = $2 AND task_id = $3;
	`

	for idx, subtaskID := range subtaskIDs {
		// Wait, pgx batching is better but looping is okay for small arrays.
		// Since we have tx, let's just loop.
		_, err := tx.Exec(ctx, query, idx, subtaskID, taskID)
		if err != nil {
			return fmt.Errorf("update position %d for subtask %d: %w", idx, subtaskID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
