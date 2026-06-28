package tasks_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	filter TasksFilter,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit <= 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w", core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w", core_errors.ErrInvalidArgument,
		)
	}

	if filter.AuthorUserID != nil && *filter.AuthorUserID <= 0 {
		return nil, fmt.Errorf(
			"author_user_id filter must be positive: %w", core_errors.ErrInvalidArgument,
		)
	}

	if filter.Priority != nil {
		if err := filter.Priority.Validate(); err != nil {
			return nil, fmt.Errorf("priority filter: %w", err)
		}
	}

	if filter.FolderID != nil && filter.AuthorUserID != nil {
		if err := s.checkFolderOwnership(ctx, filter.FolderID, *filter.AuthorUserID); err != nil {
			return nil, err
		}
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}
	return tasks, nil
}
