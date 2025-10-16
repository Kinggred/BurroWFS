package utils

import (
	"burrowfs/core/types"
	"context"
	"errors"
)

// RetrieveUser extracts the user information from the context and converts it to InternalUser.
// USER INFORMATION IS AVAILABLE ONLY IF THE REQUEST PASSED THROUGH AUTHENTICATION MIDDLEWARE.
func RetrieveUser(ctx context.Context) (types.InternalUser, error) {
	userVal := ctx.Value("user")
	if userVal == nil {
		return types.InternalUser{}, errors.New("user not found in context")
	}

	internalUser, ok := userVal.(types.InternalUser)
	if !ok {
		return types.InternalUser{}, errors.New("user in context has invalid type")
	}

	return internalUser, nil
}
