package shared

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrExist        = errors.New("already exist")
	ErrInternal     = errors.New("internal server error")
)
