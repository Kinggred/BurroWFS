package files

import (
	"burrowfs/core/db"
	"burrowfs/core/types"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

func DeleteDirectory(db *db.DB, ownerID uuid.UUID, directoryToDelete *File) ([]uuid.UUID, error) {
	filesToDelete, err := RecursiveFileSearch(db, ownerID, directoryToDelete.Path, false)
	if err != nil {
		return nil, err
	}

	var deletedIDs []uuid.UUID
	for _, file := range filesToDelete {
		deletedIDs = append(deletedIDs, file.ID)
	}

	err = DeleteFilesByIDs(db, deletedIDs)
	if err != nil {
		return nil, err
	}
	return deletedIDs, nil
}

func DeleteFileByID(db *db.DB, id uuid.UUID) error {
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

func DeleteFilesByIDs(db *db.DB, ids []uuid.UUID) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(ids) == 0 {
		return nil
	}

	query := db.Builder.Delete("files").Where(squirrel.Eq{"id": ids})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}

func DeleteFileByPath(db *db.DB, ownerID uuid.UUID, path *types.Path) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.Builder.Delete("files").Where(squirrel.Eq{"owner_id": ownerID, "path": path.Clean})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}
