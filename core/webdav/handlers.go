package webdav

import (
	"burrowfs/core/logging"
	"net/http"
)

func handlePropfind(r *http.Request) {
	var logger = logging.Get("webdav/propfind")
	user := r.Context().Value("user")
	logger.Info("PROPFIND request for user: ", user)
	//Postgress lookup for user's files and folders
	//Return XML response with files and folders
}
