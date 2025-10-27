package handlers

import (
	"burrowfs/api/schemas"
	database "burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/types"
	"time"
)

func HandleLock(user *types.InternalUser, path *types.Path, depth string, timeout string, body schemas.LockInfo) int {
	dbConn, err := database.Open()
	defer dbConn.Close()
	if err != nil {
		return 500
	}
	resourceToLock, err := models.GetFileByPath(dbConn, user.Id, path.Clean)
	if err != nil || resourceToLock == nil {
		return 404
	}

	if resourceToLock.IsLocked() {
		return 423
	}

	err = models.LockFile(dbConn, *resourceToLock.FileID, user.Id, time.Minute*30) // TODO: use timeout from header
	if err != nil {
		return 500
	}

	return 200
}
