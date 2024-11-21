package repository

import "errors"

var (
	ErrFiltersNotProvided = errors.New("filters were not provided...")
	ErrNotFound           = errors.New("no data was found...")
	ErrAlreadyExists      = errors.New("data already exists...")
)
