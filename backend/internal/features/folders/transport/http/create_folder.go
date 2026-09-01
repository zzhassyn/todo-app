package folders_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type CreateFolderRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
}

type CreateFolderResponse FolderDTOResponse

func (h *FoldersHTTPHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to create folder",
		)
		return
	}

	var request CreateFolderRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	folderDomain := domain.NewFolderUninitialized(claims.UserID, request.Title)

	folderDomain, err := h.foldersService.CreateFolder(ctx, folderDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create folder")
		return
	}

	response := CreateFolderResponse(folderDTOFromDomain(folderDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
