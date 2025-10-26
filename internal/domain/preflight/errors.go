package preflight

import "errors"

var (
	// Entity/Value Object errors
	ErrInvalidDebianVersion = errors.New("invalid debian version")
	ErrInvalidGPU           = errors.New("invalid gpu configuration")
	ErrInvalidDiskSpace     = errors.New("invalid disk space value")

	// Repository errors
	ErrSessionNotFound = errors.New("validation session not found")
)
