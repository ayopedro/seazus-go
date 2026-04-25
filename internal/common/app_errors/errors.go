package appErrors

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("incorrect email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists")
	ErrForbidden          = errors.New("forbidden")
	ErrTimeout            = errors.New("operation timed out")
)

func MapPostgresError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return ErrConflict
		case "23503":
			return ErrInvalidInput
		default:
			return ErrInternal
		}
	}
	return err
}
