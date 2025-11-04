package routes

import (
	"burrowfs/api/common"
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/config"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	methods "burrowfs/core/webdav"
	"burrowfs/core/webdav/handlers"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func WebDavRoutes() http.Handler {
	logger := logging.Get("webdav/router")
	router := chi.NewRouter()

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		schemas.JSONResponse(w, http.StatusOK, "OK")
	})

	router.Method(methods.PROPFIND, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		path := utils.RetrievePath(r.URL.Path, true)
		depth := r.Header.Get("Depth")
		if depth == "" {
			depth = "infinite"
		}
		files, err := handlers.HandlePropfind(&user, path, depth)
		if err != nil {
			common.HttpError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		fileDTOs := make([]types.FileDTO, 0, len(files))
		for _, f := range files {
			fileDTOs = append(fileDTOs, f.ToDTO())
		}
		response := webdav.ParseFilesToMultistatus(fileDTOs)

		fromRoot := r.URL.Path == "/"
		schemas.MultistatusWebDavResponse(w, response, fromRoot)
	}))

	router.Method(http.MethodPut, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}(r.Body)

		newPath := utils.RetrievePath(r.URL.Path, true)
		logger.Debug("PUT request for path: " + newPath.Clean + " by user: " + user.Name)

		status := handlers.HandlePut(r.Context(), &user, newPath, r.Body)

		if status != http.StatusCreated && status != http.StatusOK {
			common.HttpError(w, status, http.StatusText(status))
			return
		}

	}))

	router.Method(methods.MKCOL, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPath := utils.RetrievePath(r.URL.Path, true)
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		status, path := handlers.HandleMkcol(user, newPath)

		schemas.MkcolWebDavResponse(w, status, path)
	}))

	router.Method(methods.MOVE, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		destination := r.Header.Get("Destination")
		if destination == "" {
			common.HttpError(w, http.StatusBadRequest, "Destination header is required")
			return
		}
		currentPath := utils.RetrievePath(r.URL.Path, true)
		newPath := utils.RetrievePath(destination, false)
		overwrite := r.Header.Get("Overwrite") == "T"
		status := handlers.HandleMove(&user, currentPath, newPath, overwrite)
		if status != 201 && status != 204 {
			common.HttpError(w, status, http.StatusText(status))
			return
		}
		schemas.MoveWebDavResponse(w, status)
	}))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, 401, "Auth error")
			logger.Debug(err.Error())
			return
		}
		path := utils.RetrievePath(r.URL.Path, true)
		client := r.Header.Get("User-Agent")

		blockRedirect := utils.Contains(config.CONFIG.BlockRedirectList, client) || config.CONFIG.ForceDirect
		status, combined := handlers.HandleGet(r.Context(), &user, path, blockRedirect)

		if status != 200 {
			common.HttpError(w, status, http.StatusText(status))
			return
		}

		if combined == nil || combined.IsEmpty() {
			common.HttpError(w, 500, "Internal server error")
			logger.Error("Combined response is nil or empty")
			return
		}

		if combined.ReturnAsRedirect() {
			http.Redirect(w, r, combined.Address, http.StatusTemporaryRedirect)
			return
		}

		response := webdav.NewFileResponse(combined)
		schemas.FileWebDavResponse(w, status, &response)
	})

	router.Method(methods.LOCK, "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.RetrieveUser(r.Context())
		if err != nil {
			common.HttpError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		path := utils.RetrievePath(r.URL.Path, true)
		depth := r.Header.Get("Depth")
		if depth == "" {
			depth = "infinite"
		}
		timeout := r.Header.Get("Timeout")
		if timeout == "" {
			timeout = "Infinite"
		}
		var lockInfo webdav.LockInfo
		if err := common.ParseXML(r, &lockInfo); err != nil {
			common.HttpError(w, http.StatusBadRequest, "Invalid XML body")
			return
		}

		status, resp, lockTokenHeader := handlers.HandleLock(&user, path, depth, timeout, lockInfo)
		if resp == nil {
			common.HttpError(w, status, http.StatusText(status))
			return
		}

		if lockTokenHeader != "" {
			w.Header().Set("Lock-Token", lockTokenHeader)
		}
		schemas.LockWebDavResponse(w, status, *resp)
	}))

	return router
}
