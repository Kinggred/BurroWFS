package handlers

import (
	"burrowfs/core/aws"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// HandlePut processes a PUT request to upload or update a file.
func HandlePut(ctx context.Context, user *types.InternalUser, newPath *types.Path, data io.ReadCloser) int {
	var logger = logging.Get("webdav/put")
	logger.Info("PUT request for user: ", user)

	data, mime, fileSize := utils.GetFileMetadata(data)
	if mime == nil {
		logger.Error("Failed to detect MIME type")
		return http.StatusInternalServerError
	}

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
			S3Key:        utils.GenerateKey(user, fileId.String(), fileVersion),
			Size:         fileSize,
			ETag:         "", // To be generated after S3 upload
			Version:      fileVersion,
		}
		logger.Info("Created new file: ", file.Path)
	} else {
		logger.Error("File versioning not implemented yet")
		return http.StatusInternalServerError
		// TODO: Implement versioning
		//newVersion := file.Version + 1
		// Update existing file record with new version
		// create a method saving previous version
	}

	eTag, err := aws.PutFile(ctx, file, data)
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError
	}
	// Update the ETag after successful upload via lambda events possibly
	file.ETag = eTag

	_, err = models.CreateFile(dbConn, file)

	return http.StatusCreated
}
