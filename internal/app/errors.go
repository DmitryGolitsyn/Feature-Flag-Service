package app

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrInvalid   = errors.New("invalid_input")
	ErrForbidden = errors.New("forbidden")
	ErrInternal  = errors.New("internal_error")
)
