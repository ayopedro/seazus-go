package apperrors

import (
	"errors"
)

var (
	ErrInternal        = errors.New("Something went wrong")
	ErrURLNotFound     = errors.New("URL not found")
	ErrUserNotFound    = errors.New("User not found")
	ErrIdentifierTaken = errors.New("Identifier already in use")
	ErrEmailConflict   = errors.New("Email already in use")
	ErrInvalidURL      = errors.New("Invalid URL provided")
	ErrUnauthorized    = errors.New("Authentication required")
	ErrForbidden       = errors.New("You cannot access this resource")
	ErrInvalidPayload  = errors.New("Invalid payload")
)
