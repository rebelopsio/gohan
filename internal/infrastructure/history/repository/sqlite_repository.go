package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rebelopsio/gohan/internal/domain/history"
)

// SQLiteRepository is a SQLite implementation of history.Repository
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite history repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
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

	repo := &SQLiteRepository{db: db}
	if err := repo.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return repo, nil
}

// initialize creates the necessary tables
func (r *SQLiteRepository) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS installation_records (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		outcome TEXT NOT NULL,
		package_name TEXT NOT NULL,
		recorded_at DATETIME NOT NULL,
		installed_at DATETIME NOT NULL,
		data TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_records_session_id ON installation_records(session_id);
	CREATE INDEX IF NOT EXISTS idx_records_outcome ON installation_records(outcome);
	CREATE INDEX IF NOT EXISTS idx_records_package_name ON installation_records(package_name);
	CREATE INDEX IF NOT EXISTS idx_records_recorded_at ON installation_records(recorded_at);
	CREATE INDEX IF NOT EXISTS idx_records_installed_at ON installation_records(installed_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// recordStorageModel is a serializable representation of an installation record
type recordStorageModel struct {
	ID             string                    `json:"id"`
	SessionID      string                    `json:"session_id"`
	Outcome        string                    `json:"outcome"`
	Metadata       metadataDTO               `json:"metadata"`
	SystemContext  systemContextDTO          `json:"system_context"`
	FailureDetails *failureDetailsDTO        `json:"failure_details,omitempty"`
	RecordedAt     time.Time                 `json:"recorded_at"`
}

type metadataDTO struct {
	PackageName     string              `json:"package_name"`
	TargetVersion   string              `json:"target_version"`
	InstalledAt     time.Time           `json:"installed_at"`
	CompletedAt     time.Time           `json:"completed_at"`
	InstalledPackages []installedPackageDTO `json:"installed_packages"`
}

type installedPackageDTO struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	SizeBytes uint64 `json:"size_bytes"`
}

type systemContextDTO struct {
	OSVersion     string `json:"os_version"`
	KernelVersion string `json:"kernel_version"`
	GohanVersion  string `json:"gohan_version"`
	Hostname      string `json:"hostname"`
}

type failureDetailsDTO struct {
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
	Phase     string    `json:"phase"`
	ErrorCode string    `json:"error_code"`
}

// toStorageModel converts a domain record to a storage model
func toStorageModel(record history.InstallationRecord) *recordStorageModel {
	// Convert metadata
	metadata := record.Metadata()
	var pkgDTOs []installedPackageDTO
	for _, pkg := range metadata.InstalledPackages() {
		pkgDTOs = append(pkgDTOs, installedPackageDTO{
			Name:      pkg.Name(),
			Version:   pkg.Version(),
			SizeBytes: pkg.SizeBytes(),
		})
	}

	metadataDTO := metadataDTO{
		PackageName:       metadata.PackageName(),
		TargetVersion:     metadata.TargetVersion(),
		InstalledAt:       metadata.InstalledAt(),
		CompletedAt:       metadata.CompletedAt(),
		InstalledPackages: pkgDTOs,
	}

	// Convert system context
	sysCtx := record.SystemContext()
	systemContextDTO := systemContextDTO{
		OSVersion:     sysCtx.OSVersion(),
		KernelVersion: sysCtx.KernelVersion(),
		GohanVersion:  sysCtx.GohanVersion(),
		Hostname:      sysCtx.Hostname(),
	}

	// Convert failure details if present
	var failureDTO *failureDetailsDTO
	if record.HasFailureDetails() {
		fd := record.FailureDetails()
		failureDTO = &failureDetailsDTO{
			Reason:    fd.Reason(),
			FailedAt:  fd.FailedAt(),
			Phase:     fd.Phase(),
			ErrorCode: fd.ErrorCode(),
		}
	}

	return &recordStorageModel{
		ID:             record.ID().String(),
		SessionID:      record.SessionID(),
		Outcome:        record.Outcome().String(),
		Metadata:       metadataDTO,
		SystemContext:  systemContextDTO,
		FailureDetails: failureDTO,
		RecordedAt:     record.RecordedAt(),
	}
}

// fromStorageModel converts a storage model back to a domain record
func fromStorageModel(model *recordStorageModel) (history.InstallationRecord, error) {
	// Reconstruct installed packages
	var packages []history.InstalledPackage
	for _, pkgDTO := range model.Metadata.InstalledPackages {
		pkg, err := history.NewInstalledPackage(pkgDTO.Name, pkgDTO.Version, pkgDTO.SizeBytes)
		if err != nil {
			return history.InstallationRecord{}, fmt.Errorf("failed to create installed package: %w", err)
		}
		packages = append(packages, pkg)
	}

	// Reconstruct metadata
	metadata, err := history.NewInstallationMetadata(
		model.Metadata.PackageName,
		model.Metadata.TargetVersion,
		model.Metadata.InstalledAt,
		model.Metadata.CompletedAt,
		packages,
	)
	if err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to create metadata: %w", err)
	}

	// Reconstruct system context
	sysCtx, err := history.NewSystemContext(
		model.SystemContext.OSVersion,
		model.SystemContext.KernelVersion,
		model.SystemContext.GohanVersion,
		model.SystemContext.Hostname,
	)
	if err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to create system context: %w", err)
	}

	// Reconstruct failure details if present
	var failureDetails *history.FailureDetails
	if model.FailureDetails != nil {
		fd, err := history.NewFailureDetails(
			model.FailureDetails.Reason,
			model.FailureDetails.FailedAt,
			model.FailureDetails.Phase,
			model.FailureDetails.ErrorCode,
		)
		if err != nil {
			return history.InstallationRecord{}, fmt.Errorf("failed to create failure details: %w", err)
		}
		failureDetails = &fd
	}

	// Reconstruct outcome
	outcome, err := history.NewInstallationOutcome(model.Outcome)
	if err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to create outcome: %w", err)
	}

	// Parse record ID
	recordID, err := history.ParseRecordID(model.ID)
	if err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to parse record ID: %w", err)
	}

	// Reconstruct installation record with preserved ID
	record, err := history.ReconstructInstallationRecord(
		recordID,
		model.SessionID,
		outcome,
		metadata,
		sysCtx,
		failureDetails,
		model.RecordedAt,
	)
	if err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to reconstruct record: %w", err)
	}

	return record, nil
}

