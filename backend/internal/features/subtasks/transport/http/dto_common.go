package subtasks_transport_http

import (
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type SubtaskDTOResponse struct {
	ID          int        `json:"id"`
	TaskID      int        `json:"task_id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completed_at"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func subtaskDTOFromDomain(domainSubtask domain.Subtask) SubtaskDTOResponse {
	return SubtaskDTOResponse{
		ID:          domainSubtask.ID,
		TaskID:      domainSubtask.TaskID,
		Title:       domainSubtask.Title,
		CompletedAt: domainSubtask.CompletedAt,
		Position:    domainSubtask.Position,
		CreatedAt:   domainSubtask.CreatedAt,
		UpdatedAt:   domainSubtask.UpdatedAt,
	}
}

func subtasksDTOFromDomain(domainSubtasks []domain.Subtask) []SubtaskDTOResponse {
	if len(domainSubtasks) == 0 {
		return []SubtaskDTOResponse{}
	}

	result := make([]SubtaskDTOResponse, 0, len(domainSubtasks))
	for _, subtask := range domainSubtasks {
		result = append(result, subtaskDTOFromDomain(subtask))
	}

	return result
}
