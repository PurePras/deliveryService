// Package apperror defines sentinel errors shared by the repository, service,
// and handler layers so a repository-level DB error can be mapped all the way
// to an HTTP status code without those layers importing each other.
package apperror

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrValidation   = errors.New("validation failed")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)
