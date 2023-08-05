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

func PGDuplicationError(err error) (isDuplicationErr bool) {
	if data, ok := err.(*pgconn.PgError); ok && data.Code == "23505" {
		isDuplicationErr = true
	}
	return
}
