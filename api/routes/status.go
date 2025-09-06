package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
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

		response := schemas.StatusSchema{
			Status:   "ok",
			Database: dbStatus,
			Time:     time.DateTime,
			Debug:    config.CONFIG.Debug,
		}

		common.RestResponse(w, http.StatusOK, response)
	})

	return router
}
