package files

import (
	"burrowfs/core/db"
	"burrowfs/core/logging"
	"burrowfs/core/types"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"strings"
	"time"
)

func MoveFile(db *db.DB, ownerID uuid.UUID, oldPath *types.Path, newPath *types.Path, newPathParent *File) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newPathParentID *uuid.UUID = nil
	if newPathParent != nil {
		newPathParentID = &newPathParent.ID
	}

	query := db.Builder.Update("files").
		Set("path", newPath.Clean).Set("path_parent_id", newPathParentID).
		Where(squirrel.Eq{"owner_id": ownerID, "path": oldPath.Clean})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)

	return err
}

func MoveDirectory(db *db.DB, ownerID uuid.UUID, oldPath *types.Path, newPathRoot *types.Path, newPathParentID *uuid.UUID) error {
	logger := logging.Get("MoveDirectory")
	logger.Info("Moving directory from ", oldPath.Clean, " to ", newPathRoot.Clean)

	filesToMove, err := RecursiveFileSearch(db, ownerID, oldPath.Clean, false)
	if err != nil {
		return err
	}

	var filesToUpdate []File
	logger.Debug(fmt.Sprintf("Found %d files to move", len(filesToMove)))

	for index, file := range filesToMove {
		logger.Debug(fmt.Sprintf("Processing files %d: %v", index, file))
		if index == 0 {
			file.PathParentID = newPathParentID
		}
		file.Path = strings.Replace(file.Path, oldPath.Clean, newPathRoot.Clean, 1)
		filesToUpdate = append(filesToUpdate, file)
	}
	err = BatchUpdateFiles(db, filesToUpdate)
	if err != nil {
		return err
	}
	return nil
}

// LockFileWithToken sets a lock with token in lock_info as "expiry_userID_token".
func LockFileWithToken(db *db.DB, fileID uuid.UUID, userID uuid.UUID, token string, lockTimeout time.Duration) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lockInfo := fmt.Sprintf("%d_%s_%s", time.Now().Add(lockTimeout).Unix(), userID.String(), token)

	query := db.Builder.Update("files").
		Set("lock_info", lockInfo).
		Where(squirrel.Eq{"id": fileID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}

// UnlockFileWithToken clears lock_info only if the provided token matches the stored one.
func UnlockFileWithToken(db *db.DB, fileID uuid.UUID, token string) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch current lock_info
	var current string
	qsel, qargs, err := db.Builder.Select("lock_info").From("files").Where(squirrel.Eq{"id": fileID}).Limit(1).ToSql()
	if err != nil {
		return err
	}
	err = db.Pool.QueryRow(dbCtx, qsel, qargs...).Scan(&current)
	if err != nil {
		return err
	}
	_, _, storedToken, perr := ParseLockInfo(current)
	if perr != nil {
		return perr
	}
	if strings.TrimSpace(storedToken) != strings.TrimSpace(token) {
		return errors.New("unlock token does not match")
	}

	// Clear lock
	qupd, uargs, err := db.Builder.Update("files").Set("lock_info", "").Where(squirrel.Eq{"id": fileID}).ToSql()
	if err != nil {
		return err
	}
	_, err = db.Pool.Exec(dbCtx, qupd, uargs...)
	return err
}
