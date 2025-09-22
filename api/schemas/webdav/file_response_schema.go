package webdav

import (
	"burrowfs/core/webdav/handlers"
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

func NewFileResponse(responses *handlers.CombinedResponses) FileResponse {
	return FileResponse{
		ContentType:  responses.File.ContentType,
		Size:         responses.File.Size,
		ETag:         responses.File.ETag,
		Redirect:     responses.ReturnAsRedirect(),
		LastModified: responses.File.UpdatedAt.Format(time.RFC1123),
		DownloadURL:  responses.Address,
		Data:         *responses.Data,
	}
}
