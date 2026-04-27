package appErrors

import (
	"database/sql"
	"errors"
	"net/http"

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

	ErrCreatingShortURL = errors.New("unable to shorten URL")
)

func MapPostgresError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

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

func StatusCode(err error) int {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrTimeout):
		return http.StatusRequestTimeout
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false

}
