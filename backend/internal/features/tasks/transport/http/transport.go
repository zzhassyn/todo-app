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
		tagNames []string,
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

	// ArchiveTask soft-deletes a task. There is no hard-delete endpoint —
	// archiving is the only way to remove a task from the default view.
	ArchiveTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Task, error)

	UnarchiveTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Task, error)

	PatchTask(
		ctx context.Context,
		id int,
		requestingUserID int,
		patch domain.TaskPatch,
		tagNames []string,
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

	// PermanentlyDeleteTask hard-deletes a task. Only succeeds if the task
	// is already archived (see tasks_service.PermanentlyDeleteTask) — it
	// is reachable exclusively from the archive view in the UI.
	PermanentlyDeleteTask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) error

	BulkPatchTasks(
		ctx context.Context,
		requestingUserID int,
		params tasks_service.BulkPatchParams,
	) ([]domain.Task, error)

	BulkCompleteTasks(
		ctx context.Context,
		requestingUserID int,
		taskIDs []int,
	) ([]domain.Task, error)

	BulkArchiveTasks(
		ctx context.Context,
		requestingUserID int,
		taskIDs []int,
	) ([]domain.Task, error)

	CreateRecurringTask(ctx context.Context, task domain.RecurringTask) (domain.RecurringTask, error)
	GetRecurringTasks(ctx context.Context, userID int) ([]domain.RecurringTask, error)
	DeleteRecurringTask(ctx context.Context, id int, userID int) error

	GetTags(ctx context.Context) ([]domain.Tag, error)
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
			// DELETE archives rather than hard-deletes — see ArchiveTask.
			Method:     http.MethodDelete,
			Path:       "/tasks/{id}",
			Handler:    h.ArchiveTask,
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
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{id}/archive",
			Handler:    h.ArchiveTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{id}/unarchive",
			Handler:    h.UnarchiveTask,
			Middleware: mw,
		},
		{
			// Hard delete. Only succeeds for already-archived tasks — see
			// tasks_service.PermanentlyDeleteTask.
			Method:     http.MethodDelete,
			Path:       "/tasks/{id}/permanent",
			Handler:    h.PermanentlyDeleteTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/tasks/bulk",
			Handler:    h.BulkPatchTasks,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/bulk/complete",
			Handler:    h.BulkCompleteTasks,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/bulk/archive",
			Handler:    h.BulkArchiveTasks,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/recurring-tasks",
			Handler:    h.CreateRecurringTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/recurring-tasks",
			Handler:    h.GetRecurringTasks,
			Middleware: mw,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/recurring-tasks/{id}",
			Handler:    h.DeleteRecurringTask,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/tags",
			Handler:    h.GetTags,
			Middleware: mw,
		},
	}
}
