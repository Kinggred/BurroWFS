package models

import (
	database "burrowfs/db"
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `db:"id"`
	Email      string     `db:"email"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	ModifiedAt time.Time  `db:"modified_at"`
	LastLogin  *time.Time `db:"last_login"`
	IsEnabled  bool       `db:"is_enabled"`
}

func CreateUser(ctx context.Context, db *database.DB, email, name string) (uuid.UUID, error) {
	query := db.Builder.Insert("users").Columns("email", "name").Values(email, name).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = db.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

func DeleteUser(ctx context.Context, db *database.DB, id uuid.UUID) uuid.UUID {
	query := db.Builder.Delete("users").Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Println(err)
		return uuid.Nil
	}

	_, err = db.Pool.Exec(ctx, sql, args...)
	return id
}
