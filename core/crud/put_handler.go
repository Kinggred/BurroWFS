package crud

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	"burrowfs/core/webdav/handlers"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func HandlePut(user *types.InternalUser, filesToAdd *[]rest.FileInRequest, basePath *types.Path) (code int, newFiles []types.FileDTO) {
	logger := logging.Get("handlers/put")
	logger.Debug(fmt.Sprint("Putting files: ", len(*filesToAdd)))

	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, nil
	}

	var createdFiles []*models.File
	var status int

	for _, file := range *filesToAdd {
		var parentDir *models.File
		var fileID uuid.UUID
		path := basePath.MergePaths(*utils.RetrievePath(file.RelativePath, true))
		if path.IsFolder() {
			status, _ = handlers.HandleMkcol(*user, path)
		} else {
			// Verify all parent directories exist
			if !basePath.IsParentRoot() {
				parentDir, err = models.GetFileByPath(dbConn, user.Id, basePath.Clean)
				if err != nil {
					logger.Error("Parent directory does not exist: ", err)
					return http.StatusConflict, nil
				}
				if parentDir.FileID != nil {
					logger.Error("Parent path is not a directory")
					return http.StatusConflict, nil
				}
			}
			fileID = uuid.New()
			var pathParentID *uuid.UUID = nil
			if parentDir != nil {
				pathParentID = &parentDir.ID
			}
			_, mimetype, _ := utils.GetFileMetadata(file.GetReadCloser())

			fileToAdd := models.File{
				ID:           uuid.New(),
				FileID:       &fileID,
				OwnerID:      user.Id,
				PathParentID: pathParentID,
				ContentType:  mimetype.String(),
				Name:         path.Name,
				Path:         path.Clean,
				S3Key:        utils.GenerateKey(user, fileID.String(), 1),
				Size:         file.Size,
				Version:      1,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			}

			createdFiles = append(createdFiles, &fileToAdd)
			newFiles = append(newFiles, fileToAdd.ToDTO())
		}
	}
	_, err = models.CreateBatch(dbConn, createdFiles)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return status, newFiles
}
