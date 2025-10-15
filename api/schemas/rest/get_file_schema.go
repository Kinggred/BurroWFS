package rest

import (
	"fmt"

	"github.com/google/uuid"
)

type GetFileSchema struct {
	FileID   uuid.UUID `json:"file_id"`
	FilePath string    `json:"file_path"`
	Depth    string    `json:"depth"` // 0 / 1 / infinity
}

func (g *GetFileSchema) Validate() error {
	if g.Depth != "0" && g.Depth != "1" && g.Depth != "infinity" {
		return fmt.Errorf("invalid depth value: %s", g.Depth)
	}

	if g.FileID == uuid.Nil && g.FilePath == "" {
		return fmt.Errorf("either file_id or file_path must be provided")
	}

	return nil
}
