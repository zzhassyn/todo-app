package users_transport_http

import (
	"context"
	"net/http"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_http_middleware "github.com/zzhassyn/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService   UsersService
	authMiddleware core_http_middleware.Middleware
}

// UsersService is the subset of users_service.UsersService that this HTTP
// handler needs. Note that user creation is intentionally NOT exposed here:
// registration is handled exclusively by the auth feature (POST
// /auth/register), which calls usersService.CreateUser directly after
// hashing the password.
type UsersService interface {
	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		key int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService, authMiddleware core_http_middleware.Middleware) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService:   usersService,
		authMiddleware: authMiddleware,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	mw := []core_http_middleware.Middleware{h.authMiddleware}

	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/users",
			Handler:    h.GetUsers,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/users/{id}",
			Handler:    h.GetUser,
			Middleware: mw,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/users/{id}",
			Handler:    h.DeleteUser,
			Middleware: mw,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/users/{id}",
			Handler:    h.PatchUser,
			Middleware: mw,
		},
	}
}
