package middleware

import (
	"burrowfs/core/logging"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	logger := logging.Get("RequestLoggerMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Info(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		logger.Debug(fmt.Sprintf("%s", r.Body))
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		logger.Info(fmt.Sprintf("%s: %d %s", "Request processed", ww.Status(), time.Since(startTime)))
	})
}
