package auth_transport_http

import (
	"errors"
	"net/http"

	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type RefreshResponse struct {
	User UserDTOResponse `json:"user"`
}

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	cookie, err := r.Cookie(h.refreshCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			responseHandler.ErrorResponse(errors.New("refresh token cookie missing"), "unauthorized")
			return
		}
		responseHandler.ErrorResponse(err, "failed to read refresh token cookie")
		return
	}

	result, err := h.authService.Refresh(ctx, cookie.Value)
	if err != nil {
		// If refresh fails, clear the cookies so the client is forced to login
		h.clearTokensCookies(w)
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	h.setTokensCookies(w, result.Token, result.RefreshToken)

	response := RefreshResponse{User: userDTOFromDomain(result.User)}
	responseHandler.JSONResponse(response, http.StatusOK)
}
