package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes() http.Handler {
	router := chi.NewRouter()

	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		dbConn, err := db.Open()
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		defer dbConn.Close()

		var body rest.RegisterSchema
		err = common.ParseJSON(r, &body, false) // TODO
		if err != nil {
			common.HttpError(w, http.StatusUnprocessableEntity, "Bad data provided")
			return
		}

		userId, err := models.CreateUser(dbConn, body.Email, body.Password, body.Name)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		schemas.JSONResponse(w, http.StatusOK, userId)
	})

	return router
}
