package handlers

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models/files"
	"burrowfs/core/logging"
	"burrowfs/core/types"

	"github.com/google/uuid"
)

// HandleMkcol returns statusCode and files location
func HandleMkcol(user types.InternalUser, newPath *types.Path) (status int, location string) {
	logger := logging.Get("webdav/mkcol")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error(err.Error())
		return 500, ""
	}

	var pathParentID *uuid.UUID = nil

	if newPath.Clean == "/" {
		logger.Info("Tried creating root folder")
		return 403, ""
	} else if newPath.ParentPath == "/" {
		pathParentID = nil
	} else {
		dirParent, err := files.GetFileByPath(dbConn, user.Id, newPath.ParentPath)
		if err != nil {
			logger.Debug(err.Error())
			logger.Info("Parent folder not found")
			return 409, ""
		}
		pathParentID = &dirParent.ID
	}

	file := files.File{
		ID:           uuid.New(),
		OwnerID:      user.Id,
		PathParentID: pathParentID,
		Name:         newPath.Name,
		Path:         newPath.Clean,
		Size:         0,
	}

	_, err = files.CreateFile(dbConn, &file)
	if err != nil {
		logger.Error(err.Error())
		return 500, ""
	}

	return 201, newPath.Clean

}
