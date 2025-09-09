package handlers

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"net/http"
	"strings"
)

func HandlePropfind(r *http.Request) ([]models.File, error) {
	depth := r.Header.Get("Depth")
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/files") // TODO: Configurable base path
	if path[len(path)-1] == '/' {
		path = strings.TrimSuffix(path, "/")
	}
	if path == "" {
		path = "/"
	}
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logging.Get("webdav/propfind").Error("Failed to connect to database: ", err)
		return nil, err
	}

	var logger = logging.Get("webdav/propfind")
	user, err := utils.RetrieveUser(r.Context())
	if err != nil {
		logger.Error("Failed to retrieve user from context: ", err)
		return nil, err
	}
	logger.Info("PROPFIND request for user: ", user)

	files, err := models.GetUserFiles(dbConn, user.Id, path, depth)
	if err != nil {
		logger.Error("Failed to retrieve user files: ", err)
		return nil, err
	}

	return files, nil
}
