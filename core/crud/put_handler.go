package crud

import (
	"burrowfs/api/schemas/rest"
	"burrowfs/core/logging"
	"burrowfs/core/types"
)

func HandlePut(user *types.InternalUser, filesToAdd *[]rest.FileBodySchema) int {
	logger := logging.Get("handlers/put")
	logger.Error("Put handler not implemented")
	return 404
}
