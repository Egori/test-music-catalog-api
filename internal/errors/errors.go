package catalog_errors

import "errors"

var (
	ErrSongNotFound = errors.New("song not found")
	ErrInvalidPage  = errors.New("invalid page number")
)
