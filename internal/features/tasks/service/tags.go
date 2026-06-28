package tasks_service

import (
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

// maxTagsPerTask caps how many tags can be attached to one task. This is a
// sanity bound, not a domain-significant number — it exists to keep the
// task_tags table and tag pickers in the UI from growing unbounded.
const maxTagsPerTask = 20

func validateTagNames(names []string) error {
	if len(names) > maxTagsPerTask {
		return fmt.Errorf(
			"too many tags: %d (max %d): %w", len(names), maxTagsPerTask, core_errors.ErrInvalidArgument,
		)
	}

	for _, name := range names {
		tag := domain.NewTagUninitialized(name)
		if err := tag.Validate(); err != nil {
			return fmt.Errorf("tag %q: %w", name, err)
		}
	}

	return nil
}
