package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/logging"
	methods "burrowfs/core/webdav"
	"burrowfs/core/webdav/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FilesRoutes() http.Handler {
	logger := logging.Get("webdav/router")
	router := chi.NewRouter()

	router.Options("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schemas.JSONResponse(w, http.StatusOK, "OK")
	}))

	router.Method(methods.PROPFIND, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := handlers.HandlePropfind(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		response := webdav.ParseFilesToMultistatus(files)

		fromRoot := r.URL.Path == "/"
		schemas.MultistatusWebDavResponse(w, response, fromRoot)
	}))

	router.Method(http.MethodPut, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := handlers.HandlePut(r)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		logger.Debug("PUT response: %+v", file)

		schemas.JSONResponse(w, http.StatusCreated, file)
	}))

	router.Method(methods.MKCOL, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, path := handlers.HandleMkcol(w, r)

		schemas.MkcolWebDavResponse(w, status, path)
	}))

	router.Method(methods.MOVE, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := handlers.HandleMove(w, r)
		if status != 201 && status != 204 {
			common.HttpError(w, status, http.StatusText(status))
			return
		}
		schemas.MoveWebDavResponse(w, status)
	}))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		file := handlers.HandleGet(w, r)
		if file == nil {
			common.HttpError(w, 404, "File not Found")
			return
		}
		code := 200
		if file.Redirect {
			code = 307
		}
		schemas.FileWebDavResponse(w, code, file)
	})

	return router
}
