package models

import (
	"time"

	"github.com/google/uuid"
)

type FileVersion struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	FileID      *uuid.UUID `db:"file_id" json:"file_id"`     // links to main files
	FolderID    *uuid.UUID `db:"folder_id" json:"folder_id"` // optional, record folder at the time
	OwnerID     uuid.UUID  `db:"owner_id" json:"owner_id"`
	ContentType string     `db:"content_type" json:"content_type"`
	Name        string     `db:"name" json:"name"`
	S3Key       string     `db:"s3_key" json:"s3_key"`
	Size        int64      `db:"size" json:"size"`
	ETag        string     `db:"etag" json:"etag"`
	Version     int        `db:"version" json:"version"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}
