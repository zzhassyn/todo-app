package folders_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteFolder permanently deletes a folder (and, via ON DELETE CASCADE,
// every task inside it — see the repository layer's DeleteFolder comment).
// Like tasks, ownership is enforced by returning ErrNotFound rather than a
// forbidden error when requestingUserID does not own the folder, so folder
// IDs belonging to other users cannot be enumerated by probing this
// endpoint.
func (s *FoldersService) DeleteFolder(ctx context.Context, id uuid.UUID, requestingUserID int) error {
	if _, err := s.GetFolder(ctx, id, requestingUserID); err != nil {
		return err
	}

	if err := s.foldersRepository.DeleteFolder(ctx, id); err != nil {
		return fmt.Errorf("delete folder from repository: %w", err)
	}
	return nil
}
