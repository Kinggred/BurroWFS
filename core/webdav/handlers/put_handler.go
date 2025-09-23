package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/aws"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// HandlePut processes a PUT request to upload or update a file.
func HandlePut(user *rest.UserResponse, newPath *utils.Path, data io.ReadCloser) int {
	var logger = logging.Get("webdav/put")
	logger.Info("PUT request for user: ", user)
	contents, err := io.ReadAll(data)
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError
	}
	fileSize := int64(len(contents))

	mime, err := mimetype.DetectReader(data)
	if err != nil {
		logger.Error("Failed to detect MIME type: ", err)
		return http.StatusInternalServerError
	}
	logger.Debug("Detected MIME type: ", mime.String())

	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
	}

	if fileSize > aws.MaxFileSize {
		logger.Info("File size exceeds limit")
		return http.StatusRequestEntityTooLarge
	}

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
		file = &models.File{
			ID:           uuid.New(),
			FileID:       &fileId,
			OwnerID:      user.Id,
			PathParentID: pathParentID,
			ContentType:  mime.String(),
			Name:         newPath.Name,
			Path:         newPath.Clean,
			S3Key:        fmt.Sprintf("/%s/%s/%d", user.Id, fileId, fileVersion),
			Size:         fileSize,
			ETag:         "", // To be generated after S3 upload
			Version:      fileVersion,
			Permissions:  "{}", // Default permissions
			LockInfo:     "{}", // Figure out later
		}
		if err != nil {
			logger.Error("Failed to create new file record: ", err)
			return http.StatusInternalServerError
		}
		logger.Info("Created new file: ", newFile.Path)
	} else {
		logger.Error("File versioning not implemented yet")
		return http.StatusInternalServerError
		//newVersion := file.Version + 1
		// Update existing file record with new version
		// create a method saving previous version
	}

	eTag, err := aws.PutFile(file, data)
	if err != nil {
		return http.StatusInternalServerError
	}
	newFile.ETag = eTag

	_, err = models.CreateFile(dbConn, &newFile)

	return http.StatusCreated
}
