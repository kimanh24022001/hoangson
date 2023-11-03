// NOTE(duong): This file contains some of the postgres error aliases.
package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrCode_UniqueViolation = "23505"
)

func PgError(err error) *pgconn.PgError {
	var result *pgconn.PgError
	ok := errors.As(err, &result)

	if ok {
		return result
	}
	return nil
}
