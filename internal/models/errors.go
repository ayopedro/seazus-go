package models

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrInternalServerError  = errors.New("Internal server error")
	ErrUserNotFound         = errors.New("user not found!")
	ErrInvalidCredentials   = errors.New("invalid credentials. email/password incorrect")
	ErrInvalidPayload       = errors.New("invalid payload")
	ErrInvalidAuthorization = errors.New("invalid authorization header")
	ErrInvalidToken         = errors.New("invalid/expired token")
	ErrAuthentication       = errors.New("authentication error")
	ErrUnauthorized         = errors.New("user is unauthorized")
)

var (
	Conflict = "duplicate key value violates unique constraint"
)
