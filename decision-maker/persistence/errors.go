package persistence

import (
	"database/sql"
	"errors"
	"strings"
)

//TODO constant errors using function notation (removes need for a ParseDBError function)
// ErrFoo(dbType) error

const PostgresError = iota

// general

var ErrNotFound = errors.New("record not found")

// user record related

var ErrUsernameAlreadyTaken = errors.New("username already taken")

// ParseDBError will return a more generic error given an error from the database
func ParseDBError(typeOfError int, err error) error {
	switch typeOfError {
	case PostgresError:
		return parsePostgresError(err)
	}
	return err
}

func parsePostgresError(err error) error {
	if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"user_username_key\"") {
		return ErrUsernameAlreadyTaken
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
