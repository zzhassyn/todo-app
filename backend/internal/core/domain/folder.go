package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

type Folder struct {
	ID        uuid.UUID
	UserID    int
	Title     string
	CreatedAt time.Time
}

func NewFolder(id uuid.UUID, userID int, title string, createdAt time.Time) Folder {
	return Folder{
		ID:        id,
		UserID:    userID,
		Title:     title,
		CreatedAt: createdAt,
	}
}

// NewFolderUninitialized builds a Folder for creation: the ID is generated
// client-side (Go) rather than left for the database to assign, since the
// project has no auto-incrementing UUID equivalent and pgx/Postgres need a
// concrete value to insert. CreatedAt is left zero; the repository fills
// it in from the DB's DEFAULT now() on insert.
func NewFolderUninitialized(userID int, title string) Folder {
	return NewFolder(uuid.New(), userID, title, time.Time{})
}

func (f *Folder) Validate() error {
	titleLength := len([]rune(f.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("invalid `Title` length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	if f.UserID <= 0 {
		return fmt.Errorf("invalid `UserID`: %d: %w", f.UserID, core_errors.ErrInvalidArgument)
	}

	return nil
}
