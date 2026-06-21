package domain

import (
	"fmt"
	"time"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(title string, description *string, authorUserID int) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Time{},
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("invalid `Title` length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))
		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf(
				"invalid `Description` length: %d: %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.AuthorUserID <= 0 {
		return fmt.Errorf("invalid `AuthorUserID`: %d: %w", t.AuthorUserID, core_errors.ErrInvalidArgument)
	}

	if !t.Completed && t.CompletedAt != nil {
		return fmt.Errorf("`CompletedAt` must be nil when task is not completed: %w", core_errors.ErrInvalidArgument)
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("`CompletedAt` must be set when task is completed: %w", core_errors.ErrInvalidArgument)
		}

		if !t.CreatedAt.IsZero() && t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func (t *Task) Complete(now time.Time) error {
	if t.Completed {
		return fmt.Errorf("task is already completed: %w", core_errors.ErrConflict)
	}

	t.Completed = true
	t.CompletedAt = &now

	return t.Validate()
}

func (t *Task) Uncomplete() error {
	if !t.Completed {
		return fmt.Errorf("task is already not completed: %w", core_errors.ErrConflict)
	}

	t.Completed = false
	t.CompletedAt = nil

	return t.Validate()
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp
	return nil
}

func NewTaskPatch(title Nullable[string], description Nullable[string]) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
	}
}