// Save persists an installation record
func (r *SQLiteRepository) Save(ctx context.Context, record history.InstallationRecord) error {
	// Convert to storage model
	model := toStorageModel(record)

	// Serialize to JSON
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	// Upsert query
	query := `
	INSERT INTO installation_records (id, session_id, outcome, package_name, recorded_at, installed_at, data)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		session_id = excluded.session_id,
		outcome = excluded.outcome,
		package_name = excluded.package_name,
		recorded_at = excluded.recorded_at,
		installed_at = excluded.installed_at,
		data = excluded.data
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		record.ID().String(),
		record.SessionID(),
		record.Outcome().String(),
		record.PackageName(),
		record.RecordedAt(),
		record.InstalledAt(),
		string(data),
	)

	return err
}

// FindByID retrieves a record by its ID
func (r *SQLiteRepository) FindByID(ctx context.Context, id history.RecordID) (history.InstallationRecord, error) {
	query := `SELECT data FROM installation_records WHERE id = ?`

	var data string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return history.InstallationRecord{}, history.ErrRecordNotFound
		}
		return history.InstallationRecord{}, fmt.Errorf("failed to query record: %w", err)
	}

	// Deserialize storage model
	var model recordStorageModel
	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return history.InstallationRecord{}, fmt.Errorf("failed to unmarshal record data: %w", err)
	}

	// Convert to domain
	record, err := fromStorageModel(&model)
	if err != nil {
		return history.InstallationRecord{}, err
	}

	return record, nil
}

// FindAll retrieves records matching the filter
func (r *SQLiteRepository) FindAll(ctx context.Context, filter history.RecordFilter) ([]history.InstallationRecord, error) {
	query := `SELECT data FROM installation_records`
	args := []interface{}{}

	// Build WHERE clause for SQL-optimizable filters
	whereClauses := []string{}

	if filter.HasOutcomeFilter() {
		whereClauses = append(whereClauses, "outcome = ?")
		args = append(args, filter.Outcome().String())
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	var records []history.InstallationRecord

	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		var model recordStorageModel
		if err := json.Unmarshal([]byte(data), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal record data: %w", err)
		}

		record, err := fromStorageModel(&model)
		if err != nil {
			return nil, err
		}

		// Apply remaining filters in memory
		if r.matches(record, filter) {
			records = append(records, record)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating records: %w", err)
	}

	return records, nil
}

// FindRecent retrieves the most recent records up to the specified limit
func (r *SQLiteRepository) FindRecent(ctx context.Context, limit int) ([]history.InstallationRecord, error) {
	if limit <= 0 {
		return []history.InstallationRecord{}, nil
	}

	query := `SELECT data FROM installation_records ORDER BY recorded_at DESC LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent records: %w", err)
	}
	defer rows.Close()

	var records []history.InstallationRecord

	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		var model recordStorageModel
		if err := json.Unmarshal([]byte(data), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal record data: %w", err)
		}

		record, err := fromStorageModel(&model)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating records: %w", err)
	}

	return records, nil
}

