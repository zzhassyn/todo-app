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
)

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=100"`
	Description *string    `json:"description" validate:"omitempty,min=1,max=1000"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueAt       *time.Time `json:"due_at"`
	Tags        []string   `json:"tags" validate:"omitempty,max=20,dive,min=1,max=50"`
	FolderID    *uuid.UUID `json:"folder_id"`
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to create task",
		)
		return
	}

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	taskDomain := taskDomainFromCreateRequest(request, claims.UserID)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain, request.Tags)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")

		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func taskDomainFromCreateRequest(dto CreateTaskRequest, authorUserID int) domain.Task {
	priority := domain.Priority(dto.Priority)
	if priority == "" {
		priority = domain.PriorityMedium
	}

	return domain.NewTaskUninitialized(dto.Title, dto.Description, authorUserID, priority, dto.DueAt, dto.FolderID)
}
