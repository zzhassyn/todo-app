package auth_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type MeResponse UserDTOResponse

func (h *AuthHTTPHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	claims, ok := core_auth.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("no authenticated user in context: %w", core_errors.ErrUnauthorized),
			"failed to get current user",
		)
		return
	}

	user, err := h.authService.Me(ctx, claims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get current user")
		return
	}

	response := MeResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
