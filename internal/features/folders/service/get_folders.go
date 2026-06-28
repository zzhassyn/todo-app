package folders_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *FoldersService) GetFolders(ctx context.Context, userID int) ([]domain.Folder, error) {
	folders, err := s.foldersRepository.GetFolders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get folders from repository: %w", err)
	}
	return folders, nil
}
