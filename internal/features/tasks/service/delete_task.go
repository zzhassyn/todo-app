package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(ctx context.Context, id int, requestingUserID int) error {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return err
	}

	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}
	return nil
}
