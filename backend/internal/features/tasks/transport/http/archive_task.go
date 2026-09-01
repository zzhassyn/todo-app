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

type ArchiveTaskResponse TaskDTOResponse

// ArchiveTask handles both DELETE /tasks/{id} and POST /tasks/{id}/archive
// (see Routes()) — there is no hard-delete endpoint; archiving is a
// reversible soft-delete, and the response body is the archived task so
// the client can show/undo it without a follow-up GET.
func (h *TasksHTTPHandler) ArchiveTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to archive task",
		)
		return
	}

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path")
		return
	}

	task, err := h.tasksService.ArchiveTask(ctx, taskID, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to archive task")
		return
	}

	response := ArchiveTaskResponse(taskDTOFromDomain(task))
	responseHandler.JSONResponse(response, http.StatusOK)
}
