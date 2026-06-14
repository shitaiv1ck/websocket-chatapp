package core_errors

import "errors"

var (
	ErrCoockie        = errors.New("invalid cookie")
	ErrInvalidArg     = errors.New("invalid argument")
	ErrConflict       = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
	ErrUnauthenticate = errors.New("invalid username or password")
	ErrUnauthorize    = errors.New("invalid user id")
)
