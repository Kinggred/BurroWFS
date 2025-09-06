package models

import (
	database "burrowfs/core/db"
	"burrowfs/core/utils"
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `db:"id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	ModifiedAt time.Time  `db:"modified_at"`
	LastLogin  *time.Time `db:"last_login"`
	IsEnabled  bool       `db:"is_enabled"`
}

func CreateUser(db *database.DB, email, password, name string) (uuid.UUID, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pwdHash, err := utils.EncodePassword(password)
	if err != nil {
		return uuid.Nil, err
	}

	query := db.Builder.Insert("users").Columns("email", "password", "name").Values(email, pwdHash, name).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = db.Pool.QueryRow(dbCtx, sql, args...).Scan(&id)
	return id, err
}

func GetUserByEmail(db *database.DB, email string) (*User, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := db.Builder.Select("*").From("users").Where(squirrel.Eq{"email": email})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = db.Pool.QueryRow(dbCtx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.LastLogin,
		&user.IsEnabled,
	)
	return &user, err
}

func DeleteUser(db *database.DB, id uuid.UUID) uuid.UUID {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := db.Builder.Delete("users").Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Println(err)
		return uuid.Nil
	}

	_, err = db.Pool.Exec(dbCtx, sql, args...)
	return id
}
