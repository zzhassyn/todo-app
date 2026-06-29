package subtasks_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type CreateSubtaskRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Position int    `json:"position" validate:"min=0"`
}

type PatchSubtaskRequest struct {
	Title    *string `json:"title" validate:"omitempty,min=1,max=100"`
	Position *int    `json:"position" validate:"omitempty,min=0"`
}

type ReorderSubtasksRequest struct {
	SubtaskIDs []int `json:"subtask_ids" validate:"required"`
}

func (h *SubtasksHTTPHandler) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to create subtask")
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid task id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	var request CreateSubtaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	subtask, err := h.subtasksService.CreateSubtask(ctx, taskID, request.Title, request.Position, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create subtask")
		return
	}

	responseHandler.JSONResponse(subtaskDTOFromDomain(subtask), http.StatusCreated)
}

func (h *SubtasksHTTPHandler) PatchSubtask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to patch subtask")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	var request PatchSubtaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	patch := domain.NewSubtaskPatch(
		domain.NewNullableFromPointer(request.Title),
		domain.NewNullableFromPointer(request.Position),
	)

	subtask, err := h.subtasksService.PatchSubtask(ctx, id, patch, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch subtask")
		return
	}

	responseHandler.JSONResponse(subtaskDTOFromDomain(subtask), http.StatusOK)
}

func (h *SubtasksHTTPHandler) CompleteSubtask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to complete subtask")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	subtask, err := h.subtasksService.CompleteSubtask(ctx, id, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to complete subtask")
		return
	}

	responseHandler.JSONResponse(subtaskDTOFromDomain(subtask), http.StatusOK)
}

func (h *SubtasksHTTPHandler) UncompleteSubtask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to uncomplete subtask")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	subtask, err := h.subtasksService.UncompleteSubtask(ctx, id, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to uncomplete subtask")
		return
	}

	responseHandler.JSONResponse(subtaskDTOFromDomain(subtask), http.StatusOK)
}

func (h *SubtasksHTTPHandler) DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to delete subtask")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	if err := h.subtasksService.DeleteSubtask(ctx, id, claims.UserID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete subtask")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}

func (h *SubtasksHTTPHandler) ReorderSubtasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("no user: %w", core_errors.ErrUnauthorized), "failed to reorder subtasks")
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid task id: %w", core_errors.ErrInvalidArgument), "failed to parse id")
		return
	}

	var request ReorderSubtasksRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	if err := h.subtasksService.ReorderSubtasks(ctx, taskID, request.SubtaskIDs, claims.UserID); err != nil {
		responseHandler.ErrorResponse(err, "failed to reorder subtasks")
		return
	}

	responseHandler.JSONResponse(nil, http.StatusNoContent)
}
