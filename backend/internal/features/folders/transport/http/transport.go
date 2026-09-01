package folders_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_http_middleware "github.com/zzhassyn/todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/zzhassyn/todo-app/internal/core/transport/http/server"
)

type FoldersHTTPHandler struct {
	foldersService FoldersService
	authMiddleware core_http_middleware.Middleware
}

type FoldersService interface {
	CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error)
	GetFolders(ctx context.Context, userID int) ([]domain.Folder, error)
	DeleteFolder(ctx context.Context, id uuid.UUID, requestingUserID int) error
}

func NewFoldersHTTPHandler(foldersService FoldersService, authMiddleware core_http_middleware.Middleware) *FoldersHTTPHandler {
	return &FoldersHTTPHandler{
		foldersService: foldersService,
		authMiddleware: authMiddleware,
	}
}

// Routes returns all folder routes. Every route requires authentication,
// since folders are always scoped to their owner.
func (h *FoldersHTTPHandler) Routes() []core_http_server.Route {
	mw := []core_http_middleware.Middleware{h.authMiddleware}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/folders",
			Handler:    h.CreateFolder,
			Middleware: mw,
		},
		{
			Method:     http.MethodGet,
			Path:       "/folders",
			Handler:    h.GetFolders,
			Middleware: mw,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/folders/{id}",
			Handler:    h.DeleteFolder,
			Middleware: mw,
		},
	}
}
