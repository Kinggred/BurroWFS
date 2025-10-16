package handlers

import (
	"burrowfs/api/schemas"
	"burrowfs/api/schemas/rest"
	database "burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/utils"
)

func HandleLock(user *rest.UserResponse, path *utils.Path, depth string, timeout string, body schemas.LockInfo) int {
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

	err = models.LockFile(dbConn, resourceToLock, depth, timeout, body)
	if err != nil {
		return 500
	}

	return 200
}
