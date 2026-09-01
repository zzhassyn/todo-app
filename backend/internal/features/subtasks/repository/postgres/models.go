package subtasks_postgres_repository

import (
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type SubtaskModel struct {
	ID          int
	TaskID      int
	Title       string
	CompletedAt *time.Time
	Position    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func subtaskDomainFromModel(model SubtaskModel) domain.Subtask {
	return domain.NewSubtask(domain.SubtaskParams{
		ID:          model.ID,
		TaskID:      model.TaskID,
		Title:       model.Title,
		CompletedAt: model.CompletedAt,
		Position:    model.Position,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	})
}

func subtasksDomainFromModels(models []SubtaskModel) []domain.Subtask {
	if len(models) == 0 {
		return nil
	}

	result := make([]domain.Subtask, 0, len(models))
	for _, m := range models {
		result = append(result, subtaskDomainFromModel(m))
	}

	return result
}
