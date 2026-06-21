package tasks_service

import (
	"context"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
	usersChecker    UsersChecker
}

// TasksFilter is a service-level (domain-oriented) filter for listing tasks.
// It is intentionally decoupled from any repository implementation detail.
type TasksFilter struct {
	AuthorUserID *int
	Completed    *bool
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		filter TasksFilter,
		limit *int,
		offset *int,
	) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
	PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error)
}

// UsersChecker is the subset of the users feature's service that the tasks
// feature depends on. The tasks feature does not depend on the users
// repository directly; it depends on this narrow interface, which is
// satisfied by users_service.UsersService.
type UsersChecker interface {
	GetUser(ctx context.Context, id int) (domain.User, error)
}

func NewTasksService(tasksRepository TasksRepository, usersChecker UsersChecker) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
		usersChecker:    usersChecker,
	}
}
