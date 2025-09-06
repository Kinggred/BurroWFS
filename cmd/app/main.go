package main

import (
	"burrowfs/api/middleware"
	"burrowfs/api/routes"
	"burrowfs/core/config"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	builtInMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.LoadVariables()
	router := chi.NewRouter()

	router.Use(builtInMiddleware.Logger)
	router.Use(builtInMiddleware.RealIP)
	router.Use(builtInMiddleware.Timeout(60 * 10))

	router.Mount("/status", routes.StatusRoutes())
	router.Mount("/auth", routes.AuthRoutes())
	router.Mount("/user", middleware.DigestAuthMiddleware(routes.UserRoutes()))

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
