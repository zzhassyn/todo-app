package folders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (r *FoldersRepository) CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.folders (id, user_id, title)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, created_at;
	`

	row := r.pool.QueryRow(ctx, query, folder.ID, folder.UserID, folder.Title)

	var folderModel FolderModel
	if err := row.Scan(
		&folderModel.ID,
		&folderModel.UserID,
		&folderModel.Title,
		&folderModel.CreatedAt,
	); err != nil {
		return domain.Folder{}, fmt.Errorf("scan error: %w", err)
	}

	return folderDomainFromModel(folderModel), nil
}
