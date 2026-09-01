package subtasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (s *SubtasksService) checkTaskAccess(ctx context.Context, taskID int, requestingUserID int) error {
	task, err := s.tasksChecker.GetTask(ctx, taskID, requestingUserID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if task.AuthorUserID != requestingUserID {
		return fmt.Errorf("task belongs to another user: %w", core_errors.ErrUnauthorized)
	}

	return nil
}

func (s *SubtasksService) CreateSubtask(
	ctx context.Context,
	taskID int,
	title string,
	position int,
	requestingUserID int,
) (domain.Subtask, error) {
	if err := s.checkTaskAccess(ctx, taskID, requestingUserID); err != nil {
		return domain.Subtask{}, err
	}

	subtask := domain.NewSubtaskUninitialized(taskID, title, position)
	if err := subtask.Validate(); err != nil {
		return domain.Subtask{}, fmt.Errorf("validate subtask: %w", err)
	}

	created, err := s.subtasksRepository.CreateSubtask(ctx, subtask)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("create subtask: %w", err)
	}

	return created, nil
}

func (s *SubtasksService) PatchSubtask(
	ctx context.Context,
	id int,
	patch domain.SubtaskPatch,
	requestingUserID int,
) (domain.Subtask, error) {
	subtask, err := s.subtasksRepository.GetSubtask(ctx, id)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("get subtask: %w", err)
	}

	if err := s.checkTaskAccess(ctx, subtask.TaskID, requestingUserID); err != nil {
		return domain.Subtask{}, err
	}

	if err := subtask.ApplyPatch(patch); err != nil {
		return domain.Subtask{}, fmt.Errorf("apply patch: %w", err)
	}

	updated, err := s.subtasksRepository.PatchSubtask(ctx, id, subtask)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("patch subtask: %w", err)
	}

	return updated, nil
}

func (s *SubtasksService) CompleteSubtask(
	ctx context.Context,
	id int,
	requestingUserID int,
) (domain.Subtask, error) {
	subtask, err := s.subtasksRepository.GetSubtask(ctx, id)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("get subtask: %w", err)
	}

	if err := s.checkTaskAccess(ctx, subtask.TaskID, requestingUserID); err != nil {
		return domain.Subtask{}, err
	}

	if err := subtask.Complete(time.Now()); err != nil {
		return domain.Subtask{}, fmt.Errorf("complete subtask: %w", err)
	}

	updated, err := s.subtasksRepository.PatchSubtask(ctx, id, subtask)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("patch subtask: %w", err)
	}

	return updated, nil
}

func (s *SubtasksService) UncompleteSubtask(
	ctx context.Context,
	id int,
	requestingUserID int,
) (domain.Subtask, error) {
	subtask, err := s.subtasksRepository.GetSubtask(ctx, id)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("get subtask: %w", err)
	}

	if err := s.checkTaskAccess(ctx, subtask.TaskID, requestingUserID); err != nil {
		return domain.Subtask{}, err
	}

	if err := subtask.Uncomplete(); err != nil {
		return domain.Subtask{}, fmt.Errorf("uncomplete subtask: %w", err)
	}

	updated, err := s.subtasksRepository.PatchSubtask(ctx, id, subtask)
	if err != nil {
		return domain.Subtask{}, fmt.Errorf("patch subtask: %w", err)
	}

	return updated, nil
}

func (s *SubtasksService) DeleteSubtask(
	ctx context.Context,
	id int,
	requestingUserID int,
) error {
	subtask, err := s.subtasksRepository.GetSubtask(ctx, id)
	if err != nil {
		return fmt.Errorf("get subtask: %w", err)
	}

	if err := s.checkTaskAccess(ctx, subtask.TaskID, requestingUserID); err != nil {
		return err
	}

	if err := s.subtasksRepository.DeleteSubtask(ctx, id); err != nil {
		return fmt.Errorf("delete subtask: %w", err)
	}

	return nil
}

func (s *SubtasksService) ReorderSubtasks(
	ctx context.Context,
	taskID int,
	subtaskIDs []int,
	requestingUserID int,
) error {
	if err := s.checkTaskAccess(ctx, taskID, requestingUserID); err != nil {
		return err
	}

	if err := s.subtasksRepository.ReorderSubtasks(ctx, taskID, subtaskIDs); err != nil {
		return fmt.Errorf("reorder subtasks: %w", err)
	}

	return nil
}
