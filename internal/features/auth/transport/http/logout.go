package auth_transport_http

import (
	"net/http"

	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	h.clearTokenCookie(w)

	responseHandler.NoContentResponse()
}
