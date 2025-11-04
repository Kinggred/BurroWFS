package handlers

import (
	"burrowfs/core/aws"
	"burrowfs/core/db"
	"burrowfs/core/db/models/files"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// HandleGet processes a GET request to retrieve a files's metadata and content.
func HandleGet(ctx context.Context, user *types.InternalUser, path *types.Path, blockRedirect bool) (status int, combinedResponse *types.CombinedResponses) {
	logger := logging.Get("handlers/get")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, nil
	}

	file, err := files.GetFileByPath(dbConn, user.Id, path.Clean)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return http.StatusNotFound, nil
		} else {
			logger.Error(err.Error())
			return http.StatusInternalServerError, nil
		}
	}

	fileDto := file.ToDTO()
	combinedResponse = utils.NewCombinedResponses(&fileDto, nil, "")
	if blockRedirect {
		stream, err := aws.GetFileStream(ctx, file.S3Key)
		if err != nil {
			logger.Error("Failed to get S3 files stream: ", err)
			return http.StatusInternalServerError, nil
		}
		combinedResponse.Data = stream
	} else {
		combinedResponse.Address, err = aws.GetFileURL(ctx, file.S3Key)
		if err != nil {
			logger.Error("Failed to get S3 files URL: ", err)
			return http.StatusInternalServerError, nil
		}
	}

	return http.StatusOK, combinedResponse
}
