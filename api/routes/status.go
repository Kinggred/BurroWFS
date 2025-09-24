package routes

import (
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	"burrowfs/core/aws"
	"burrowfs/core/config"
	"burrowfs/core/db"
	"context"
	"net/http"
	"time"

	s4 "github.com/aws/aws-sdk-go-v2/service/s3"
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

		schemas.JSONResponse(w, http.StatusOK, response)
	})

	router.Get("/buckets", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s3 := aws.InitS3()
		if s3 == nil {
			schemas.JSONResponse(w, http.StatusOK, []string{})
			return
		}

		buckets, err := s3.ListBuckets(ctx, &s4.ListBucketsInput{})
		if err != nil {

			schemas.JSONResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		schemas.JSONResponse(w, http.StatusOK, buckets)
	})

	return router
}
