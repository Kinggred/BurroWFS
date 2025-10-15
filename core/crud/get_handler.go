package crud

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"net/http"
)

func HandleGet(user *types.InternalUser, path *types.Path, depth string) ([]models.File, int) {
	logger := logging.Get("handlers/get")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error(err.Error())
		return nil, http.StatusInternalServerError
	}

	files, err := models.GetUserFiles(dbConn, user.Id, path.Clean, depth)
	if err != nil {
		logger.Error(err.Error())
		return nil, http.StatusInternalServerError
	}
	if len(files) == 0 {
		return nil, http.StatusNotFound
	}

	return files, http.StatusOK
}
