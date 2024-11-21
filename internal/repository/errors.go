package repository

import "errors"

var (
	ErrNotFound      = errors.New("no data was found...")
	ErrAlreadyExists = errors.New("data already exists...")
)
