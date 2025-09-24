package main

import (
	"burrowfs/api/middleware"
	"burrowfs/api/routes"
	"burrowfs/core/config"
	"burrowfs/core/logging"
	"burrowfs/core/webdav"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	builtInMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.LoadVariables()
	logging.Init()
	registerWebDAVMethods()

	logger := logging.Get("main")
	router := chi.NewRouter()

	router.Use(middleware.RequestLoggerMiddleware)
	router.Use(builtInMiddleware.RealIP)
	router.Use(builtInMiddleware.Timeout(5 * time.Minute))

	router.Mount("/", webdavMiddlewares(routes.FilesRoutes()))
	router.Mount(pathWithAPIPrefix("/status"), routes.StatusRoutes())
	router.Mount(pathWithAPIPrefix("/auth"), routes.AuthRoutes())
	router.Mount(pathWithAPIPrefix("/user"), middleware.DigestAuthMiddleware(routes.UserRoutes()))

	logger.Info("listening on port: " + config.CONFIG.Port)
	logger.Debug("Debug mode is enabled")
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.CONFIG.Port), router)
	if err != nil {
		logger.Fatal("[*] " + err.Error())
	}
}

func registerWebDAVMethods() {
	chi.RegisterMethod(webdav.PROPFIND)
	chi.RegisterMethod(webdav.MKCOL)
	chi.RegisterMethod(webdav.COPY)
	chi.RegisterMethod(webdav.MOVE)
	chi.RegisterMethod(webdav.LOCK)
	chi.RegisterMethod(webdav.UNLOCK)
	chi.RegisterMethod(webdav.PROPPATCH)
}

func pathWithAPIPrefix(path string) string {
	return config.CONFIG.APIPrefix + path
}

func webdavMiddlewares(handler http.Handler) http.Handler {
	handler = middleware.DigestAuthMiddleware(handler)
	handler = middleware.RestrictPathsMiddleware(handler)
	return handler
}
