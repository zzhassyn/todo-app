package tasks_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, id int, requestingUserID int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
