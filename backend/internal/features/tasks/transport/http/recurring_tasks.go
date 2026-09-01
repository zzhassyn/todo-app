package tasks_transport_http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type RecurringTaskRequest struct {
	Title          string           `json:"title" validate:"required,min=1,max=100"`
	Description    *string          `json:"description,omitempty" validate:"omitempty,min=1,max=1000"`
	Priority       domain.Priority  `json:"priority" validate:"required,oneof=low medium high"`
	FolderID       *uuid.UUID       `json:"folder_id,omitempty"`
	Tags           []string         `json:"tags,omitempty"`
	CronExpression string           `json:"cron_expression" validate:"required"`
}

type RecurringTaskResponse struct {
	ID             int              `json:"id"`
	Title          string           `json:"title"`
	Description    *string          `json:"description,omitempty"`
	Priority       domain.Priority  `json:"priority"`
	FolderID       *uuid.UUID       `json:"folder_id,omitempty"`
	Tags           []string         `json:"tags"`
	CronExpression string           `json:"cron_expression"`
	LastRunAt      *time.Time       `json:"last_run_at,omitempty"`
	NextRunAt      time.Time        `json:"next_run_at"`
	CreatedAt      time.Time        `json:"created_at"`
}

func recurringTaskDTOFromDomain(t domain.RecurringTask) RecurringTaskResponse {
	tags := t.Tags
	if tags == nil {
		tags = []string{}
	}
	return RecurringTaskResponse{
		ID:             t.ID,
		Title:          t.Title,
		Description:    t.Description,
		Priority:       t.Priority,
		FolderID:       t.FolderID,
		Tags:           tags,
		CronExpression: t.CronExpression,
		LastRunAt:      t.LastRunAt,
		NextRunAt:      t.NextRunAt,
		CreatedAt:      t.CreatedAt,
	}
}

func (h *TasksHTTPHandler) CreateRecurringTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	var req RecurringTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "invalid request")
		return
	}

	task := domain.RecurringTask{
		AuthorUserID:   claims.UserID,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		FolderID:       req.FolderID,
		Tags:           req.Tags,
		CronExpression: req.CronExpression,
	}

	created, err := h.tasksService.CreateRecurringTask(ctx, task)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create recurring task")
		return
	}

	responseHandler.JSONResponse(recurringTaskDTOFromDomain(created), http.StatusCreated)
}

func (h *TasksHTTPHandler) GetRecurringTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	tasks, err := h.tasksService.GetRecurringTasks(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to fetch recurring tasks")
		return
	}

	dtos := make([]RecurringTaskResponse, 0, len(tasks))
	for _, t := range tasks {
		dtos = append(dtos, recurringTaskDTOFromDomain(t))
	}

	responseHandler.JSONResponse(dtos, http.StatusOK)
}

func (h *TasksHTTPHandler) DeleteRecurringTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("unauthorized"), "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "invalid id")
		return
	}

	if err := h.tasksService.DeleteRecurringTask(ctx, id, claims.UserID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete recurring task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
