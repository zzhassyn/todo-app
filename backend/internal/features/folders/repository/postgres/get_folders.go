package folders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (r *FoldersRepository) GetFolders(ctx context.Context, userID int) ([]domain.Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, title, created_at
		FROM todoapp.folders
		WHERE user_id = $1
		ORDER BY created_at ASC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select folders: %w", err)
	}
	defer rows.Close()

	var folderModels []FolderModel
	for rows.Next() {
		var folderModel FolderModel
		if err := rows.Scan(
			&folderModel.ID,
			&folderModel.UserID,
			&folderModel.Title,
			&folderModel.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan folder: %w", err)
		}
		folderModels = append(folderModels, folderModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return folderDomainsFromModels(folderModels), nil
}
