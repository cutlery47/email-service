package service

import "errors"

var (
	ErrCacheNotFound = errors.New("cached data was not found")
	ErrWrongCode     = errors.New("provided code doesn't match the stored one")
)
