package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	tasks_service "github.com/zzhassyn/todo-app/internal/features/tasks/service"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	filter tasks_service.TasksFilter,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var (
		conditions []string
		args       []any
	)

	if filter.AuthorUserID != nil {
		args = append(args, *filter.AuthorUserID)
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", len(args)))
	}

	if filter.Completed != nil {
		args = append(args, *filter.Completed)
		conditions = append(conditions, fmt.Sprintf("completed = $%d", len(args)))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))

	args = append(args, offset)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args))

	query := fmt.Sprintf(`
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
		%s
		ORDER BY id ASC
		LIMIT %s OFFSET %s;
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		if err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return taskDomainsFromModels(taskModels), nil
}
