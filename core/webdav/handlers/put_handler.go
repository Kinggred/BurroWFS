package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// HandlePut processes a PUT request to upload or update a file.
func HandlePut(user *rest.UserResponse, newPath *utils.Path, data *[]byte) int {
	fileSize := int64(len(*data))
	var logger = logging.Get("webdav/put")
	logger.Info("PUT request for user: ", user)

	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
	}

	// Limit file size to 10MB for now
	if fileSize > 10*1024*1024 {
		logger.Info("File size exceeds limit")
		return http.StatusRequestEntityTooLarge
	}

	// TODO: File content waits for S3 integration
	//	logger.Error("Failed to read request body: ", err)
	var pathParentID *uuid.UUID = nil
	var newFile models.File

	file, err := models.GetFileByPath(dbConn, user.Id, newPath.Clean)
	if err == nil {
	} else if errors.Is(err, pgx.ErrNoRows) {
		if newPath.ParentPath != "/" {
			dirParent, err := models.GetFileByPath(dbConn, user.Id, newPath.ParentPath)
			if err != nil {
				logger.Error("Parent folder not found: ", err)
				return http.StatusConflict
			}
			pathParentID = &dirParent.ID
		} else {
			pathParentID = nil
		}
	} else {
		logger.Error("Failed to check existing file: ", err)
		return http.StatusInternalServerError
	}

	if file == nil {
		fileId := uuid.New()
		fileVersion := 1
		newFile = models.File{
			ID:           uuid.New(),
			FileID:       &fileId,
			OwnerID:      user.Id,
			PathParentID: pathParentID,
			ContentType:  "text/yaml", // TODO: Detect content type
			Name:         newFile.Name,
			Path:         newPath.Clean,
			S3Key:        fmt.Sprintf("/%s/%s/%d", user.Id, fileId, fileVersion),
			Size:         fileSize,
			ETag:         "", // To be generated after S3 upload
			Version:      fileVersion,
			Permissions:  "{}", // Default permissions
			LockInfo:     "{}", // Figure out later
		}
		_, err = models.CreateFile(dbConn, &newFile)
		if err != nil {
			logger.Error("Failed to create new file record: ", err)
			return http.StatusInternalServerError
		}
		logger.Info("Created new file: ", newFile.Path)
		return http.StatusCreated
	} else {
		logger.Error("File versioning not implemented yet")
		return http.StatusInternalServerError
		//newVersion := file.Version + 1
		// Update existing file record with new version
		// create a method saving previous version
	}

}
