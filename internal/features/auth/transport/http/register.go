package auth_transport_http

import (
	"net/http"

	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

type RegisterRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8,max=72"`
}

type RegisterResponse struct {
	User UserDTOResponse `json:"user"`
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	result, err := h.authService.Register(
		ctx,
		request.FullName,
		request.PhoneNumber,
		request.Email,
		request.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	h.setTokenCookie(w, result.Token)

	response := RegisterResponse{User: userDTOFromDomain(result.User)}
	responseHandler.JSONResponse(response, http.StatusCreated)
}
