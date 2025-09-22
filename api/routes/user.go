package routes

import (
	"burrowfs/api/schemas"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UserRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		schemas.JSONResponse(w, http.StatusOK, r.Context().Value("user"))
	})

	return router
}
