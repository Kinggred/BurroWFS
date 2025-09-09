package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/mappers"
	webdav2 "burrowfs/api/schemas"
	"burrowfs/core/logging"
	"burrowfs/core/webdav"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func FilesRoutes() http.Handler {
	logger := logging.Get("webdav/files")
	router := chi.NewRouter()

	router.Options("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webdav2.RestResponse(w, http.StatusOK, "OK")
	}))

	router.Method(webdav.PROPFIND, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := webdav.HandlePropfind(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		response := mappers.FilesToWebDAVResponses(files)

		fromRoot := strings.TrimPrefix(r.URL.Path, "/files") == "/"

		webdav2.MultistatusWebDavResponse(w, response, fromRoot)
	}))

	router.Method(http.MethodPut, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := webdav.HandlePut(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		logger.Debug("PUT response: %+v", file)

		webdav2.RestResponse(w, http.StatusCreated, file)
	}))

	router.Method(webdav.MKCOL, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webdav2.RestResponse(w, http.StatusNotImplemented, "MKCOL not implemented yet")
	}))

	return router
}
