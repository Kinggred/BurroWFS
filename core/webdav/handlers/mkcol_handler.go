package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"

	"github.com/google/uuid"
)

// HandleMkcol returns statusCode and file location
func HandleMkcol(user rest.UserResponse, newPath *utils.Path) (status int, location string) {
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
		dirParent, err := models.GetFileByPath(dbConn, user.Id, newPath.ParentPath)
		if err != nil {
			logger.Debug(err.Error())
			logger.Info("Parent folder not found")
			return 409, ""
		}
		pathParentID = &dirParent.ID
	}

	file := models.File{
		ID:           uuid.New(),
		OwnerID:      user.Id,
		PathParentID: pathParentID,
		Name:         newPath.Name,
		Path:         newPath.Clean,
		Size:         0,
		Permissions:  "{}",
		LockInfo:     "{}",
	}

	_, err = models.CreateFile(dbConn, &file)
	if err != nil {
		logger.Error(err.Error())
		return 500, ""
	}

	return 201, newPath.Clean

}
