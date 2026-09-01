package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_request "github.com/zzhassyn/todo-app/internal/core/transport/http/request"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
	core_http_types "github.com/zzhassyn/todo-app/internal/core/transport/http/types"
	core_http_utils "github.com/zzhassyn/todo-app/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` cannot be null")
		}

		fullNamelen := len([]rune(*r.FullName.Value))
		if fullNamelen < 3 || fullNamelen > 100 {
			return fmt.Errorf("'FullName' must be between 3 and 100 characters long")
		}
	}
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumber := len([]rune(*r.PhoneNumber.Value))
			if phoneNumber < 10 || phoneNumber > 15 {
				return fmt.Errorf("'PhoneNumber' must be between 10 and 15 characters long")
			}

			if !strings.HasPrefix((*r.PhoneNumber.Value), "+") {
				return fmt.Errorf("'PhoneNumber' must start with a '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID from path parameter")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("invalid request body: %v: %w", err, core_errors.ErrInvalidArgument), "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
