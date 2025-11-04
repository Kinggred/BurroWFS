package files

import (
	"burrowfs/core/db"
	"context"
	"github.com/google/uuid"
	"time"
)

func CreateFile(db *db.DB, file *File) (uuid.UUID, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.Builder.Insert("files").
		Columns(
			"id", "file_id", "version_parent_id", "owner_id", "path_parent_id", "name", "path", "s3_key", "size", "etag", "version", "created_at", "updated_at", "lock_info", "content_type",
		).
		Values(
			file.ID, file.FileID, file.VersionParentID, file.OwnerID, file.PathParentID, file.Name, file.Path, file.S3Key, file.Size, file.ETag, file.Version, file.CreatedAt, file.UpdatedAt, file.LockInfo, file.ContentType,
		)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return file.ID, err
}

func CreateBatch(db *db.DB, files []*File) ([]uuid.UUID, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(files) == 0 {
		return nil, nil
	}

	query := db.Builder.Insert("files").
		Columns(
			"id", "file_id", "version_parent_id", "owner_id", "path_parent_id",
			"name", "path", "s3_key", "size", "etag", "version",
			"created_at", "updated_at", "lock_info", "content_type",
		)

	ids := make([]uuid.UUID, 0, len(files))

	for _, file := range files {
		query = query.Values(
			file.ID, file.FileID,
			file.VersionParentID,
			file.OwnerID,
			file.PathParentID,
			file.Name, file.Path, file.S3Key, file.Size, file.ETag, file.Version,
			file.CreatedAt, file.UpdatedAt, file.LockInfo,
			file.ContentType,
		)
		ids = append(ids, file.ID)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
