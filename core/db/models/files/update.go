package files

import (
	"burrowfs/core/db"
	"burrowfs/core/logging"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"time"
)

// UpdateFile updates a files record.
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
		Set("lock_info", file.LockInfo).
		Where(squirrel.Eq{"id": file.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return err
}

func BatchUpdateFiles(db *db.DB, files []File) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger := logging.Get("BatchUpdateFiles")
	logger.Debug(fmt.Sprintf("BatchUpdateFiles called with %d files", len(files)))

	if len(files) == 0 {
		return nil
	}

	tx, err := db.Pool.Begin(dbCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(dbCtx)

	batch := &pgx.Batch{}
	for _, file := range files {
		sql, args, err := db.Builder.
			Update("files").
			Set("path", file.Path).
			Set("path_parent_id", file.PathParentID).
			Set("updated_at", file.UpdatedAt).
			Where(squirrel.Eq{"id": file.ID}).
			ToSql()
		if err != nil {
			return err
		}
		batch.Queue(sql, args...)
		logger.Debug(fmt.Sprintf("Updating files: %v", file))
	}

	if batch.Len() == 0 {
		return tx.Commit(dbCtx) // Nothing to do
	}

	br := tx.SendBatch(dbCtx, batch)
	for range files {
		if _, err := br.Exec(); err != nil {
			err := br.Close()
			if err != nil {
				return err
			}
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(dbCtx)
}

// UpdateFilePartial updates only specified fields of a files record by its ID.
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
