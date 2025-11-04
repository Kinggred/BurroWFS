package utils

import (
	"burrowfs/core/types"
	"fmt"
)

func GenerateKey(user *types.InternalUser, fileId string, version int) string {
	return fmt.Sprintf("%s/%s/%d", user.Id, fileId, version)
}
