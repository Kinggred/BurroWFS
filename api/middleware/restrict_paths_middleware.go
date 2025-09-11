package middleware

import (
	"burrowfs/core/config"
	"net/http"
	"strings"
)

// RestrictPathsMiddleware restricts access to certain paths based on configuration.
// It checks if the request path starts with the configured API prefix and denies access if it does.
// This prevents webdav clients from accessing API routes.
func RestrictPathsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.SplitAfter(r.URL.Path, "/")[1]
		if path == config.CONFIG.APIPrefix {
			http.Error(w, "Access to this path is restricted", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
