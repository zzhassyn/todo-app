package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/zzhassyn/todo-app/internal/core/transport/http/types"
	tasks_service "github.com/zzhassyn/todo-app/internal/features/tasks/service"
)

type BulkActionRequest struct {
	TaskIDs []int `json:"task_ids" validate:"required,min=1"`
}

type BulkPatchRequest struct {
	TaskIDs     []int                               `json:"task_ids" validate:"required,min=1"`
	Title       core_http_types.Nullable[string]    `json:"title"`
	Description core_http_types.Nullable[string]    `json:"description"`
	Priority    core_http_types.Nullable[string]    `json:"priority"`
	DueAt       core_http_types.Nullable[time.Time] `json:"due_at"`
	Tags        []string                            `json:"tags" validate:"omitempty,max=20,dive,min=1,max=50"`
	FolderID    core_http_types.Nullable[uuid.UUID] `json:"folder_id"`
}

func (h *TasksHTTPHandler) BulkPatchTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	var request BulkPatchRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode request")
		return
	}

	var priority domain.Nullable[domain.Priority]
	if request.Priority.Set {
		priority.Set = true
		if request.Priority.Value != nil {
			p := domain.Priority(*request.Priority.Value)
			priority.Value = &p
		}
	}

	patch := domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		priority,
		request.DueAt.ToDomain(),
		request.FolderID.ToDomain(),
		domain.Nullable[float64]{Set: false, Value: nil}, // Bulk position patching not supported here
	)

	// In patch logic, if Tags is provided (non-nil array), it's applied
	params := tasks_service.BulkPatchParams{
		TaskIDs: request.TaskIDs,
		Patch:   patch,
		Tags:    request.Tags,
	}

	patchedTasks, err := h.tasksService.BulkPatchTasks(ctx, claims.UserID, params)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch tasks")
		return
	}

	var taskDTOs []TaskDTOResponse
	for _, task := range patchedTasks {
		taskDTOs = append(taskDTOs, taskDTOFromDomain(task))
	}

	responseHandler.JSONResponse(BulkResponse{Tasks: taskDTOs}, http.StatusOK)
}

type BulkResponse struct {
	Tasks []TaskDTOResponse `json:"tasks"`
}

func (h *TasksHTTPHandler) BulkCompleteTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	var request BulkActionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode request")
		return
	}

	patchedTasks, err := h.tasksService.BulkCompleteTasks(ctx, claims.UserID, request.TaskIDs)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to complete tasks")
		return
	}

	var taskDTOs []TaskDTOResponse
	for _, task := range patchedTasks {
		taskDTOs = append(taskDTOs, taskDTOFromDomain(task))
	}

	responseHandler.JSONResponse(BulkResponse{Tasks: taskDTOs}, http.StatusOK)
}

func (h *TasksHTTPHandler) BulkArchiveTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	var request BulkActionRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode request")
		return
	}

	archivedTasks, err := h.tasksService.BulkArchiveTasks(ctx, claims.UserID, request.TaskIDs)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to archive tasks")
		return
	}

	var taskDTOs []TaskDTOResponse
	for _, task := range archivedTasks {
		taskDTOs = append(taskDTOs, taskDTOFromDomain(task))
	}

	responseHandler.JSONResponse(BulkResponse{Tasks: taskDTOs}, http.StatusOK)
}
