package models

import "errors"

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrDuplicateEmail      = errors.New("email already exists")
	ErrInternalServerError = errors.New("Internal server error")
	ErrUserNotFound        = errors.New("user not found!")
)
