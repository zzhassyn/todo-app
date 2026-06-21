package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	core_logger "github.com/zzhassyn/todo-app/internal/core/logger"
	core_pgx_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/zzhassyn/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
	auth_service "github.com/zzhassyn/todo-app/internal/features/auth/service"
	auth_transport_http "github.com/zzhassyn/todo-app/internal/features/auth/transport/http"
	tasks_postgres_repository "github.com/zzhassyn/todo-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/zzhassyn/todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/zzhassyn/todo-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/zzhassyn/todo-app/internal/features/users/repository/postgres"
	users_service "github.com/zzhassyn/todo-app/internal/features/users/service"
	users_transport_http "github.com/zzhassyn/todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx, core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("Failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	authConfig := core_auth.NewConfigMust()
	authMiddleware := core_http_middleware.Auth(authConfig)

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, authMiddleware)

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	authService := auth_service.NewAuthService(usersService, auth_service.Config{
		JWTSecret: authConfig.JWTSecret,
		TokenTTL:  authConfig.TokenTTL,
	})
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(
		authService,
		authConfig.CookieName,
		authConfig.TokenTTL,
		authConfig.CookieSecure,
	)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository, usersService)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService, authMiddleware)

	logger.Debug("Initializing HTTP server")

	httpServerConfig := core_http_server.NewConfigMust()

	httpServer := core_http_server.NewHTTPServer(
		httpServerConfig,
		logger,
		core_http_middleware.CORS(httpServerConfig.AllowedOrigin),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVErsionV1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(authTransportHTTP.Routes()...)

	authMW := []core_http_middleware.Middleware{authMiddleware}
	for _, route := range authTransportHTTP.AuthenticatedRoutes() {
		route.Middleware = append(route.Middleware, authMW...)
		apiVersionRouter.RegisterRoutes(route)
	}

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server start error", zap.Error(err))
	}
}
