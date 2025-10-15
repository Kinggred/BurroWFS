package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas/rest"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// FileRoutes Provides more Rest like access to files and directories
// It expects and returns JSON data
// Allows for simpler implementation on the client side for basic web views
func FileRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var fileRequest rest.GetFileSchema

		err := common.ParseBody(r, &fileRequest, true)
		if err != nil {
			common.HttpError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	})
	return router
}
