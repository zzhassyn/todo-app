package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

// GetTags returns all tags in the system, ordered by name. Tags are not
// scoped per-user (they are a shared vocabulary, same as in most todo
// apps), so this is a simple unfiltered listing used to power
// autocomplete/suggestions in the UI.
func (r *TasksRepository) GetTags(ctx context.Context) ([]domain.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT id, name FROM todoapp.tags ORDER BY name ASC;`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select tags: %w", err)
	}
	defer rows.Close()

	var tagModels []TagModel
	for rows.Next() {
		var tag TagModel
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tagModels = append(tagModels, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return tagDomainsFromModels(tagModels), nil
}
