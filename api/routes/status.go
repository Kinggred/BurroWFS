package routes

import (
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	"burrowfs/core/config"
	"burrowfs/core/db"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func StatusRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		dbConn, _ := db.Open()
		defer dbConn.Close()

		dbStatus := "ok"
		if dbConn == nil {
			dbStatus = "error"
		}

		response := rest.StatusSchema{
			Status:   "ok",
			Database: dbStatus,
			Time:     time.DateTime,
			Debug:    config.CONFIG.Debug,
		}

		schemas.RestResponse(w, http.StatusOK, response)
	})

	return router
}
