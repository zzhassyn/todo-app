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
	authService       AuthService
	cookieName        string
	tokenTTL          time.Duration
	refreshCookieName string
	refreshTokenTTL   time.Duration
	cookieSecure      bool
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

	Refresh(ctx context.Context, refreshToken string) (auth_service.RefreshResult, error)

	Me(ctx context.Context, userID int) (domain.User, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	cookieName string,
	tokenTTL time.Duration,
	refreshCookieName string,
	refreshTokenTTL time.Duration,
	cookieSecure bool,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService:       authService,
		cookieName:        cookieName,
		tokenTTL:          tokenTTL,
		refreshCookieName: refreshCookieName,
		refreshTokenTTL:   refreshTokenTTL,
		cookieSecure:      cookieSecure,
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
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.Logout,
		},
	}
}

func (h *AuthHTTPHandler) AuthenticatedRoutes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/auth/me",
			Handler: h.Me,
		},
	}
}

func (h *AuthHTTPHandler) setTokensCookies(w http.ResponseWriter, accessToken string, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.tokenTTL.Seconds()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshCookieName,
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.refreshTokenTTL.Seconds()),
	})
}

func (h *AuthHTTPHandler) clearTokensCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshCookieName,
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
