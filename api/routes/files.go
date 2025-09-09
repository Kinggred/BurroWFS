package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/logging"
	methods "burrowfs/core/webdav"
	"burrowfs/core/webdav/handlers"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func FilesRoutes() http.Handler {
	logger := logging.Get("webdav/files")
	router := chi.NewRouter()

	router.Options("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schemas.RestResponse(w, http.StatusOK, "OK")
	}))

	router.Method(methods.PROPFIND, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := handlers.HandlePropfind(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		response := webdav.ParseFilesToMultistatus(files)

		fromRoot := strings.TrimPrefix(r.URL.Path, "/files") == "/"

		schemas.MultistatusWebDavResponse(w, response, fromRoot)
	}))

	router.Method(http.MethodPut, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := handlers.HandlePut(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		logger.Debug("PUT response: %+v", file)

		schemas.RestResponse(w, http.StatusCreated, file)
	}))

	router.Method(methods.MKCOL, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schemas.RestResponse(w, http.StatusNotImplemented, "MKCOL not implemented yet")
	}))

	return router
}
