package folders_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *FoldersService) CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error) {
	if err := folder.Validate(); err != nil {
		return domain.Folder{}, fmt.Errorf("validate folder domain: %w", err)
	}

	folder, err := s.foldersRepository.CreateFolder(ctx, folder)
	if err != nil {
		return domain.Folder{}, fmt.Errorf("create folder in repository: %w", err)
	}

	return folder, nil
}
