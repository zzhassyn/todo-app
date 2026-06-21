package core_http_middleware

import (
	"fmt"
	"net/http"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_http_response "github.com/zzhassyn/todo-app/internal/core/transport/http/response"
)

// Auth returns a middleware that requires a valid JWT, supplied via an
// httpOnly cookie (config.CookieName). On success, the decoded claims are
// stored in the request context and can be retrieved downstream via
// core_auth.FromContext. On failure, it responds with 401 and does not
// call the next handler.
func Auth(config core_auth.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			cookie, err := r.Cookie(config.CookieName)
			if err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("missing auth cookie: %v: %w", err, core_errors.ErrUnauthorized),
					"failed to authenticate request",
				)
				return
			}

			claims, err := core_auth.ParseToken(config.JWTSecret, cookie.Value)
			if err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("invalid auth token: %v: %w", err, core_errors.ErrUnauthorized),
					"failed to authenticate request",
				)
				return
			}

			ctx = core_auth.ToContext(ctx, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
