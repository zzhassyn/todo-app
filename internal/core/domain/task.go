package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
	Priority     Priority
	DueAt        *time.Time
	ArchivedAt   *time.Time
	Tags         []Tag
	FolderID     *uuid.UUID
}

// TaskParams groups the fields needed to construct a Task. The struct
// exists because Task has grown enough fields (priority, due date,
// archival, tags, folder) that a long positional constructor would be
// error-prone to call correctly; named fields make call sites
// self-documenting.
type TaskParams struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
	Priority     Priority
	DueAt        *time.Time
	ArchivedAt   *time.Time
	Tags         []Tag
	FolderID     *uuid.UUID
}

func NewTask(p TaskParams) Task {
	return Task{
		ID:           p.ID,
		Version:      p.Version,
		Title:        p.Title,
		Description:  p.Description,
		Completed:    p.Completed,
		CreatedAt:    p.CreatedAt,
		CompletedAt:  p.CompletedAt,
		AuthorUserID: p.AuthorUserID,
		Priority:     p.Priority,
		DueAt:        p.DueAt,
		ArchivedAt:   p.ArchivedAt,
		Tags:         p.Tags,
		FolderID:     p.FolderID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
	priority Priority,
	dueAt *time.Time,
	folderID *uuid.UUID,
) Task {
	return NewTask(TaskParams{
		ID:           UninitializedID,
		Version:      UninitializedVersion,
		Title:        title,
		Description:  description,
		Completed:    false,
		CreatedAt:    time.Time{},
		CompletedAt:  nil,
		AuthorUserID: authorUserID,
		Priority:     priority,
		DueAt:        dueAt,
		ArchivedAt:   nil,
		Tags:         nil,
		FolderID:     folderID,
	})
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

	if err := t.Priority.Validate(); err != nil {
		return fmt.Errorf("validate priority: %w", err)
	}

	if t.ArchivedAt != nil && !t.CreatedAt.IsZero() && t.ArchivedAt.Before(t.CreatedAt) {
		return fmt.Errorf("`ArchivedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	if t.FolderID != nil && *t.FolderID == uuid.Nil {
		return fmt.Errorf("`FolderID` can't be the nil UUID: %w", core_errors.ErrInvalidArgument)
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

// IsArchived reports whether the task has been soft-deleted.
func (t *Task) IsArchived() bool {
	return t.ArchivedAt != nil
}

func (t *Task) Archive(now time.Time) error {
	if t.IsArchived() {
		return fmt.Errorf("task is already archived: %w", core_errors.ErrConflict)
	}

	t.ArchivedAt = &now

	return t.Validate()
}

func (t *Task) Unarchive() error {
	if !t.IsArchived() {
		return fmt.Errorf("task is not archived: %w", core_errors.ErrConflict)
	}

	t.ArchivedAt = nil

	return t.Validate()
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Priority    Nullable[Priority]
	DueAt       Nullable[time.Time]
	// FolderID follows the same Nullable convention as DueAt: Set=true,
	// Value=nil moves the task out of any folder (back to the unfiled
	// "buffer"); Set=true, Value!=nil moves it into that folder.
	FolderID Nullable[uuid.UUID]
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Priority.Set {
		if p.Priority.Value == nil {
			return fmt.Errorf("`Priority` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}
		if err := p.Priority.Value.Validate(); err != nil {
			return fmt.Errorf("validate patched priority: %w", err)
		}
	}

	if p.FolderID.Set && p.FolderID.Value != nil && *p.FolderID.Value == uuid.Nil {
		return fmt.Errorf("`FolderID` can't be the nil UUID: %w", core_errors.ErrInvalidArgument)
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

	if patch.Priority.Set {
		tmp.Priority = *patch.Priority.Value
	}

	// DueAt follows the same Nullable convention as Description: Set=true,
	// Value=nil clears the due date; Set=true, Value!=nil sets it.
	if patch.DueAt.Set {
		tmp.DueAt = patch.DueAt.Value
	}

	if patch.FolderID.Set {
		tmp.FolderID = patch.FolderID.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp
	return nil
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	priority Nullable[Priority],
	dueAt Nullable[time.Time],
	folderID Nullable[uuid.UUID],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Priority:    priority,
		DueAt:       dueAt,
		FolderID:    folderID,
	}
}
