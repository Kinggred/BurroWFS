package routes

import (
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	"burrowfs/core/crud"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// FileRoutes Provides more Rest like access to files and directories
// It expects and returns JSON data
// Allows for simpler implementation on the client side for basic web views
func FileRoutes() http.Handler {
	logger := logging.Get("routes/files")
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		filePath := r.URL.Query().Get("filePath")
		depth := r.URL.Query().Get("depth")

		path := utils.RetrievePath(filePath, true)
		logger.Debug(depth + " " + filePath)

		if depth != "0" && depth != "1" && depth != "infinity" {
			http.Error(w, "invalid depth value", http.StatusBadRequest)
			return
		}

		files, code := crud.HandleGet(&user, path, depth)

		if code != http.StatusOK {
			http.Error(w, http.StatusText(code), code)
			return
		}

		rest.NewGetResponseSchema(files)

		schemas.JSONResponse(w, code, files)
	})
	return router
}
