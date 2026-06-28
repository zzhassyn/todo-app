package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

// upsertTagsTx finds or creates a tag row for each name and returns their
// IDs. Tag names are deduplicated; ON CONFLICT DO NOTHING + a follow-up
// SELECT handles the find-or-create race safely under concurrent inserts
// of the same tag name.
func upsertTagsTx(ctx context.Context, tx core_postgres_pool.Tx, names []string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}

	deduped := dedupeStrings(names)

	insertQuery := `
		INSERT INTO todoapp.tags (name)
		VALUES (unnest($1::varchar[]))
		ON CONFLICT (name) DO NOTHING;
	`
	if _, err := tx.Exec(ctx, insertQuery, deduped); err != nil {
		return nil, fmt.Errorf("insert tags: %w", err)
	}

	selectQuery := `
		SELECT id FROM todoapp.tags WHERE name = ANY($1::varchar[]);
	`
	rows, err := tx.Query(ctx, selectQuery, deduped)
	if err != nil {
		return nil, fmt.Errorf("select tag ids: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tag id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return ids, nil
}

// replaceTaskTagsTx replaces the full set of tags attached to taskID with
// tagIDs (delete-then-insert; simplest correct approach for the small,
// per-task tag counts this app deals with).
func replaceTaskTagsTx(ctx context.Context, tx core_postgres_pool.Tx, taskID int, tagIDs []int) error {
	if _, err := tx.Exec(ctx, `DELETE FROM todoapp.task_tags WHERE task_id = $1;`, taskID); err != nil {
		return fmt.Errorf("delete existing task tags: %w", err)
	}

	if len(tagIDs) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO todoapp.task_tags (task_id, tag_id)
		SELECT $1, unnest($2::int[]);
	`
	if _, err := tx.Exec(ctx, insertQuery, taskID, tagIDs); err != nil {
		return fmt.Errorf("insert task tags: %w", err)
	}

	return nil
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
