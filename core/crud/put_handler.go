package crud

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"burrowfs/core/utils"
	"burrowfs/core/webdav/handlers"
	"context"
	"fmt"
	"net/http"
)

func HandlePut(ctx context.Context, user *types.InternalUser, filesToAdd *[]rest.FileBodySchema) int {
	logger := logging.Get("handlers/put")
	logger.Debug(fmt.Sprint("Putting files: ", len(*filesToAdd)))
	var createdFilePaths []*types.Path

	for _, file := range *filesToAdd {
		path := utils.RetrievePath(file.RelativePath, true)
		status := handlers.HandlePut(ctx, user, path, file.GetReadCloser())
		if status != http.StatusCreated {
			logger.Error(fmt.Sprintf("Failed to put file %s with status %d", file.RelativePath, status))
			// TODO: Create proper rollback mechanism once delete handler is implemented
			// Rollback previously created files
			// for _, createdPath := range createdFilePaths {
			// handlers.HandleDelete(ctx, user, createdPath)
			// }
		} else {
			createdFilePaths = append(createdFilePaths, path)
		}
	}
	return http.StatusCreated
}
