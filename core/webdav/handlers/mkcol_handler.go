package handlers

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// HandleMkcol returns statusCode and file location
func HandleMkcol(w http.ResponseWriter, r *http.Request) (int, string) {
	logger := logging.Get("webdav/mkcol")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		logger.Debug(err.Error())
		return 500, ""
	}

	user, err := utils.RetrieveUser(r.Context())
	if err != nil {
		logger.Debug(err.Error())
		return 401, ""
	}

	newDirPath := r.URL.Path
	dirParentSliced := strings.Split(newDirPath, "/")
	dirParentPath := strings.Join(dirParentSliced, "/")

	var pathParentID *uuid.UUID = nil
	var name = ""

	if len(dirParentSliced) < 3 {
		dirParent, err := models.GetFileByPath(dbConn, user.Id, dirParentPath[2:])
		if err != nil {
			logger.Debug(err.Error())
			logger.Info("File not found")
			return 404, ""
		}
		pathParentID = &dirParent.ID
		name = dirParentSliced[len(dirParentSliced)-2]
	} else {
		name = strings.TrimPrefix(strings.TrimSuffix(newDirPath, "/"), "/")
		pathParentID = nil
	}

	file := models.File{
		ID:           uuid.New(),
		OwnerID:      user.Id,
		PathParentID: pathParentID,
		Name:         name,
		Path:         strings.TrimSuffix(newDirPath, "/"),
		Size:         0,
		Permissions:  "{}",
		LockInfo:     "{}",
	}

	_, err = models.CreateFile(dbConn, &file)
	if err != nil {
		logger.Debug(err.Error())
		return 500, ""
	}

	return 201, newDirPath

}
