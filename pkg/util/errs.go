package util

import "errors"

var (
	ErrChunkMissing = errors.New("chunk(s) missing")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
)
