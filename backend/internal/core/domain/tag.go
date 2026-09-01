package domain

import (
	"fmt"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

type Tag struct {
	ID   int
	Name string
}

func NewTag(id int, name string) Tag {
	return Tag{ID: id, Name: name}
}

func NewTagUninitialized(name string) Tag {
	return NewTag(UninitializedID, name)
}

func (t *Tag) Validate() error {
	nameLength := len([]rune(t.Name))
	if nameLength < 1 || nameLength > 50 {
		return fmt.Errorf("invalid `Name` length: %d: %w", nameLength, core_errors.ErrInvalidArgument)
	}
	return nil
}
