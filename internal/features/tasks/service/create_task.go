package tasks_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	if _, err := s.usersChecker.GetUser(ctx, task.AuthorUserID); err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return domain.Task{}, fmt.Errorf(
				"author user with id='%d' does not exist: %w", task.AuthorUserID, core_errors.ErrInvalidArgument,
			)
		}

		return domain.Task{}, fmt.Errorf("check author user existence: %w", err)
	}

	task, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task in repository: %w", err)
	}

	return task, nil
}
