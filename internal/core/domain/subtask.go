package domain

import (
	"fmt"
	"time"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

type Subtask struct {
	ID          int
	TaskID      int
	Title       string
	CompletedAt *time.Time
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubtaskParams struct {
	ID          int
	TaskID      int
	Title       string
	CompletedAt *time.Time
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSubtask(p SubtaskParams) Subtask {
	return Subtask{
		ID:          p.ID,
		TaskID:      p.TaskID,
		Title:       p.Title,
		CompletedAt: p.CompletedAt,
		Position:    p.Position,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func NewSubtaskUninitialized(
	taskID int,
	title string,
	position int,
) Subtask {
	return NewSubtask(SubtaskParams{
		ID:          UninitializedID,
		TaskID:      taskID,
		Title:       title,
		CompletedAt: nil,
		Position:    position,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	})
}

func (s *Subtask) Validate() error {
	titleLength := len([]rune(s.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("invalid `Title` length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	if s.TaskID <= 0 {
		return fmt.Errorf("invalid `TaskID`: %d: %w", s.TaskID, core_errors.ErrInvalidArgument)
	}

	if s.CompletedAt != nil && !s.CreatedAt.IsZero() && s.CompletedAt.Before(s.CreatedAt) {
		return fmt.Errorf("`CompletedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (s *Subtask) Complete(now time.Time) error {
	if s.CompletedAt != nil {
		return fmt.Errorf("subtask is already completed: %w", core_errors.ErrConflict)
	}
	s.CompletedAt = &now
	return s.Validate()
}

func (s *Subtask) Uncomplete() error {
	if s.CompletedAt == nil {
		return fmt.Errorf("subtask is already not completed: %w", core_errors.ErrConflict)
	}
	s.CompletedAt = nil
	return s.Validate()
}

type SubtaskPatch struct {
	Title    Nullable[string]
	Position Nullable[int]
}

func (p *SubtaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	if p.Position.Set && p.Position.Value == nil {
		return fmt.Errorf("`Position` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (s *Subtask) ApplyPatch(patch SubtaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate subtask patch: %w", err)
	}

	tmp := *s

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Position.Set {
		tmp.Position = *patch.Position.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched subtask: %w", err)
	}

	*s = tmp
	return nil
}

func NewSubtaskPatch(
	title Nullable[string],
	position Nullable[int],
) SubtaskPatch {
	return SubtaskPatch{
		Title:    title,
		Position: position,
	}
}
