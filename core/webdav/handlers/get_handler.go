package handlers

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/aws"
	"burrowfs/core/db"
	"burrowfs/core/db/models"
	"burrowfs/core/logging"
	"burrowfs/core/utils"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// CombinedResponses holds the possible responses for a GET request.
// It hat to contain a File metadata and either raw Data or a redirect Address.
type CombinedResponses struct {
	File    *models.File
	Data    io.ReadCloser
	Address string
}

func (c *CombinedResponses) IsEmpty() bool {
	return c.File == nil && c.Data == nil && c.Address == ""
}

func (c *CombinedResponses) ReturnAsRedirect() bool {
	return c.Address != ""
}

func NewCombinedResponses(file *models.File, data io.ReadCloser, address string) *CombinedResponses {
	return &CombinedResponses{
		File:    file,
		Data:    data,
		Address: address,
	}
}

// HandleGet processes a GET request to retrieve a file's metadata and content.
func HandleGet(ctx context.Context, user *rest.UserResponse, path *utils.Path, blockRedirect bool) (status int, combinedResponse *CombinedResponses) {
	logger := logging.Get("handlers/get")
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		logger.Error(err.Error())
		return http.StatusInternalServerError, nil
	}

	file, err := models.GetFileByPath(dbConn, user.Id, path.Clean)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return http.StatusNotFound, nil
		} else {
			logger.Error(err.Error())
			return http.StatusInternalServerError, nil
		}
	}

	combinedResponse = NewCombinedResponses(file, nil, "")
	if blockRedirect {
		stream, err := aws.GetFileStream(ctx, file.S3Key)
		if err != nil {
			logger.Error("Failed to get S3 file stream: ", err)
			return http.StatusInternalServerError, nil
		}
		combinedResponse.Data = stream
	} else {
		combinedResponse.Address, err = aws.GetFileURL(ctx, file.S3Key)
		if err != nil {
			logger.Error("Failed to get S3 file URL: ", err)
			return http.StatusInternalServerError, nil
		}
	}

	return http.StatusOK, combinedResponse
}
