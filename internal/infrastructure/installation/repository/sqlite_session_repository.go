package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// SQLiteSessionRepository is a SQLite implementation of InstallationSessionRepository
type SQLiteSessionRepository struct {
	db *sql.DB
}

// NewSQLiteSessionRepository creates a new SQLite session repository
func NewSQLiteSessionRepository(dbPath string) (*SQLiteSessionRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteSessionRepository{db: db}
	if err := repo.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return repo, nil
}

// initialize creates the necessary tables
func (r *SQLiteSessionRepository) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS installation_sessions (
		id TEXT PRIMARY KEY,
		configuration TEXT NOT NULL,
		status TEXT NOT NULL,
		snapshot TEXT,
		installed_components TEXT,
		started_at DATETIME NOT NULL,
		completed_at DATETIME,
		failure_reason TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_status ON installation_sessions(status);
	CREATE INDEX IF NOT EXISTS idx_started_at ON installation_sessions(started_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// sessionDTO represents the database model
type sessionDTO struct {
	ID                  string
	Configuration       string
	Status              string
	Snapshot            sql.NullString
	InstalledComponents string
	StartedAt           time.Time
	CompletedAt         sql.NullTime
	FailureReason       sql.NullString
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// Save persists an installation session
func (r *SQLiteSessionRepository) Save(ctx context.Context, session *installation.InstallationSession) error {
	// Serialize configuration
	configData, err := json.Marshal(session.Configuration())
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Serialize snapshot if present
	var snapshotData sql.NullString
	if snapshot := session.Snapshot(); snapshot != nil {
		data, err := json.Marshal(snapshot)
		if err != nil {
			return fmt.Errorf("failed to marshal snapshot: %w", err)
		}
		snapshotData = sql.NullString{String: string(data), Valid: true}
	}

	// Serialize installed components
	componentsData, err := json.Marshal(session.InstalledComponents())
	if err != nil {
		return fmt.Errorf("failed to marshal components: %w", err)
	}

	// Handle completed_at
	var completedAt sql.NullTime
	if !session.CompletedAt().IsZero() {
		completedAt = sql.NullTime{Time: session.CompletedAt(), Valid: true}
	}

	// Handle failure reason
	var failureReason sql.NullString
	if session.FailureReason() != "" {
		failureReason = sql.NullString{String: session.FailureReason(), Valid: true}
	}

	now := time.Now()

	// Upsert query
	query := `
	INSERT INTO installation_sessions (
		id, configuration, status, snapshot, installed_components,
		started_at, completed_at, failure_reason, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		configuration = excluded.configuration,
		status = excluded.status,
		snapshot = excluded.snapshot,
		installed_components = excluded.installed_components,
		started_at = excluded.started_at,
		completed_at = excluded.completed_at,
		failure_reason = excluded.failure_reason,
		updated_at = excluded.updated_at
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		session.ID(),
		string(configData),
		string(session.Status()),
		snapshotData,
		string(componentsData),
		session.StartedAt(),
		completedAt,
		failureReason,
		now,
		now,
	)

	return err
}

// FindByID retrieves an installation session by its ID
func (r *SQLiteSessionRepository) FindByID(ctx context.Context, id string) (*installation.InstallationSession, error) {
	query := `
	SELECT id, configuration, status, snapshot, installed_components,
	       started_at, completed_at, failure_reason, created_at, updated_at
	FROM installation_sessions
	WHERE id = ?
	`

	var dto sessionDTO
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dto.ID,
		&dto.Configuration,
		&dto.Status,
		&dto.Snapshot,
		&dto.InstalledComponents,
		&dto.StartedAt,
		&dto.CompletedAt,
		&dto.FailureReason,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, installation.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	return r.dtoToSession(&dto)
}

// List retrieves all installation sessions
func (r *SQLiteSessionRepository) List(ctx context.Context) ([]*installation.InstallationSession, error) {
	query := `
	SELECT id, configuration, status, snapshot, installed_components,
	       started_at, completed_at, failure_reason, created_at, updated_at
	FROM installation_sessions
	ORDER BY started_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*installation.InstallationSession
	for rows.Next() {
		var dto sessionDTO
		err := rows.Scan(
			&dto.ID,
			&dto.Configuration,
			&dto.Status,
			&dto.Snapshot,
			&dto.InstalledComponents,
			&dto.StartedAt,
			&dto.CompletedAt,
			&dto.FailureReason,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		session, err := r.dtoToSession(&dto)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// dtoToSession converts a database DTO to a domain session
// Note: This is a simplified reconstruction. In a real system, you might need
// to use reflection or a more sophisticated approach to fully reconstruct the aggregate
func (r *SQLiteSessionRepository) dtoToSession(dto *sessionDTO) (*installation.InstallationSession, error) {
	// For now, we'll return an error indicating this needs proper implementation
	// In production, you'd need to properly deserialize and reconstruct the aggregate
	return nil, fmt.Errorf("session reconstruction from SQLite not yet fully implemented")
}

// Close closes the database connection
func (r *SQLiteSessionRepository) Close() error {
	return r.db.Close()
}
