package tasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *TasksService) ArchiveTask(ctx context.Context, id int, requestingUserID int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := requireOwner(task, requestingUserID); err != nil {
		return domain.Task{}, err
	}

	if err := task.Archive(time.Now().UTC()); err != nil {
		return domain.Task{}, fmt.Errorf("archive task: %w", err)
	}

	archivedTask, err := s.tasksRepository.ArchiveTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("archive task in repository: %w", err)
	}

	return archivedTask, nil
}
