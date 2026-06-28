package folders_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type FoldersService struct {
	foldersRepository FoldersRepository
}

type FoldersRepository interface {
	CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error)
	GetFolders(ctx context.Context, userID int) ([]domain.Folder, error)
	GetFolder(ctx context.Context, id uuid.UUID) (domain.Folder, error)
	DeleteFolder(ctx context.Context, id uuid.UUID) error
}

func NewFoldersService(foldersRepository FoldersRepository) *FoldersService {
	return &FoldersService{
		foldersRepository: foldersRepository,
	}
}
