package main

import (
	"burrowfs/api/middleware"
	"burrowfs/api/routes"
	"burrowfs/core/config"
	"burrowfs/core/logging"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	builtInMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.LoadVariables()
	logger := logging.Get("main")
	router := chi.NewRouter()

	router.Use(middleware.RequestLoggerMiddleware)
	router.Use(builtInMiddleware.RealIP)
	router.Use(builtInMiddleware.Timeout(60 * 10))

	router.Mount("/status", routes.StatusRoutes())
	router.Mount("/auth", routes.AuthRoutes())
	router.Mount("/user", middleware.DigestAuthMiddleware(routes.UserRoutes()))

	logger.Info("listening on port: " + config.CONFIG.Port)
	logger.Debug("Debug mode is enabled")
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.CONFIG.Port), router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
