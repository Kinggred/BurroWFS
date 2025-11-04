package files

import (
	"burrowfs/core/db"
	"burrowfs/core/logging"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

// RecursiveFileSearch performs a recursive search for files and folders starting from a given path.
// If startDeeper is true, it starts from the parent of starting path; otherwise, it includes the starting path itself.
// If startingPath is "/", it retrieves all files and folders for the user.
// If startingPath is not found, it returns an empty list.
func RecursiveFileSearch(db *db.DB, ownerID uuid.UUID, startingPath string, startDeeper bool) ([]File, error) {
	logger := logging.Get("RecursiveFileSearch")
	logger.Debug("Starting recursive files search")
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := ""

	var startingPathParentID *uuid.UUID = nil
	var startingFile *File
	var err error
	if startingPath != "/" {
		startingFile, err = GetFileByPath(db, ownerID, startingPath)
		if err != nil || startingFile == nil {
			logger.Debug("Starting path not found: " + startingPath)
			return nil, nil
		}
		if startDeeper {
			startingPathParentID = startingFile.PathParentID
		} else {
			startingPathParentID = &startingFile.ID
		}
	}
	logger.Debug(fmt.Sprintf("Starting from files: %v, startingPathParentID: %s, startingFileID: %s", startingFile, startingPathParentID, startingFile.ID.String()))

	query = `
			WITH RECURSIVE file_tree AS (
    		SELECT *, ARRAY[id] AS path_ids
    		FROM files WHERE owner_id = $1 AND id = $2
    		UNION ALL

    		SELECT f.*, ft.path_ids || f.id
    		FROM files f JOIN file_tree ft ON f.path_parent_id = ft.id
    		WHERE f.owner_id = $1
			)
			SELECT id, file_id, version_parent_id, owner_id, path_parent_id,
    		name, path, s3_key, size, content_type, etag, version,
    		created_at, updated_at, lock_info 
			FROM file_tree
			ORDER BY
    		path_ids,                  -- ensures parent-first traversal
    		(file_id IS NOT NULL) ASC;
			`

	rows, err := db.Pool.Query(dbCtx, query, ownerID, startingPathParentID)
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
			&file.LockInfo,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	for index, file := range files {
		logger.Debug(fmt.Sprintf("%d: %v", index, file))
	}
	return files, nil
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
    		created_at, updated_at, lock_info 
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
		&file.ContentType,
		&file.ETag,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.LockInfo,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFileByPath retrieves the latest version of a files by its path and owner ID.
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
		&file.ContentType,
		&file.ETag,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.LockInfo,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
