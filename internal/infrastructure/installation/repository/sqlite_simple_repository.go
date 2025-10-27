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

// SQLiteSimpleSessionRepository is a simple SQLite implementation that stores sessions as JSON
type SQLiteSimpleSessionRepository struct {
	db *sql.DB
}

// NewSQLiteSimpleSessionRepository creates a new simple SQLite session repository
func NewSQLiteSimpleSessionRepository(dbPath string) (*SQLiteSimpleSessionRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	_, err = db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	repo := &SQLiteSimpleSessionRepository{db: db}
	if err := repo.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return repo, nil
}

// initialize creates the necessary tables
func (r *SQLiteSimpleSessionRepository) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		data TEXT NOT NULL,
		started_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
	CREATE INDEX IF NOT EXISTS idx_sessions_started_at ON sessions(started_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// sessionRecord represents a row in the database
type sessionRecord struct {
	ID        string
	Status    string
	Data      string
	StartedAt time.Time
	UpdatedAt time.Time
}

// sessionStorageModel is a serializable representation of a session
// Uses DTOs with public fields for JSON serialization
type sessionStorageModel struct {
	ID                  string                     `json:"id"`
	Configuration       configurationDTO           `json:"configuration"`
	Status              string                     `json:"status"`
	Snapshot            *snapshotDTO               `json:"snapshot"`
	InstalledComponents []installedComponentDTO    `json:"installed_components"`
	StartedAt           time.Time                  `json:"started_at"`
	CompletedAt         time.Time                  `json:"completed_at"`
	FailureReason       string                     `json:"failure_reason"`
}

// configurationDTO is a serializable version of InstallationConfiguration
type configurationDTO struct {
	Components         []componentSelectionDTO `json:"components"`
	GPUVendor          string                  `json:"gpu_vendor,omitempty"`
	GPURequiresDriver  bool                    `json:"gpu_requires_driver,omitempty"`
	GPUDriverComponent string                  `json:"gpu_driver_component,omitempty"`
	DiskAvailable      uint64                  `json:"disk_available"`
	DiskRequired       uint64                  `json:"disk_required"`
	MergeExistingConf  bool                    `json:"merge_existing_conf"`
}

// componentSelectionDTO is a serializable version of ComponentSelection
type componentSelectionDTO struct {
	Component    string   `json:"component"`
	Version      string   `json:"version"`
	PackageName  string   `json:"package_name,omitempty"`
	SizeBytes    uint64   `json:"size_bytes,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// snapshotDTO is a serializable version of SystemSnapshot
type snapshotDTO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Path      string    `json:"path"`
	DiskAvailable uint64 `json:"disk_available"`
	DiskRequired  uint64 `json:"disk_required"`
	Packages  []string  `json:"packages"`
	Corrupted bool      `json:"corrupted"`
	CorruptionReason string `json:"corruption_reason,omitempty"`
}

// installedComponentDTO is a serializable version of InstalledComponent
type installedComponentDTO struct {
	ID           string    `json:"id"`
	Component    string    `json:"component"`
	Version      string    `json:"version"`
	PackageName  string    `json:"package_name,omitempty"`
	SizeBytes    uint64    `json:"size_bytes,omitempty"`
	Dependencies []string  `json:"dependencies,omitempty"`
	InstalledAt  time.Time `json:"installed_at"`
	Verified     bool      `json:"verified"`
	VerifiedAt   time.Time `json:"verified_at,omitempty"`
}

// toStorageModel converts a domain session to a storage model
func toStorageModel(session *installation.InstallationSession) *sessionStorageModel {
	// Convert configuration
	config := session.Configuration()
	configDTO := configurationDTO{
		Components:        make([]componentSelectionDTO, 0),
		DiskAvailable:     config.DiskSpace().Available(),
		DiskRequired:      config.DiskSpace().Required(),
		MergeExistingConf: config.MergeExistingConfig(),
	}

	// Convert components
	for _, comp := range config.Components() {
		compDTO := componentSelectionDTO{
			Component: string(comp.Component()),
			Version:   comp.Version(),
		}
		if comp.HasPackageInfo() {
			compDTO.PackageName = comp.PackageInfo().Name()
			compDTO.SizeBytes = comp.PackageInfo().SizeBytes()
			compDTO.Dependencies = comp.PackageInfo().Dependencies()
		}
		configDTO.Components = append(configDTO.Components, compDTO)
	}

	// Convert GPU support if present
	if config.HasGPUSupport() {
		configDTO.GPUVendor = config.GPUSupport().Vendor()
		configDTO.GPURequiresDriver = config.GPUSupport().RequiresDriver()
		configDTO.GPUDriverComponent = string(config.GPUSupport().DriverComponent())
	}

	// Convert snapshot if present
	var snapDTO *snapshotDTO
	if snapshot := session.Snapshot(); snapshot != nil {
		snapDTO = &snapshotDTO{
			ID:               snapshot.ID(),
			CreatedAt:        snapshot.CreatedAt(),
			Path:             snapshot.Path(),
			DiskAvailable:    snapshot.DiskSpace().Available(),
			DiskRequired:     snapshot.DiskSpace().Required(),
			Packages:         snapshot.Packages(),
			Corrupted:        snapshot.IsCorrupted(),
			CorruptionReason: snapshot.CorruptionReason(),
		}
	}

	// Convert installed components
	installedDTOs := make([]installedComponentDTO, 0)
	for _, comp := range session.InstalledComponents() {
		compDTO := installedComponentDTO{
			ID:          comp.ID(),
			Component:   string(comp.Component()),
			Version:     comp.Version(),
			InstalledAt: comp.InstalledAt(),
			Verified:    comp.IsVerified(),
		}
		if comp.HasPackageInfo() {
			compDTO.PackageName = comp.PackageInfo().Name()
			compDTO.SizeBytes = comp.PackageInfo().SizeBytes()
			compDTO.Dependencies = comp.PackageInfo().Dependencies()
		}
		if comp.IsVerified() {
			compDTO.VerifiedAt = comp.VerifiedAt()
		}
		installedDTOs = append(installedDTOs, compDTO)
	}

	return &sessionStorageModel{
		ID:                  session.ID(),
		Configuration:       configDTO,
		Status:              string(session.Status()),
		Snapshot:            snapDTO,
		InstalledComponents: installedDTOs,
		StartedAt:           session.StartedAt(),
		CompletedAt:         session.CompletedAt(),
		FailureReason:       session.FailureReason(),
	}
}

// Save persists an installation session
func (r *SQLiteSimpleSessionRepository) Save(ctx context.Context, session *installation.InstallationSession) error {
	// Convert to storage model
	model := toStorageModel(session)

	// Serialize to JSON
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	now := time.Now()

	// Upsert query
	query := `
	INSERT INTO sessions (id, status, data, started_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status,
		data = excluded.data,
		updated_at = excluded.updated_at
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		session.ID(),
		string(session.Status()),
		string(data),
		session.StartedAt(),
		now,
	)

	return err
}

// fromStorageModel converts a storage DTO back to domain objects
func fromStorageModel(model *sessionStorageModel) (*installation.InstallationSession, error) {
	// Reconstruct component selections
	components := make([]installation.ComponentSelection, 0, len(model.Configuration.Components))
	for _, compDTO := range model.Configuration.Components {
		// Create package info if present
		var packageInfo *installation.PackageInfo
		if compDTO.PackageName != "" {
			pkgInfo, err := installation.NewPackageInfo(
				compDTO.PackageName,
				compDTO.Version,
				compDTO.SizeBytes,
				compDTO.Dependencies,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create package info: %w", err)
			}
			packageInfo = &pkgInfo
		}

		// Create component selection
		compSel, err := installation.NewComponentSelection(
			installation.ComponentName(compDTO.Component),
			compDTO.Version,
			packageInfo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create component selection: %w", err)
		}
		components = append(components, compSel)
	}

	// Reconstruct disk space
	diskSpace, err := installation.NewDiskSpace(
		model.Configuration.DiskAvailable,
		model.Configuration.DiskRequired,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk space: %w", err)
	}

	// Reconstruct GPU support if present
	var gpuSupport *installation.GPUSupport
	if model.Configuration.GPUVendor != "" {
		gpu, err := installation.NewGPUSupport(
			model.Configuration.GPUVendor,
			model.Configuration.GPURequiresDriver,
			installation.ComponentName(model.Configuration.GPUDriverComponent),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create GPU support: %w", err)
		}
		gpuSupport = &gpu
	}

	// Reconstruct configuration
	config, err := installation.NewInstallationConfiguration(
		components,
		gpuSupport,
		diskSpace,
		model.Configuration.MergeExistingConf,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	// Reconstruct snapshot if present
	var snapshot *installation.SystemSnapshot
	if model.Snapshot != nil {
		snapDiskSpace, err := installation.NewDiskSpace(
			model.Snapshot.DiskAvailable,
			model.Snapshot.DiskRequired,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot disk space: %w", err)
		}

		snapshot, err = installation.NewSystemSnapshot(
			model.Snapshot.Path,
			snapDiskSpace,
			model.Snapshot.Packages,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot: %w", err)
		}

		if model.Snapshot.Corrupted {
			snapshot.MarkAsCorrupted(model.Snapshot.CorruptionReason)
		}
	}

	// Reconstruct installed components
	installedComponents := make([]*installation.InstalledComponent, 0, len(model.InstalledComponents))
	for _, compDTO := range model.InstalledComponents {
		// Create package info if present
		var packageInfo *installation.PackageInfo
		if compDTO.PackageName != "" {
			pkgInfo, err := installation.NewPackageInfo(
				compDTO.PackageName,
				compDTO.Version,
				compDTO.SizeBytes,
				compDTO.Dependencies,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create package info: %w", err)
			}
			packageInfo = &pkgInfo
		}

		// Create installed component
		comp, err := installation.NewInstalledComponent(
			installation.ComponentName(compDTO.Component),
			compDTO.Version,
			packageInfo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create installed component: %w", err)
		}

		if compDTO.Verified {
			comp.MarkAsVerified()
		}

		installedComponents = append(installedComponents, comp)
	}

	// Reconstruct session using domain factory
	session, err := installation.ReconstructInstallationSession(
		model.ID,
		config,
		installation.InstallationStatus(model.Status),
		snapshot,
		installedComponents,
		model.StartedAt,
		model.CompletedAt,
		model.FailureReason,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct session: %w", err)
	}

	return session, nil
}

// FindByID retrieves an installation session by its ID
func (r *SQLiteSimpleSessionRepository) FindByID(ctx context.Context, id string) (*installation.InstallationSession, error) {
	query := `SELECT id, status, data, started_at, updated_at FROM sessions WHERE id = ?`

	var record sessionRecord
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.Status,
		&record.Data,
		&record.StartedAt,
		&record.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found: %s", id)
		}
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// Deserialize storage model
	var model sessionStorageModel
	if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	// Convert to domain using helper
	session, err := fromStorageModel(&model)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// List retrieves all installation sessions
func (r *SQLiteSimpleSessionRepository) List(ctx context.Context) ([]*installation.InstallationSession, error) {
	query := `SELECT id, status, data, started_at, updated_at FROM sessions ORDER BY started_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*installation.InstallationSession

	for rows.Next() {
		var record sessionRecord
		err := rows.Scan(
			&record.ID,
			&record.Status,
			&record.Data,
			&record.StartedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

		// Deserialize storage model
		var model sessionStorageModel
		if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
		}

		// Convert to domain using helper
		session, err := fromStorageModel(&model)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// Close closes the database connection
func (r *SQLiteSimpleSessionRepository) Close() error {
	return r.db.Close()
}

// Count returns the number of sessions (useful for testing)
func (r *SQLiteSimpleSessionRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions").Scan(&count)
	return count, err
}

// Clear removes all sessions (useful for testing)
func (r *SQLiteSimpleSessionRepository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions")
	return err
}
