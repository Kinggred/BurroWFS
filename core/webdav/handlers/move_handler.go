package handlers

import (
	"burrowfs/core/db"
	"burrowfs/core/db/models/files"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// HandleMove processes a MOVE request to move or rename a files or directory.
func HandleMove(user *types.InternalUser, currentPath *types.Path, newPath *types.Path, overwrite bool) int {
	logger := logging.Get("handlers/move")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Debug(err.Error())
		return 500
	}

	var newPathParentId *uuid.UUID = nil

	if currentPath.Clean == newPath.Clean {
		logger.Info("Tried moving folder into itself")
		return 403
	}

	newPathParent, err := files.GetFileByPath(dbConn, user.Id, newPath.ParentPath)
	if err != nil && newPath.ParentPath != "/" {
		logger.Debug(err.Error())
		return 409
	} else if newPathParent != nil {
		newPathParentId = &newPathParent.ID
	}

	destinationFile, err := files.GetFileByPath(dbConn, user.Id, newPath.Clean)
	if err == nil {
		if !overwrite {
			logger.Debug("File already exists at destination: ", destinationFile.Path)
			return 412
		}
	} else if errors.Is(err, pgx.ErrNoRows) {
		logger.Debug("File not found at destination, proceeding")
	} else {
		logger.Debug(err.Error())
		return 500
	}

	if currentPath.IsFolder() {
		if destinationFile != nil {
			deletedFilesIDs, err := files.DeleteDirectory(dbConn, user.Id, destinationFile)
			if err != nil {
				logger.Debug(err.Error())
				return 500
			}
			// Log deleted files IDs for further processing (e.g., S3 cleanup)
			// This is a placeholder for future S3 integration
			// where we would delete the files from S3 storage as well
			// For now, we just log the IDs
			for _, fileID := range deletedFilesIDs {
				logger.Debug("Deleted files ID: ", fileID)
				// TODO: S3 Integration
			}
		}
		err := files.MoveDirectory(dbConn, user.Id, currentPath, newPath, newPathParentId)
		if err != nil {
			logger.Error(err.Error())
			return 500
		}
		return 204

	} else {
		if destinationFile != nil {
			err = files.DeleteFileByID(dbConn, destinationFile.ID)
			if err != nil {
				logger.Debug(err.Error())
				return 412
			}
		}
		err = files.MoveFile(dbConn, user.Id, currentPath, newPath, newPathParent)
		if err != nil {
			logger.Debug(err.Error())
			return 500
		}
		return 201

	}
}