// Count returns the number of records matching the filter
func (r *SQLiteRepository) Count(ctx context.Context, filter history.RecordFilter) (int, error) {
	query := `SELECT COUNT(*) FROM installation_records`
	args := []interface{}{}

	// Build WHERE clause for SQL-optimizable filters
	whereClauses := []string{}

	if filter.HasOutcomeFilter() {
		whereClauses = append(whereClauses, "outcome = ?")
		args = append(args, filter.Outcome().String())
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}

	// If we have period or package filters, we need to count in memory
	if filter.HasPeriodFilter() || filter.HasPackageFilter() {
		records, err := r.FindAll(ctx, filter)
		if err != nil {
			return 0, err
		}
		return len(records), nil
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// Delete removes a record by its ID
func (r *SQLiteRepository) Delete(ctx context.Context, id history.RecordID) error {
	query := `DELETE FROM installation_records WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return history.ErrRecordNotFound
	}

	return nil
}

// PurgeOlderThan removes all records with recordedAt before the cutoff date
func (r *SQLiteRepository) PurgeOlderThan(ctx context.Context, cutoffDate time.Time) (int, error) {
	query := `DELETE FROM installation_records WHERE recorded_at < ?`

	result, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to purge records: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return int(rows), nil
}

// ExportRecords exports records matching the filter to a serialized format
func (r *SQLiteRepository) ExportRecords(ctx context.Context, filter history.RecordFilter) ([]byte, error) {
	// TODO: Implement serialization when needed
	return nil, nil
}

// ImportRecords imports records from a serialized format
func (r *SQLiteRepository) ImportRecords(ctx context.Context, data []byte) (int, error) {
	// TODO: Implement deserialization when needed
	return 0, nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// matches checks if a record matches the filter criteria (for in-memory filtering)
func (r *SQLiteRepository) matches(record history.InstallationRecord, filter history.RecordFilter) bool {
	// Empty filter matches everything
	if filter.IsEmpty() {
		return true
	}

	// Check outcome filter
	if !filter.MatchesOutcome(record.Outcome()) {
		return false
	}

	// Check metadata filter (period and package)
	if !filter.MatchesMetadata(record.Metadata()) {
		return false
	}

	return true
}
