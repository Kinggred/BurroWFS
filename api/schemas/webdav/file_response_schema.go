package webdav

import (
	"burrowfs/core/logging"
	"strings"
	"time"
)

type FileResponse struct {
	ContentType  string
	Size         int64
	ETag         string
	Redirect     bool
	LastModified string
	DownloadURL  string
	Data         []byte
}

func NewFileResponse(contentType string, eTag string, size int64, updatedAt time.Time, downloadURL string, data []byte) FileResponse {
	var redirect bool = true
	if strings.TrimSpace(downloadURL) == "" {
		if data == nil {
			logging.Get("webdav/file_response").Error("Both downloadURL and data are empty in FileResponse")
		}
		redirect = false
	}

	return FileResponse{
		ContentType:  contentType,
		Size:         size,
		ETag:         eTag,
		Redirect:     redirect,
		LastModified: updatedAt.Format(time.RFC1123),
		DownloadURL:  downloadURL,
		Data:         data,
	}
}
