package tasks_postgres_repository

import (
	"context"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (r *TasksRepository) CreateRecurringTask(ctx context.Context, task domain.RecurringTask) (domain.RecurringTask, error) {
	sql := `
		INSERT INTO todoapp.recurring_tasks (
			author_user_id, title, description, priority, folder_id, tags, cron_expression, next_run_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(
		ctx, sql,
		task.AuthorUserID, task.Title, task.Description, task.Priority, task.FolderID, task.Tags, task.CronExpression, task.NextRunAt,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	
	if err != nil {
		return domain.RecurringTask{}, err
	}
	return task, nil
}

func (r *TasksRepository) GetRecurringTasks(ctx context.Context, userID int) ([]domain.RecurringTask, error) {
	sql := `
		SELECT id, author_user_id, title, description, priority, folder_id, tags, cron_expression, last_run_at, next_run_at, created_at, updated_at
		FROM todoapp.recurring_tasks
		WHERE author_user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.RecurringTask
	for rows.Next() {
		var t domain.RecurringTask
		if err := rows.Scan(
			&t.ID, &t.AuthorUserID, &t.Title, &t.Description, &t.Priority, &t.FolderID, &t.Tags,
			&t.CronExpression, &t.LastRunAt, &t.NextRunAt, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TasksRepository) DeleteRecurringTask(ctx context.Context, id int, userID int) error {
	sql := `DELETE FROM todoapp.recurring_tasks WHERE id = $1 AND author_user_id = $2`
	tag, err := r.pool.Exec(ctx, sql, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}
