package handlers

import (
	"burrowfs/api/common"
	schemas "burrowfs/api/schemas/webdav"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"net/http"
)

func HandleGet(w http.ResponseWriter, r *http.Request) *schemas.FileResponse {
	logger := logging.Get("handlers/get")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		common.HttpError(w, 500, "Error")
		logger.Debug(err.Error())
		return nil
	}
	filePath := r.URL.Path
	user, err := utils.RetrieveUser(r.Context())
	if err != nil {
		common.HttpError(w, 401, "Auth error")
		logger.Debug(err.Error())
		return nil
	}

	file, err := models.GetFileByPath(dbConn, user.Id, filePath)
	if err != nil {
		common.HttpError(w, 404, "Not Found")
		logger.Debug(err.Error())
		return nil
	}

	// Placeholder
	data := []byte(`
	some:
		ymltest
`)

	// TODO: AWS S3 Integration

	fileResponse := schemas.NewFileResponse(file.ContentType, file.ETag, file.Size, file.UpdatedAt, "", data)
	fileResponse.Size = int64(len(data))
	return &fileResponse
}
