package rest

import (
	"bytes"
	"io"
)

// PutFilesInputSchema defines the expected input for uploading a file.
type PutFilesInputSchema struct {
	Items []FileBodySchema `json:"items"`
}

type FileBodySchema struct {
	RelativePath string `json:"relative_path"`
	ContentType  string `json:"content_type"`
	Data         []byte `json:"data"`
}

func (f FileBodySchema) GetReadCloser() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(f.Data))
}
