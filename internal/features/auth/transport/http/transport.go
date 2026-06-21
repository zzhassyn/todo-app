package auth_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
	auth_service "github.com/zzhassyn/todo-app/internal/features/auth/service"
)

type AuthHTTPHandler struct {
	authService AuthService
	cookieName  string
	tokenTTL    time.Duration
	// cookieSecure controls the cookie's Secure attribute. It should be
	// true in production (HTTPS only); kept false by default for local
	// HTTP development.
	cookieSecure bool
}

type AuthService interface {
	Register(
		ctx context.Context,
		fullName string,
		phoneNumber *string,
		email string,
		plaintextPassword string,
	) (auth_service.RegisterResult, error)

	Login(
		ctx context.Context,
		email string,
		plaintextPassword string,
	) (auth_service.LoginResult, error)

	Me(ctx context.Context, userID int) (domain.User, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	cookieName string,
	tokenTTL time.Duration,
	cookieSecure bool,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService:  authService,
		cookieName:   cookieName,
		tokenTTL:     tokenTTL,
		cookieSecure: cookieSecure,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}

// AuthenticatedRoutes returns routes that require a valid session. The
// caller is expected to register these behind the Auth middleware (see
// cmd/todoapp/main.go), separately from Routes(), which are public.
func (h *AuthHTTPHandler) AuthenticatedRoutes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/auth/me",
			Handler: h.Me,
		},
	}
}

func (h *AuthHTTPHandler) setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.tokenTTL.Seconds()),
	})
}

func (h *AuthHTTPHandler) clearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
