package files

import (
	"burrowfs/core/types"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	FileID          *uuid.UUID `db:"file_id" json:"file_id"`                               // Logical files identifier for versioning
	VersionParentID *uuid.UUID `db:"version_parent_id" json:"version_parent_id,omitempty"` // Previous version (nullable for first version)
	OwnerID         uuid.UUID  `db:"owner_id" json:"owner_id"`
	PathParentID    *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"` // Parent folder ID (nullable for root)
	ContentType     string     `db:"content_type" json:"content_type"`
	Name            string     `db:"name" json:"name"`
	Path            string     `db:"path" json:"path"`
	S3Key           string     `db:"s3_key" json:"s3_key"`
	Size            int64      `db:"size" json:"size"`
	ETag            string     `db:"etag" json:"etag"`
	Version         int        `db:"version" json:"version"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	LockInfo        *string    `db:"lock_info" json:"lock_info"`
}

func (f File) String() string {
	return fmt.Sprintf("fileName: %s, filePath: %s", f.Name, f.Path)
}

// ToDTO converts a File model to a FileDTO for use in API or transfer layers.
func (f File) ToDTO() types.FileDTO {
	return types.FileDTO{
		ID:           f.ID,
		FileID:       f.FileID,
		OwnerID:      f.OwnerID,
		PathParentID: f.PathParentID,
		ContentType:  f.ContentType,
		Name:         f.Name,
		Path:         f.Path,
		S3Key:        f.S3Key,
		Size:         f.Size,
		ETag:         f.ETag,
		Version:      f.Version,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
		LockInfo:     f.LockInfo,
	}
}

// IsLocked checks if the files is currently locked based on the LockInfo field.
// Uses format: "expiry_userID_token".
func (f File) IsLocked() bool {
	if f.LockInfo == nil || strings.TrimSpace(*f.LockInfo) == "" {
		return false
	}
	expiry, _, _, err := ParseLockInfo(*f.LockInfo)
	if err != nil {
		return true
	}
	return time.Now().Before(expiry)
}

// SetLockWithToken sets lock with a token stored in lock_info as "expiry_userID_token".
func (f File) SetLockWithToken(userID uuid.UUID, token string, lockTimeout time.Duration) {
	lockInfo := fmt.Sprintf("%d_%s_%s", time.Now().Add(lockTimeout).Unix(), userID.String(), token)
	f.LockInfo = &lockInfo
}

// ParseLockInfo parses lock_info string into expiry time, userID and optional token.
func ParseLockInfo(s string) (expiry time.Time, userID uuid.UUID, token string, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, uuid.Nil, "", errors.New("empty lock_info")
	}
	parts := strings.Split(s, "_")
	if len(parts) < 2 {
		return time.Time{}, uuid.Nil, "", errors.New("invalid lock_info format")
	}
	expiryUnix, perr := strconv.ParseInt(parts[0], 10, 64)
	if perr != nil {
		return time.Time{}, uuid.Nil, "", perr
	}
	uid, uerr := uuid.Parse(parts[1])
	if uerr != nil {
		return time.Time{}, uuid.Nil, "", uerr
	}
	if len(parts) >= 3 {
		token = strings.Join(parts[2:], "_") // in case token contains underscores
	}
	return time.Unix(expiryUnix, 0), uid, token, nil
}
