package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

// loadTagsByTaskID fetches all tags attached to the given task IDs and
// groups them by task ID. Returns an empty map (not an error) if taskIDs
// is empty or none of the tasks have tags.
func (r *TasksRepository) loadTagsByTaskID(ctx context.Context, taskIDs []int) (map[int][]domain.Tag, error) {
	result := make(map[int][]domain.Tag)
	if len(taskIDs) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT tt.task_id, t.id, t.name
		FROM todoapp.task_tags tt
		JOIN todoapp.tags t ON t.id = tt.tag_id
		WHERE tt.task_id = ANY($1)
		ORDER BY t.name ASC;
	`

	rows, err := r.pool.Query(ctx, query, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("select task tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			taskID int
			tag    TagModel
		)
		if err := rows.Scan(&taskID, &tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("scan task tag: %w", err)
		}
		result[taskID] = append(result[taskID], tagDomainFromModel(tag))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return result, nil
}
