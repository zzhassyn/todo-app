package tasks_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// checkFolderOwnership verifies that folderID, if set, both exists and
// belongs to userID. A nil folderID (the task is unfiled) always passes.
// Returns ErrInvalidArgument (not ErrNotFound) on failure: from the task
// caller's point of view, an unusable folder_id is a bad request body,
// not a missing resource of their own.
func (s *TasksService) checkFolderOwnership(ctx context.Context, folderID *uuid.UUID, userID int) error {
	if folderID == nil {
		return nil
	}

	if _, err := s.foldersChecker.GetFolder(ctx, *folderID, userID); err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return fmt.Errorf(
				"folder with id='%s' does not exist: %w", *folderID, core_errors.ErrInvalidArgument,
			)
		}

		return fmt.Errorf("check folder ownership: %w", err)
	}

	return nil
}
