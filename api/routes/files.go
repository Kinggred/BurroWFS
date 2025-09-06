package routes

import (
	"burrowfs/api/common"
	"burrowfs/core/webdav"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FilesRoutes() http.Handler {
	router := chi.NewRouter()

	router.Method(webdav.PROPFIND, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common.RestResponse(w, http.StatusNotImplemented, "PROPFIND not implemented yet")
	}))

	router.Method(webdav.MKCOL, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common.RestResponse(w, http.StatusNotImplemented, "MKCOL not implemented yet")
	}))

	return router
}
