package domain

import (
	"fmt"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) Validate() error {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return fmt.Errorf("invalid `Priority`: %q: %w", string(p), core_errors.ErrInvalidArgument)
	}
}
