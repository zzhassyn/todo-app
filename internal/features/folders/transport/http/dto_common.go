package folders_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type FolderDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func folderDTOFromDomain(folder domain.Folder) FolderDTOResponse {
	return FolderDTOResponse{
		ID:        folder.ID,
		Title:     folder.Title,
		CreatedAt: folder.CreatedAt,
	}
}

func foldersDTOFromDomains(folders []domain.Folder) []FolderDTOResponse {
	foldersDTO := make([]FolderDTOResponse, len(folders))
	for i, folder := range folders {
		foldersDTO[i] = folderDTOFromDomain(folder)
	}
	return foldersDTO
}
