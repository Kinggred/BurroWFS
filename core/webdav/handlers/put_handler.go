package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/aws"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func getFileMetadata(data io.ReadCloser) (newData io.ReadCloser, mime *mimetype.MIME, fileSize int64) {
	var logger = logging.Get("webdav/fileMetadata")
	buf, err := io.ReadAll(io.LimitReader(data, 4096))
	if err != nil {
		logger.Error("Failed to read data: ", err)
		return newData, nil, 0
	}
	mime = mimetype.Detect(buf)

	remainingData, err := io.ReadAll(data)
	if err != nil {
		logger.Error("Failed to read remaining data: ", err)
	}
	fullContents := append(buf, remainingData...)
	fileSize = int64(len(fullContents))
	data = io.NopCloser(bytes.NewReader(fullContents))

	logger.Debug("Detected MIME type: ", mime.String())
	return data, mime, fileSize
}

// HandlePut processes a PUT request to upload or update a file.
func HandlePut(ctx context.Context, user *rest.UserResponse, newPath *utils.Path, data io.ReadCloser) int {
	var logger = logging.Get("webdav/put")
	logger.Info("PUT request for user: ", user)

	data, mime, fileSize := getFileMetadata(data)
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
			S3Key:        fmt.Sprintf("%s/%s/%d", user.Id, fileId, fileVersion),
			Size:         fileSize,
			ETag:         "", // To be generated after S3 upload
			Version:      fileVersion,
		}
		logger.Info("Created new file: ", file.Path)
	} else {
		logger.Error("File versioning not implemented yet")
		return http.StatusInternalServerError
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
