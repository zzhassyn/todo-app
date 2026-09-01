package subtasks_transport_http

import (
	"context"
	"net/http"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_http_middleware "github.com/zzhassyn/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
)

type SubtasksHTTPHandler struct {
	subtasksService SubtasksService
	authMiddleware  core_http_middleware.Middleware
}

type SubtasksService interface {
	CreateSubtask(
		ctx context.Context,
		taskID int,
		title string,
		position int,
		requestingUserID int,
	) (domain.Subtask, error)

	PatchSubtask(
		ctx context.Context,
		id int,
		patch domain.SubtaskPatch,
		requestingUserID int,
	) (domain.Subtask, error)

	CompleteSubtask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Subtask, error)

	UncompleteSubtask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) (domain.Subtask, error)

	DeleteSubtask(
		ctx context.Context,
		id int,
		requestingUserID int,
	) error

	ReorderSubtasks(
		ctx context.Context,
		taskID int,
		subtaskIDs []int,
		requestingUserID int,
	) error
}

func NewSubtasksHTTPHandler(subtasksService SubtasksService, authMiddleware core_http_middleware.Middleware) *SubtasksHTTPHandler {
	return &SubtasksHTTPHandler{
		subtasksService: subtasksService,
		authMiddleware:  authMiddleware,
	}
}

func (h *SubtasksHTTPHandler) Routes() []core_http_server.Route {
	mw := []core_http_middleware.Middleware{h.authMiddleware}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{taskId}/subtasks",
			Handler:    h.CreateSubtask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/subtasks/{id}",
			Handler:    h.PatchSubtask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/subtasks/{id}/complete",
			Handler:    h.CompleteSubtask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/subtasks/{id}/uncomplete",
			Handler:    h.UncompleteSubtask,
			Middleware: mw,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/subtasks/{id}",
			Handler:    h.DeleteSubtask,
			Middleware: mw,
		},
		{
			Method:     http.MethodPost,
			Path:       "/tasks/{taskId}/subtasks/reorder",
			Handler:    h.ReorderSubtasks,
			Middleware: mw,
		},
	}
}
