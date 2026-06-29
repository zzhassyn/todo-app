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
		joins      []string
	)

	if filter.AuthorUserID != nil {
		args = append(args, *filter.AuthorUserID)
		conditions = append(conditions, fmt.Sprintf("t.author_user_id = $%d", len(args)))
	}

	if filter.Completed != nil {
		args = append(args, *filter.Completed)
		conditions = append(conditions, fmt.Sprintf("t.completed = $%d", len(args)))
	}

	// Archived semantics: nil or false => only non-archived tasks (the
	// "normal" view); true => only archived tasks (the "archive" view).
	// There is intentionally no way to request "both" — callers wanting
	// that should issue two requests, since mixing them is rarely useful
	// and would complicate the default behavior.
	if filter.Archived != nil && *filter.Archived {
		conditions = append(conditions, "t.archived_at IS NOT NULL")
	} else {
		conditions = append(conditions, "t.archived_at IS NULL")
	}

	if filter.Priority != nil {
		args = append(args, string(*filter.Priority))
		conditions = append(conditions, fmt.Sprintf("t.priority = $%d", len(args)))
	}

	if filter.Tag != nil {
		joins = append(joins, "JOIN todoapp.task_tags ft ON ft.task_id = t.id JOIN todoapp.tags tg ON tg.id = ft.tag_id")
		args = append(args, *filter.Tag)
		conditions = append(conditions, fmt.Sprintf("tg.name = $%d", len(args)))
	}

	// Folder semantics: NoFolder=true means "only unfiled tasks" (the
	// default buffer view); FolderID!=nil means "only tasks in this
	// folder". Both unset (the zero TasksFilter) means "don't filter by
	// folder at all" — neither is mutually exclusive at the type level,
	// but the service layer only ever sets one of them at a time.
	if filter.NoFolder {
		conditions = append(conditions, "t.folder_id IS NULL")
	} else if filter.FolderID != nil {
		args = append(args, *filter.FolderID)
		conditions = append(conditions, fmt.Sprintf("t.folder_id = $%d", len(args)))
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
		SELECT DISTINCT t.id, t.version, t.title, t.description, t.completed, t.created_at,
		       t.completed_at, t.author_user_id, t.priority, t.due_at, t.archived_at, t.folder_id, t.position
		FROM todoapp.tasks t
		%s
		%s
		ORDER BY t.position ASC, t.id ASC
		LIMIT %s OFFSET %s;
	`, strings.Join(joins, "\n"), whereClause, limitPlaceholder, offsetPlaceholder)

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
			&taskModel.Priority,
			&taskModel.DueAt,
			&taskModel.ArchivedAt,
			&taskModel.FolderID,
			&taskModel.Position,
		); err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskIDs := make([]int, len(taskModels))
	for i, tm := range taskModels {
		taskIDs[i] = tm.ID
	}

	tagsByTaskID, err := r.loadTagsByTaskID(ctx, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("load tags: %w", err)
	}

	subtasksByTaskID, err := r.loadSubtasksByTaskID(ctx, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("load subtasks: %w", err)
	}

	return taskDomainsFromModels(taskModels, tagsByTaskID, subtasksByTaskID), nil
}
