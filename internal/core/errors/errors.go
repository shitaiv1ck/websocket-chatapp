package core_errors

import "errors"

var (
	ErrCoockie    = errors.New("invalid coockie's files")
	ErrInvalidArg = errors.New("invalid argument")
	ErrConflict   = errors.New("already exists")
	ErrNotFound   = errors.New("not found")
)
