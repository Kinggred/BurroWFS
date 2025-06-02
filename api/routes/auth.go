package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes() http.Handler {
	router := chi.NewRouter()

	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		dbConn, err := db.Open()
		if err != nil {
			// Handle
			panic(err)
		}
		defer dbConn.Close()
		var body schemas.RegisterSchema
		common.ParseBody(r, &body)

		_, err = models.CreateUser(r.Context(), dbConn, body.Email, body.Password, body.Name)
		if err != nil {
			// Handle
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		return
	})

	router.Get("/self", func(w http.ResponseWriter, r *http.Request) {
		var body schemas.RegisterSchema
		common.ParseBody(r, &body)

		dbConn, err := db.Open()
		if err != nil {
			// TODO: Handle
			panic(err)
		}
		defer dbConn.Close()
		user, err := models.GetUserByEmail(r.Context(), dbConn, body.Email)
		if err != nil {
			panic(err)
		}

		if !utils.CheckPassword(body.Password, user.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		usrStr, err := json.Marshal(user)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(usrStr)
		if err != nil {
			return
		}
	})

	return router
}
