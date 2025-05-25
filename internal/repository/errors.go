package repository

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExist    = errors.New("exist")
	ErrNotExist = errors.New("not exist")
)
