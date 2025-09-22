package schemas

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/logging"
	"net/http"
	"strconv"
)

func FileWebDavResponse(w http.ResponseWriter, status int, file *webdav.FileResponse) {
	w.WriteHeader(status)
	if file.Redirect {
		redirectFileResponse(w, file)
		return
	}
	directFileResponse(w, file)
}
func directFileResponse(w http.ResponseWriter, file *webdav.FileResponse) {
	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	w.Header().Set("ETag", file.ETag)
	w.Header().Set("Last-Modified", file.LastModified)
	w.Header().Set("Accept-Ranges", "bytes")
	_, err := w.Write(file.Data)

	if err != nil {
		common.HttpError(w, 500, "Failed writing response")
		logging.Get("responses/direct").Error(err.Error())
	}
}

func redirectFileResponse(w http.ResponseWriter, file *webdav.FileResponse) {
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Location", file.DownloadURL)
}
