package tasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type BulkPatchParams struct {
	TaskIDs []int
	Patch   domain.TaskPatch
	Tags    []string
}

func (s *TasksService) BulkPatchTasks(
	ctx context.Context,
	requestingUserID int,
	params BulkPatchParams,
) ([]domain.Task, error) {
	var patchedTasks []domain.Task

	// Validate tags once if provided
	if params.Tags != nil {
		if err := validateTagNames(params.Tags); err != nil {
			return nil, fmt.Errorf("validate tags: %w", err)
		}
	}

	for _, id := range params.TaskIDs {
		// Patch each task individually using existing logic to ensure all rules apply
		patched, err := s.PatchTask(ctx, id, requestingUserID, params.Patch, params.Tags)
		if err != nil {
			// Fail fast on error. In a more robust system, we might return a partial success response.
			return nil, fmt.Errorf("patch task %d: %w", id, err)
		}
		patchedTasks = append(patchedTasks, patched)
	}

	return patchedTasks, nil
}

func (s *TasksService) BulkCompleteTasks(
	ctx context.Context,
	requestingUserID int,
	taskIDs []int,
) ([]domain.Task, error) {
	var completedTasks []domain.Task
	now := time.Now().UTC()

	for _, id := range taskIDs {
		task, err := s.tasksRepository.GetTask(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get task %d: %w", id, err)
		}

		if err := requireOwner(task, requestingUserID); err != nil {
			return nil, err
		}

		if err := task.Complete(now); err != nil {
			// Skip already completed tasks or invalid ones
			continue 
		}

		patched, err := s.tasksRepository.PatchTask(ctx, id, task, nil)
		if err != nil {
			return nil, fmt.Errorf("patch task %d: %w", id, err)
		}
		completedTasks = append(completedTasks, patched)
	}
	return completedTasks, nil
}

func (s *TasksService) BulkArchiveTasks(
	ctx context.Context,
	requestingUserID int,
	taskIDs []int,
) ([]domain.Task, error) {
	var archivedTasks []domain.Task
	now := time.Now().UTC()

	for _, id := range taskIDs {
		task, err := s.tasksRepository.GetTask(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get task %d: %w", id, err)
		}

		if err := requireOwner(task, requestingUserID); err != nil {
			return nil, err
		}

		if err := task.Archive(now); err != nil {
			continue 
		}

		archived, err := s.tasksRepository.ArchiveTask(ctx, id, task)
		if err != nil {
			return nil, fmt.Errorf("archive task %d: %w", id, err)
		}
		archivedTasks = append(archivedTasks, archived)
	}
	return archivedTasks, nil
}
