package rest

import (
	"bytes"
	"io"
)

// PutFilesInputSchema defines the expected input for uploading a files.
type PutFilesInputSchema struct {
	Items []FileInRequest `json:"items"`
}

type FileInRequest struct {
	RelativePath string `json:"relative_path"`
	Size         int64  `json:"size"`
	FilePart     []byte `json:"file_part"`
}

func (f FileInRequest) GetReadCloser() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(f.FilePart))
}

type PutFilesResponseSchema struct {
	Items []FileInResponse `json:"items"`
}

type FileInResponse struct {
	FileID       string `json:"file_id"`
	Path         string `json:"path"`
	PresignedURL string `json:"presigned_url"`
}
