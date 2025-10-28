package types

import (
	"io"
)

// CombinedResponses holds the possible responses for a GET request.
// It has to contain a File metadata and either raw Data or a redirect Address.
type CombinedResponses struct {
	File    *FileDTO
	Data    io.ReadCloser
	Address string
}

func (c *CombinedResponses) IsEmpty() bool {
	return c.File == nil && c.Data == nil && c.Address == ""
}

func (c *CombinedResponses) ReturnAsRedirect() bool {
	return c.Address != ""
}
