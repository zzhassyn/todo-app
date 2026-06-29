package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (r *TasksRepository) loadSubtasksByTaskID(ctx context.Context, taskIDs []int) (map[int][]domain.Subtask, error) {
	result := make(map[int][]domain.Subtask)
	if len(taskIDs) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, task_id, title, completed_at, position, created_at, updated_at
		FROM todoapp.subtasks
		WHERE task_id = ANY($1)
		ORDER BY position ASC;
	`

	rows, err := r.pool.Query(ctx, query, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("select subtasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var model SubtaskModel
		if err := rows.Scan(
			&model.ID,
			&model.TaskID,
			&model.Title,
			&model.CompletedAt,
			&model.Position,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan subtask: %w", err)
		}
		result[model.TaskID] = append(result[model.TaskID], subtaskDomainFromModel(model))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return result, nil
}
