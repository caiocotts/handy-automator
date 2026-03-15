package persistence

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

// general

var ErrNotFound = errors.New("record not found")

// user record related

var ErrUsernameAlreadyTaken = errors.New("username already taken")

// ParseDBError will return a more generic error given an error from the database
func ParseDBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == uniqueViolation && pgErr.ConstraintName == "user_username_key" {
			return ErrUsernameAlreadyTaken
		}
	}

	return err
}
