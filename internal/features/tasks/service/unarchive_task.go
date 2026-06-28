package tasks_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) UnarchiveTask(ctx context.Context, id int, requestingUserID int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return domain.Task{}, err
	}

	if err := task.Unarchive(); err != nil {
		return domain.Task{}, fmt.Errorf("unarchive task: %w", err)
	}

	unarchivedTask, err := s.tasksRepository.ArchiveTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("unarchive task in repository: %w", err)
	}

	return unarchivedTask, nil
}
