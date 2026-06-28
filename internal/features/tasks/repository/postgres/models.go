package tasks_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
	Priority     string
	DueAt        *time.Time
	ArchivedAt   *time.Time
	FolderID     *uuid.UUID
}

type TagModel struct {
	ID   int
	Name string
}

func tagDomainFromModel(tag TagModel) domain.Tag {
	return domain.NewTag(tag.ID, tag.Name)
}

func tagDomainsFromModels(tags []TagModel) []domain.Tag {
	tagDomains := make([]domain.Tag, len(tags))
	for i, tag := range tags {
		tagDomains[i] = tagDomainFromModel(tag)
	}
	return tagDomains
}

// taskDomainFromModel converts a row plus its already-loaded tags into a
// domain.Task. Tags are fetched separately (see get_task.go/get_tasks.go)
// rather than via a SQL-side JSON aggregate, keeping the query simple and
// the row-scanning code uniform across all task queries.
func taskDomainFromModel(task TaskModel, tags []domain.Tag) domain.Task {
	return domain.NewTask(domain.TaskParams{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
		Priority:     domain.Priority(task.Priority),
		DueAt:        task.DueAt,
		ArchivedAt:   task.ArchivedAt,
		Tags:         tags,
		FolderID:     task.FolderID,
	})
}

// taskDomainsFromModels converts task rows into domain.Tasks, attaching
// tags from tagsByTaskID (keyed by task ID; missing entries mean "no tags").
func taskDomainsFromModels(tasks []TaskModel, tagsByTaskID map[int][]domain.Tag) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		taskDomains[i] = taskDomainFromModel(task, tagsByTaskID[task.ID])
	}
	return taskDomains
}
