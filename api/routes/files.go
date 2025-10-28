package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	"burrowfs/core/crud"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	"fmt"
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

		fileDTOs := make([]types.FileDTO, 0, len(files))
		for _, f := range files {
			fileDTOs = append(fileDTOs, f.ToDTO())
		}

		response := rest.NewGetResponseSchema(fileDTOs)

		schemas.JSONResponse(w, code, response)
	})

	router.Put("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var body rest.PutFilesInputSchema
		err = common.ParseJSON(r, &body, false)
		if err != nil {
			logger.Debug(fmt.Sprintf("%s: invalid request body", err))
			http.Error(w, "Bad data provided", http.StatusUnprocessableEntity)
			return
		}

		code := crud.HandlePut(&user, body)

		schemas.JSONResponse(w, code, nil)
	})

	return router
}
