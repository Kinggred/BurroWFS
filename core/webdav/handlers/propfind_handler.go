package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
)

// HandlePropfind processes a PROPFIND request and retrieves the list of files for the user.
func HandlePropfind(user *rest.UserResponse, path *utils.Path, depth string) ([]models.File, error) {
	logger := logging.Get("webdav/propfind")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logging.Get("webdav/propfind").Error("Failed to connect to database: ", err)
		return nil, err
	}

	logger.Info("PROPFIND request for user: ", user)

	files, err := models.GetUserFiles(dbConn, user.Id, path.Clean, depth)
	if err != nil {
		logger.Error("Failed to retrieve user files: ", err)
		return nil, err
	}

	return files, nil
}
