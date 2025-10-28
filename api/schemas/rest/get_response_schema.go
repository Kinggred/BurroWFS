package rest

import (
	"burrowfs/core/types"
)

type GetResponseSchema struct {
	Files []FileJsonResponse `json:"files"`
}

type FileJsonResponse struct {
	Id           string `json:"id"`
	FileID       string `json:"file_id"`
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

func NewFileJsonResponse(file types.FileDTO) FileJsonResponse {
	var pathParentId string
	if file.PathParentID != nil {
		pathParentId = file.PathParentID.String()
	} else {
		pathParentId = ""
	}

	return FileJsonResponse{
		Id:           file.ID.String(),
		OwnerId:      file.OwnerID.String(),
		PathParentId: pathParentId,
		ContentType:  file.ContentType,
		Name:         file.Name,
		Path:         file.Path,
		Size:         file.Size,
		Version:      file.Version,
		CreatedAt:    file.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func NewGetResponseSchema(files []types.FileDTO) GetResponseSchema {
	response := GetResponseSchema{}
	for _, file := range files {
		fileResponse := NewFileJsonResponse(file)
		response.Files = append(response.Files, fileResponse)
	}
	return response
}
