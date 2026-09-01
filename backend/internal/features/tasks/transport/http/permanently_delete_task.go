package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/zzhassyn/todo-app/internal/core/transport/http/utils"
)

// PermanentlyDeleteTask handles DELETE /tasks/{id}/permanent. Unlike
// DELETE /tasks/{id} (which archives), this hard-deletes the row and only
// succeeds for tasks that are already archived — see
// tasks_service.PermanentlyDeleteTask for why.
func (h *TasksHTTPHandler) PermanentlyDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to permanently delete task",
		)
		return
	}

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path")
		return
	}

	if err := h.tasksService.PermanentlyDeleteTask(ctx, taskID, claims.UserID); err != nil {
		responseHandler.ErrorResponse(err, "failed to permanently delete task")
		return
	}

	responseHandler.NoContentResponse()
}
