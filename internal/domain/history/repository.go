package history

import (
	"context"
	"time"
)

// Repository defines the interface for persisting and querying installation records
type Repository interface {
	// Save persists an installation record
	Save(ctx context.Context, record InstallationRecord) error

	// FindByID retrieves a record by its ID
	FindByID(ctx context.Context, id RecordID) (InstallationRecord, error)

	// FindAll retrieves records matching the filter
	// If filter is empty, returns all records
	FindAll(ctx context.Context, filter RecordFilter) ([]InstallationRecord, error)

	// FindRecent retrieves the most recent records up to the specified limit
	// Records are ordered by recordedAt descending (newest first)
	FindRecent(ctx context.Context, limit int) ([]InstallationRecord, error)

	// Count returns the number of records matching the filter
	Count(ctx context.Context, filter RecordFilter) (int, error)

	// Delete removes a record by its ID
	Delete(ctx context.Context, id RecordID) error

	// PurgeOlderThan removes all records with recordedAt before the cutoff date
	// Returns the number of records deleted
	PurgeOlderThan(ctx context.Context, cutoffDate time.Time) (int, error)

	// ExportRecords exports records matching the filter to a serialized format
	ExportRecords(ctx context.Context, filter RecordFilter) ([]byte, error)

	// ImportRecords imports records from a serialized format
	// Returns the number of records imported
	ImportRecords(ctx context.Context, data []byte) (int, error)
}
