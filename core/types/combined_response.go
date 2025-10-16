package types

import (
	"burrowfs/core/db/models"
	"io"
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
