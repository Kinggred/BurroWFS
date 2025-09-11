package webdav

import (
	"burrowfs/core/db/models"
	"time"
)

type MultistatusSchema struct {
	Responses []ResourceSchema
}

func ParseFilesToMultistatus(files []models.File) MultistatusSchema {
	multistatus := MultistatusSchema{}
	for _, file := range files {
		response := ParseFileToResponse(file)
		multistatus.Responses = append(multistatus.Responses, response)
	}
	return multistatus
}

type ResourceSchema struct {
	Href          string
	DisplayName   string
	IsCollection  bool
	LastModified  string // RFC 1123 format
	CreatedAt     string // RFC 3339 format
	ContentType   string
	ContentLength int64
	Status        string
}

func ParseFileToResponse(file models.File) ResourceSchema {
	return ResourceSchema{
		Href:          file.Path,
		DisplayName:   file.Name,
		IsCollection:  file.FileID == nil,
		LastModified:  file.UpdatedAt.Format(time.RFC1123),
		CreatedAt:     file.CreatedAt.Format(time.RFC3339),
		ContentType:   file.ContentType,
		ContentLength: file.Size,
		Status:        "HTTP/1.1 200 OK",
	}
}
