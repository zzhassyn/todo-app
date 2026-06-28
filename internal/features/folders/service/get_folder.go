package folders_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// GetFolder returns the folder if it exists and belongs to
// requestingUserID, otherwise ErrNotFound — the same "404, not 403"
// pattern used everywhere else in this project so folder IDs belonging to
// other users can't be distinguished from nonexistent ones.
func (s *FoldersService) GetFolder(ctx context.Context, id uuid.UUID, requestingUserID int) (domain.Folder, error) {
	folder, err := s.foldersRepository.GetFolder(ctx, id)
	if err != nil {
		return domain.Folder{}, fmt.Errorf("get folder: %w", err)
	}

	if folder.UserID != requestingUserID {
		return domain.Folder{}, fmt.Errorf("folder with id='%s': %w", id, core_errors.ErrNotFound)
	}

	return folder, nil
}
