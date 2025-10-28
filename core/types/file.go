package types

import (
	"time"

	"github.com/google/uuid"
)

// FileDTO is a shadow type for transferring file data between core and API layers.
// Not sure how to avoid this duplication without creating a cyclic dependency.
type FileDTO struct {
	ID           uuid.UUID  `json:"id"`
	FileID       *uuid.UUID `json:"file_id"`
	OwnerID      uuid.UUID  `json:"owner_id"`
	PathParentID *uuid.UUID `json:"parent_id"`
	ContentType  string     `json:"content_type"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	S3Key        string     `json:"s3_key"`
	Size         int64      `json:"size"`
	ETag         string     `json:"etag"`
	Version      int        `json:"version"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LockInfo     *string    `json:"lock_info"`
}
