package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/utils"
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

		var body schemas.RegisterSchema
		err = common.ParseBody(r, &body)
		if err != nil {
			common.HttpError(w, http.StatusUnprocessableEntity, "Bad data provided")
			return
		}

		userId, err := models.CreateUser(r.Context(), dbConn, body.Email, body.Password, body.Name)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		common.StandardizedResponse(w, http.StatusOK, userId)
	})

	router.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		var body schemas.RegisterSchema
		err := common.ParseBody(r, &body)
		if err != nil {
			common.HttpError(w, http.StatusUnprocessableEntity, "Bad data provided")
			return
		}

		dbConn, err := db.Open()
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		defer dbConn.Close()
		user, err := models.GetUserByEmail(r.Context(), dbConn, body.Email)
		if err != nil || user == nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if !utils.CheckPassword(body.Password, user.Password) {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		common.StandardizedResponse(w, http.StatusOK, user)
	})

	return router
}
