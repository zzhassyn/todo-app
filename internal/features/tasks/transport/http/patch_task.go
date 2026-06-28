package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/zzhassyn/todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/zzhassyn/todo-app/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string]    `json:"title"`
	Description core_http_types.Nullable[string]    `json:"description"`
	Priority    core_http_types.Nullable[string]    `json:"priority"`
	DueAt       core_http_types.Nullable[time.Time] `json:"due_at"`
	// Tags follows the same nil-vs-empty convention used throughout the
	// tags feature: omitted or JSON null => leave tags untouched; "tags":
	// [] => clear all tags; "tags": [...] => replace with this set.
	Tags []string `json:"tags" validate:"omitempty,max=20,dive,min=1,max=50"`
	// FolderID follows the same Nullable convention as DueAt: omitted =>
	// leave the task's folder untouched; explicit JSON null => move the
	// task out of any folder (back to the unfiled buffer); a UUID =>
	// move the task into that folder.
	FolderID core_http_types.Nullable[uuid.UUID] `json:"folder_id"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` cannot be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100 characters long")
		}
	}

	if r.Description.Set && r.Description.Value != nil {
		descriptionLen := len([]rune(*r.Description.Value))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("'Description' must be between 1 and 1000 characters long")
		}
	}

	if r.Priority.Set {
		if r.Priority.Value == nil {
			return fmt.Errorf("`Priority` cannot be null")
		}
		switch *r.Priority.Value {
		case "low", "medium", "high":
		default:
			return fmt.Errorf("'Priority' must be one of: low, medium, high")
		}
	}

	if r.FolderID.Set && r.FolderID.Value != nil && *r.FolderID.Value == uuid.Nil {
		return fmt.Errorf("`FolderID` can't be the nil UUID")
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to patch task",
		)
		return
	}

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path parameter")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid request body: %v: %w", err, core_errors.ErrInvalidArgument), "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, claims.UserID, taskPatch, request.Tags)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	var priority domain.Nullable[domain.Priority]
	if request.Priority.Set {
		priority.Set = true
		if request.Priority.Value != nil {
			p := domain.Priority(*request.Priority.Value)
			priority.Value = &p
		}
	}

	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		priority,
		request.DueAt.ToDomain(),
		request.FolderID.ToDomain(),
	)
}
