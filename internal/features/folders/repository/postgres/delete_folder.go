package folders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// DeleteFolder permanently deletes a folder. Unlike tasks, folders have no
// soft-delete/archive concept. Because tasks.folder_id has ON DELETE
// CASCADE (see migration 000004_folders), this also permanently deletes
// every task inside the folder — including archived ones — bypassing the
// usual archive-first safety net for tasks entirely. This was a deliberate
// choice (not the safer "set folder_id to NULL" alternative); callers of
// this feature should make sure the UI warns users before deleting a
// non-empty folder.
func (r *FoldersRepository) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.folders WHERE id = $1;`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("folder with id='%s': %w", id, core_errors.ErrNotFound)
	}
	return nil
}
