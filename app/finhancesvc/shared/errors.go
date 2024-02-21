package shared

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrExist        = errors.New("already exist")
	ErrInternal     = errors.New("internal server error")
)

// This only used for demo purpose
// add comment for issue tracking in github

func PGDuplicationError(err error) (isDuplicationErr bool) {
	if data, ok := err.(*pgconn.PgError); ok && data.Code == "23505" {
		isDuplicationErr = true
	}

	return
}
