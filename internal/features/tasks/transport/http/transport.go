package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_http_middleware "github.com/zzhassyn/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
	tasks_service "github.com/zzhassyn/todo-app/internal/features/tasks/service"
)

type TasksHTTPHandler struct {
	tasksService   TasksService
	authMiddleware core_http_middleware.Middleware
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		filter tasks_service.TasksFilter,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) error

	PatchTask(
		ctx context.Context,
		id int,
		requestingUserID int,
		patch domain.TaskPatch,
	) (domain.Task, error)

	CompleteTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Task, error)

	UncompleteTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(tasksService TasksService, authMiddleware core_http_middleware.Middleware) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService:   tasksService,
		authMiddleware: authMiddleware,
	}
}

// Routes returns all task routes. Every route requires authentication
// (a valid JWT cookie), since tasks are always scoped to their author.
func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	mw := []core_http_middleware.Middleware{h.authMiddleware}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/tasks",
			Handler:    h.CreateTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks",
			Handler:    h.GetTasks,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks/{id}",
			Handler:    h.GetTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/tasks/{id}",
			Handler:    h.DeleteTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/tasks/{id}",
			Handler:    h.PatchTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{id}/complete",
			Handler:    h.CompleteTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{id}/uncomplete",
			Handler:    h.UncompleteTask,
			Middleware: mw,
		},
	}
}
