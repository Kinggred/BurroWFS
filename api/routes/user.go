package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/middleware"
	"burrowfs/api/schemas"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func UserRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Digest ") {
			common.HttpError(w, http.StatusUnauthorized, "Unsupported auth schema")
			return
		}

		digestFields := middleware.ParseDigestHeader(authHeader[7:])
		username := digestFields["username"]

		if username == "" {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		doConn, err := db.Open()
		defer doConn.Close()
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := models.GetUserByEmail(doConn, username)
		if err != nil || user == nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userResponse := schemas.UserResponse{
			UserId:    user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.ModifiedAt.Format(time.RFC3339),
		}

		common.StandardizedResponse(w, http.StatusOK, userResponse)
	})

	return router
}
