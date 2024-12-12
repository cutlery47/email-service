package service

import "errors"

var (
	ErrCacheNotFound = errors.New("cached data was not found")
)
