package rest

import "burrowfs/core/db/models"

type GetResponseSchema struct {
	Files []FileJsonResponse `json:"files"`
}

type FileJsonResponse struct {
	Id           string `json:"id"`
	OwnerId      string `json:"owner_id"`
	PathParentId string `json:"path_parent_id"`
	ContentType  string `json:"content_type"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	Version      int    `json:"version"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func NewFileJsonResponse(file models.File) FileJsonResponse {
	return FileJsonResponse{
		Id:           file.ID.String(),
		OwnerId:      file.OwnerID.String(),
		PathParentId: file.PathParentID.String(),
		ContentType:  file.ContentType,
		Name:         file.Name,
		Path:         file.Path,
		Size:         file.Size,
		Version:      file.Version,
		CreatedAt:    file.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func NewGetResponseSchema(files []models.File) GetResponseSchema {
	response := GetResponseSchema{}
	for _, file := range files {
		fileResponse := NewFileJsonResponse(file)
		response.Files = append(response.Files, fileResponse)
	}
	return response
}
