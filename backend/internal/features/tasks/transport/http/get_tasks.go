package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
	core_http_utils "github.com/zzhassyn/todo-app/internal/core/transport/http/utils"
	tasks_service "github.com/zzhassyn/todo-app/internal/features/tasks/service"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks always scopes results to the authenticated user: the
// author_user_id is taken from the JWT, never from a query parameter, so
// one user cannot enumerate another user's tasks by guessing IDs. Client-
// supplied filters are `completed`, `archived`, `priority`, `tag`,
// `folder_id`, `limit`, and `offset`. `archived` defaults to false (only
// non-archived tasks). `folder_id` accepts either a UUID (only tasks in
// that folder) or the literal "none" (only unfiled tasks); omitting it
// does not filter by folder at all.
func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to get tasks",
		)
		return
	}

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit' and 'offset' query params")
		return
	}

	completed, err := core_http_utils.GetBoolQueryParam(r, "completed")
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("get `completed` query param: %w", err),
			"failed to get tasks filter query params",
		)
		return
	}

	archived, err := core_http_utils.GetBoolQueryParam(r, "archived")
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("get `archived` query param: %w", err),
			"failed to get tasks filter query params",
		)
		return
	}

	priority, err := getPriorityQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks filter query params")
		return
	}

	var tag *string
	if v := r.URL.Query().Get("tag"); v != "" {
		tag = &v
	}

	var search *string
	if v := r.URL.Query().Get("search"); v != "" {
		search = &v
	}

	folderID, folderIDProvided, err := core_http_utils.GetUUIDQueryParam(r, "folder_id")
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("get `folder_id` query param: %w", err),
			"failed to get tasks filter query params",
		)
		return
	}

	authorUserID := claims.UserID
	filter := tasks_service.TasksFilter{
		AuthorUserID: &authorUserID,
		Completed:    completed,
		Archived:     archived,
		Priority:     priority,
		Tag:          tag,
		FolderID:     folderID,
		NoFolder:     folderIDProvided && folderID == nil,
		Search:       search,
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, filter, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(taskDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_utils.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` query param: %w", err)
	}

	return limit, offset, nil
}

func getPriorityQueryParam(r *http.Request) (*domain.Priority, error) {
	v := r.URL.Query().Get("priority")
	if v == "" {
		return nil, nil
	}

	priority := domain.Priority(v)
	if err := priority.Validate(); err != nil {
		return nil, fmt.Errorf("get `priority` query param: %w", err)
	}

	return &priority, nil
}
