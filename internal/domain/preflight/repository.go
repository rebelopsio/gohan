package preflight

import (
	"context"
)

// ValidationSessionRepository defines persistence operations for validation sessions
type ValidationSessionRepository interface {
	// Save persists a validation session
	Save(ctx context.Context, session *ValidationSession) error

	// FindByID retrieves a session by ID
	FindByID(ctx context.Context, id string) (*ValidationSession, error)

	// FindLatest retrieves the most recent session
	FindLatest(ctx context.Context) (*ValidationSession, error)

	// List retrieves all sessions
	List(ctx context.Context) ([]*ValidationSession, error)
}
