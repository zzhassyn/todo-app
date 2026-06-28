package tasks_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) UncompleteTask(ctx context.Context, id int, requestingUserID int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return domain.Task{}, err
	}

	if err := task.Uncomplete(); err != nil {
		return domain.Task{}, fmt.Errorf("uncomplete task: %w", err)
	}

	uncompletedTask, err := s.tasksRepository.PatchTask(ctx, id, task, nil)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return uncompletedTask, nil
}
