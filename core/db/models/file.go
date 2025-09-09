package models

import (
	"burrowfs/core/db"
	"burrowfs/core/logging"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type File struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	FileID          *uuid.UUID `db:"file_id" json:"file_id"`                               // Logical file identifier for versioning
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
	Permissions     string     `db:"permissions" json:"permissions"` // JSON string for now
	LockInfo        string     `db:"lock_info" json:"lock_info"`     // JSON string for now
}

func nullableUUID(u uuid.UUID) interface{} {
	if u == uuid.Nil {
		return nil
	}
	return u
}

func CreateFile(db *db.DB, file *File) (uuid.UUID, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.Builder.Insert("files").
		Columns(
			"id", "file_id", "version_parent_id", "owner_id", "path_parent_id", "name", "path", "s3_key", "size", "etag", "version", "created_at", "updated_at", "permissions", "lock_info", "content_type",
		).
		Values(
			file.ID, file.FileID, file.VersionParentID, file.OwnerID, file.PathParentID, file.Name, file.Path, file.S3Key, file.Size, file.ETag, file.Version, file.CreatedAt, file.UpdatedAt, file.Permissions, file.LockInfo, file.ContentType,
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
			"created_at", "updated_at", "permissions", "lock_info", "content_type",
		)

	ids := make([]uuid.UUID, 0, len(files))

	for _, file := range files {
		query = query.Values(
			file.ID, file.FileID,
			file.VersionParentID,
			file.OwnerID,
			file.PathParentID,
			file.Name, file.Path, file.S3Key, file.Size, file.ETag, file.Version,
			file.CreatedAt, file.UpdatedAt, file.Permissions, file.LockInfo,
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

// GetUserFiles retrieves files and folders for a user starting from a given path.
// The depth parameter controls how deep to traverse:
// - "0": only the specified folder
// - "1": the folder and its immediate children
// - "infinity" or any other value: full recursive traversal
// Note: This function assumes that the startingPath exists and is a folder.
// If startingPath is "/", it retrieves from the root.
// If startingPath is not found, it returns an empty list.
func GetUserFiles(db *db.DB, ownerID uuid.UUID, startingPath string, depth string) ([]File, error) {
	startingPath = "/files" + startingPath
	logger := logging.Get("GETUserFiles")
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := ""
	var args []interface{}
	var err error

	logger.Debug("GetUserFiles called with ownerID: " + ownerID.String() + ", startingPath: " + startingPath + ", depth: " + depth)
	var parentID *uuid.UUID
	if startingPath != "/" && depth != "0" {
		query, args, err := db.Builder.Select("id").From("files").Where(squirrel.Eq{"owner_id": ownerID, "path": startingPath}).Limit(1).ToSql()
		if err != nil {
			return nil, err
		}
		err = db.Pool.QueryRow(dbCtx, query, args...).Scan(&parentID)
		if err != nil {
			parentID = nil
		}
	}

	switch depth {
	case "0":
		// Only folder itself
		query, args, err = db.Builder.Select("*").From("files").Where(squirrel.Eq{"owner_id": ownerID, "path": startingPath}).ToSql()
		if err != nil {
			return nil, err
		}
		break
	case "1":
		// Folder and its immediate children

		if parentID == nil {
			// If parentID is nil, we are looking for root-level items
			query, args, err = db.Builder.Select("*").From("files").Where(squirrel.Eq{"owner_id": ownerID.String()}).Where(squirrel.Or{squirrel.Eq{"path": startingPath},
				squirrel.Or{squirrel.Expr("path_parent_id IS NULL")}}).ToSql()
		} else {
			query, args, err = db.Builder.Select("*").From("files").Where(squirrel.Eq{"owner_id": ownerID.String()}).Where(squirrel.Or{squirrel.Eq{"path": startingPath},
				squirrel.Eq{"path_parent_id": parentID}}).ToSql()
		}
		break
	default:
		// Full recursive
		// Using a recursive CTE to get all files and folders under the starting path
		// Note: This assumes that the startingPath exists and is a folder
		// If startingPath is not root, we need to find its ID first
		var baseWhere string
		if parentID == nil {
			baseWhere = "path_parent_id IS NULL"
			args = []interface{}{ownerID}
		} else {
			baseWhere = "path_parent_id = $2"

			args = []interface{}{ownerID, *parentID}
		}

		query = fmt.Sprintf(`
			WITH RECURSIVE file_tree AS (
    		SELECT *, ARRAY[id] AS path_ids
    		FROM files WHERE owner_id = $1 AND %s
    		UNION ALL

    		SELECT f.*, ft.path_ids || f.id
    		FROM files f JOIN file_tree ft ON f.path_parent_id = ft.id
    		WHERE f.owner_id = $1
			)
			SELECT id, file_id, version_parent_id, owner_id, path_parent_id,
    		name, path, s3_key, size, content_type, etag, version,
    		created_at, updated_at, permissions, lock_info 
			FROM file_tree
			ORDER BY
    		path_ids,                  -- ensures parent-first traversal
    		(file_id IS NOT NULL) ASC;
			`, baseWhere)
		break
	}

	rows, err := db.Pool.Query(dbCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.FileID,
			&file.VersionParentID,
			&file.OwnerID,
			&file.PathParentID,
			&file.Name,
			&file.Path,
			&file.S3Key,
			&file.Size,
			&file.ContentType,
			&file.ETag,
			&file.Version,
			&file.CreatedAt,
			&file.UpdatedAt,
			&file.Permissions,
			&file.LockInfo,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func GetFileByID(db *db.DB, id uuid.UUID) (*File, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.Builder.Select("*").From("files").Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var file File
	err = db.Pool.QueryRow(dbCtx, sql, args...).Scan(
		&file.ID,
		&file.FileID,
		&file.VersionParentID,
		&file.OwnerID,
		&file.PathParentID,
		&file.Name,
		&file.Path,
		&file.S3Key,
		&file.Size,
		&file.ETag,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.Permissions,
		&file.LockInfo,
		&file.ContentType,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFileByPath retrieves the latest version of a file by its path and owner ID.
// Note: This function currently does not handle shared files.
func GetFileByPath(db *db.DB, ownerID uuid.UUID, path string) (*File, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Won't work for shared files, need to add shared logic later
	// Maybe a separate function GetFileByPathWithAccess(userID, path string) that checks both owner and shared access
	query := db.Builder.Select("*").From("files").Where(squirrel.Eq{"owner_id": ownerID, "path": path}).OrderBy("version DESC").Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var file File
	err = db.Pool.QueryRow(dbCtx, sql, args...).Scan(
		&file.ID,
		&file.FileID,
		&file.VersionParentID,
		&file.OwnerID,
		&file.PathParentID,
		&file.Name,
		&file.Path,
		&file.S3Key,
		&file.Size,
		&file.ETag,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.Permissions,
		&file.LockInfo,
		&file.ContentType,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func UpdateFile(db *db.DB, file *File) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file.UpdatedAt = time.Now()

	query := db.Builder.Update("files").
		Set("name", file.Name).
		Set("path", file.Path).
		Set("s3_key", file.S3Key).
		Set("size", file.Size).
		Set("etag", file.ETag).
		Set("version", file.Version).
		Set("updated_at", file.UpdatedAt).
		Set("permissions", file.Permissions).
		Set("lock_info", file.LockInfo).
		Where(squirrel.Eq{"id": file.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}

// UpdateFilePartial updates only specified fields of a file record by its ID.
func UpdateFilePartial(db *db.DB, id string, updates map[string]interface{}) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(updates) == 0 {
		return nil // nothing to update
	}

	updates["updated_at"] = time.Now()

	query := db.Builder.Update("files")
	for k, v := range updates {
		query = query.Set(k, v)
	}
	query = query.Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}

func DeleteFile(db *db.DB, id string) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.Builder.Delete("files").Where(squirrel.Eq{"id": id})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}
