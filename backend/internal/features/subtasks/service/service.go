package subtasks_service

import (
	"context"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type SubtasksService struct {
	subtasksRepository SubtasksRepository
	tasksChecker       TasksChecker
}

type SubtasksRepository interface {
	CreateSubtask(ctx context.Context, subtask domain.Subtask) (domain.Subtask, error)
	PatchSubtask(ctx context.Context, id int, subtask domain.Subtask) (domain.Subtask, error)
	DeleteSubtask(ctx context.Context, id int) error
	GetSubtask(ctx context.Context, id int) (domain.Subtask, error)
	ReorderSubtasks(ctx context.Context, taskID int, subtaskIDs []int) error
}

type TasksChecker interface {
	GetTask(ctx context.Context, id int, requestingUserID int) (domain.Task, error)
}

func NewSubtasksService(subtasksRepository SubtasksRepository, tasksChecker TasksChecker) *SubtasksService {
	return &SubtasksService{
		subtasksRepository: subtasksRepository,
		tasksChecker:       tasksChecker,
	}
}
