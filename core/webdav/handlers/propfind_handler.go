package handlers

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models/files"
	"burrowfs/core/logging"
	"burrowfs/core/types"
)

// HandlePropfind processes a PROPFIND request and retrieves the list of files for the user.
func HandlePropfind(user *types.InternalUser, path *types.Path, depth string) ([]files.File, error) {
	logger := logging.Get("webdav/propfind")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logging.Get("webdav/propfind").Error("Failed to connect to database: ", err)
		return nil, err
	}

	logger.Info("PROPFIND request for user: ", user)

	userFiles, err := files.GetUserFiles(dbConn, user.Id, path.Clean, depth)
	if err != nil {
		logger.Error("Failed to retrieve user userFiles: ", err)
		return nil, err
	}

	return userFiles, nil
}
