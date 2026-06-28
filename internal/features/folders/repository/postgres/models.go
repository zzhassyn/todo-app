package folders_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type FolderModel struct {
	ID        uuid.UUID
	UserID    int
	Title     string
	CreatedAt time.Time
}

func folderDomainFromModel(folder FolderModel) domain.Folder {
	return domain.NewFolder(folder.ID, folder.UserID, folder.Title, folder.CreatedAt)
}

func folderDomainsFromModels(folders []FolderModel) []domain.Folder {
	folderDomains := make([]domain.Folder, len(folders))
	for i, folder := range folders {
		folderDomains[i] = folderDomainFromModel(folder)
	}
	return folderDomains
}
