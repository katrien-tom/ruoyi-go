package errors

import "errors"

// Common domain errors
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
	ErrConflict      = errors.New("resource conflict")
	ErrInternal      = errors.New("internal error")
	ErrUnavailable   = errors.New("service unavailable")
	ErrTimeout       = errors.New("operation timeout")
)
