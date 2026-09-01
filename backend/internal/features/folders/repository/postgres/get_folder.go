package folders_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

func (r *FoldersRepository) GetFolder(ctx context.Context, id uuid.UUID) (domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, title, created_at
		FROM todoapp.folders
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var folderModel FolderModel
	err := row.Scan(
		&folderModel.ID,
		&folderModel.UserID,
		&folderModel.Title,
		&folderModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Folder{}, fmt.Errorf("folder with id='%s': %w", id, core_errors.ErrNotFound)
		}

		return domain.Folder{}, fmt.Errorf("scan folder row: %w", err)
	}

	return folderDomainFromModel(folderModel), nil
}
