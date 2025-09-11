package utils

import (
	"burrowfs/api/schemas/rest"
	"context"
	"errors"
)

// RetrieveUser extracts the user information from the context.
// USER INFORMATION IS AVAILABLE ONLY IF THE REQUEST PASSED THROUGH AUTHENTICATION MIDDLEWARE.
func RetrieveUser(ctx context.Context) (rest.UserResponse, error) {
	user := ctx.Value("user")
	if user == nil {
		return rest.UserResponse{}, errors.New("user not found in context")
	}

	return user.(rest.UserResponse), nil
}
