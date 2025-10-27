package repository

import (
	"context"
	"sync"

	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// MemorySessionRepository is an in-memory implementation of InstallationSessionRepository
// Useful for testing and development. In production, use a persistent storage implementation.
type MemorySessionRepository struct {
	sessions map[string]*installation.InstallationSession
	mu       sync.RWMutex
}

// NewMemorySessionRepository creates a new in-memory session repository
func NewMemorySessionRepository() *MemorySessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[string]*installation.InstallationSession),
	}
}

// Save persists an installation session
func (r *MemorySessionRepository) Save(ctx context.Context, session *installation.InstallationSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID()] = session
	return nil
}

// FindByID retrieves an installation session by its ID
func (r *MemorySessionRepository) FindByID(ctx context.Context, id string) (*installation.InstallationSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, installation.ErrSessionNotFound
	}

	return session, nil
}

// List retrieves all installation sessions
func (r *MemorySessionRepository) List(ctx context.Context) ([]*installation.InstallationSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]*installation.InstallationSession, 0, len(r.sessions))
	for _, session := range r.sessions {
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// Count returns the number of sessions stored (useful for testing)
func (r *MemorySessionRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.sessions)
}

// Clear removes all sessions (useful for testing)
func (r *MemorySessionRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions = make(map[string]*installation.InstallationSession)
}
