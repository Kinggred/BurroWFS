package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// createFile creates a new file record
func createFile(user rest.UserResponse, filepath string, fileData []byte, pathParentId *uuid.UUID) *models.File {
	fileId := uuid.New()
	fileVersion := 1
	newFile := models.File{
		ID:           uuid.New(),
		FileID:       &fileId,
		OwnerID:      user.Id,
		PathParentID: pathParentId,
		ContentType:  "text/yaml", // TODO: Detect content type
		Name:         filepath,
		Path:         filepath,
		S3Key:        fmt.Sprintf("/%s/%s/%d", user.Id, fileId, fileVersion),
		Size:         int64(len(fileData)),
		ETag:         "", // To be generated after S3 upload
		Version:      fileVersion,
		Permissions:  "{}", // Default permissions
		LockInfo:     "{}", // Figure out later
	}
	return &newFile
}

// HandlePut processes a PUT request to upload or update a file.
// It should not create folders leading to put, those should be created with MKCOL
// but if folders do not exist, it will create them automatically, for now
func HandlePut(r *http.Request) (*models.File, error) {
	var logger = logging.Get("webdav/put")
	user, err := utils.RetrieveUser(r.Context())
	if err != nil {
		logger.Error("Failed to retrieve user from context: ", err)
		return &models.File{}, err
	}
	logger.Info("PUT request for user: ", user)

	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
	}

	// TODO: File content waits for S3 integration
	fileData, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read request body: ", err)
		return &models.File{}, err
	}
	logger.Debug(string(fileData))

	filepath := r.URL.Path
	logger.Debug("File path: ", filepath)
	file, _ := models.GetFileByPath(dbConn, user.Id, filepath)

	if file == nil {
		var pathParentID *uuid.UUID = nil

		fileRoute := strings.Split(filepath, "/")[1:]
		if len(fileRoute) > 1 {
			var resources []*models.File
			for depth, filename := range fileRoute {
				if filename == "" {
					continue
				}
				logger.Debug(fmt.Sprintf("Depth %d: %s", depth, filename))
				if len(fileRoute)-1 == depth {
					resources = append(resources, createFile(user, filepath, fileData, pathParentID))
					break
				}
				// Everything except the last segment is a folder
				folderPath := "/" + strings.Join(fileRoute[:depth+1], "/")
				existingFolder, _ := models.GetFileByPath(dbConn, user.Id, folderPath) // TODO: File tree check can optimize DB queries, this will do for now
				folderID := uuid.New()
				if existingFolder == nil {
					folder := models.File{
						ID:           folderID,
						OwnerID:      user.Id,
						PathParentID: pathParentID,
						Name:         filename,
						Path:         folderPath,
						LockInfo:     "{}", // Figure out later
						Permissions:  "{}", // Default permissions
					}
					resources = append(resources, &folder)
					pathParentID = &folderID
				}
				if existingFolder != nil {
					pathParentID = &existingFolder.ID
				}
			}
			_, err := models.CreateBatch(dbConn, resources)
			if err != nil {
				logger.Error("Failed to create folder structure: ", err)
				return &models.File{}, err
			}
			logger.Info("Created folder structure for path: ", filepath)
			return resources[len(resources)-1], nil // Return the actual file
		}
		newFile := createFile(user, filepath, fileData, pathParentID)
		fileId, err := models.CreateFile(dbConn, newFile)
		if err != nil {
			logger.Error("Failed to create new file record: ", err)
			return &models.File{}, err
		}
		logger.Info("Created new file with ID: ", fileId)
		// File upload to S3 can be handled asynchronously
		// go uploadToS3(newFile.S3Key, fileData)
		return newFile, nil

	} else {
		// File exists, update logic can be implemented here
		logger.Info("File already exists with ID: ", file.FileID)
		logger.Error("NOT IMPLEMENTED")
		return &models.File{}, nil
	}
}
