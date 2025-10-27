package types

import "github.com/google/uuid"

// InternalUser defines only the user data used internally in the app.
type InternalUser struct {
	Id    uuid.UUID
	Email string
	Name  string
}
