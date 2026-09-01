package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type GetTagsResponse []TagDTOResponse

func (h *TasksHTTPHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	tags, err := h.tasksService.GetTags(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tags")
		return
	}

	response := GetTagsResponse(tagsDTOFromDomains(tags))
	responseHandler.JSONResponse(response, http.StatusOK)
}
