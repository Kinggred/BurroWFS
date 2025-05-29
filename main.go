package main

import (
	"burrowfs/api/routes"
	"burrowfs/config"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadVariables()
	router := chi.NewRouter()

	router.Mount("/status", routes.StatusRoutes())

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
