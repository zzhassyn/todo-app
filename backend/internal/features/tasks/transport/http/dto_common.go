package tasks_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type TagDTOResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func tagDTOFromDomain(tag domain.Tag) TagDTOResponse {
	return TagDTOResponse{ID: tag.ID, Name: tag.Name}
}

type SubtaskDTOResponse struct {
	ID          int        `json:"id"`
	TaskID      int        `json:"task_id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completed_at"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func subtasksDTOFromDomains(subtasks []domain.Subtask) []SubtaskDTOResponse {
	if len(subtasks) == 0 {
		return []SubtaskDTOResponse{}
	}
	result := make([]SubtaskDTOResponse, len(subtasks))
	for i, st := range subtasks {
		result[i] = SubtaskDTOResponse{
			ID:          st.ID,
			TaskID:      st.TaskID,
			Title:       st.Title,
			CompletedAt: st.CompletedAt,
			Position:    st.Position,
			CreatedAt:   st.CreatedAt,
			UpdatedAt:   st.UpdatedAt,
		}
	}
	return result
}

func tagsDTOFromDomains(tags []domain.Tag) []TagDTOResponse {
	tagsDTO := make([]TagDTOResponse, len(tags))
	for i, tag := range tags {
		tagsDTO[i] = tagDTOFromDomain(tag)
	}
	return tagsDTO
}

type TaskDTOResponse struct {
	ID           int              `json:"id"`
	Version      int              `json:"version"`
	Title        string           `json:"title"`
	Description  *string          `json:"description"`
	Completed    bool             `json:"completed"`
	CreatedAt    time.Time        `json:"created_at"`
	CompletedAt  *time.Time       `json:"completed_at"`
	AuthorUserID int              `json:"author_user_id"`
	Priority     string           `json:"priority"`
	DueAt        *time.Time       `json:"due_at"`
	ArchivedAt   *time.Time       `json:"archived_at"`
	Tags         []TagDTOResponse     `json:"tags"`
	FolderID     *uuid.UUID           `json:"folder_id"`
	Position     float64              `json:"position"`
	Subtasks     []SubtaskDTOResponse `json:"subtasks"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
		Priority:     string(task.Priority),
		DueAt:        task.DueAt,
		ArchivedAt:   task.ArchivedAt,
		Tags:         tagsDTOFromDomains(task.Tags),
		FolderID:     task.FolderID,
		Position:     task.Position,
		Subtasks:     subtasksDTOFromDomains(task.Subtasks),
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))
	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}

	return tasksDTO
}
