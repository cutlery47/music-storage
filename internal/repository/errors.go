package repository

import "errors"

var (
	ErrFiltersNotProvided = errors.New("filters were not provided...")
	ErrSongNotFound       = errors.New("song was not found...")
)
