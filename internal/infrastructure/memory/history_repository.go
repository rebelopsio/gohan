package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
)

// HistoryRepository is an in-memory implementation of history.Repository
type HistoryRepository struct {
	mu      sync.RWMutex
	records map[string]history.InstallationRecord
}

// NewHistoryRepository creates a new in-memory history repository
func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{
		records: make(map[string]history.InstallationRecord),
	}
}

// Save persists an installation record
func (r *HistoryRepository) Save(ctx context.Context, record history.InstallationRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[record.ID().String()] = record
	return nil
}

// FindByID retrieves a record by its ID
func (r *HistoryRepository) FindByID(ctx context.Context, id history.RecordID) (history.InstallationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id.String()]
	if !exists {
		return history.InstallationRecord{}, history.ErrRecordNotFound
	}

	return record, nil
}

// FindAll retrieves records matching the filter
func (r *HistoryRepository) FindAll(ctx context.Context, filter history.RecordFilter) ([]history.InstallationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []history.InstallationRecord

	for _, record := range r.records {
		if r.matches(record, filter) {
			results = append(results, record)
		}
	}

	return results, nil
}

// FindRecent retrieves the most recent records up to the specified limit
func (r *HistoryRepository) FindRecent(ctx context.Context, limit int) ([]history.InstallationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 {
		return []history.InstallationRecord{}, nil
	}

	// Collect all records
	var allRecords []history.InstallationRecord
	for _, record := range r.records {
		allRecords = append(allRecords, record)
	}

	// Sort by recordedAt descending (newest first)
	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].RecordedAt().After(allRecords[j].RecordedAt())
	})

	// Return up to limit
	if len(allRecords) > limit {
		return allRecords[:limit], nil
	}

	return allRecords, nil
}

// Count returns the number of records matching the filter
func (r *HistoryRepository) Count(ctx context.Context, filter history.RecordFilter) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, record := range r.records {
		if r.matches(record, filter) {
			count++
		}
	}

	return count, nil
}

// Delete removes a record by its ID
func (r *HistoryRepository) Delete(ctx context.Context, id history.RecordID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.records[id.String()]; !exists {
		return history.ErrRecordNotFound
	}

	delete(r.records, id.String())
	return nil
}

// PurgeOlderThan removes all records with recordedAt before the cutoff date
func (r *HistoryRepository) PurgeOlderThan(ctx context.Context, cutoffDate time.Time) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	deleted := 0
	for id, record := range r.records {
		if record.RecordedAt().Before(cutoffDate) {
			delete(r.records, id)
			deleted++
		}
	}

	return deleted, nil
}

// ExportRecords exports records matching the filter to a serialized format
func (r *HistoryRepository) ExportRecords(ctx context.Context, filter history.RecordFilter) ([]byte, error) {
	// TODO: Implement serialization when needed
	return nil, nil
}

// ImportRecords imports records from a serialized format
func (r *HistoryRepository) ImportRecords(ctx context.Context, data []byte) (int, error) {
	// TODO: Implement deserialization when needed
	return 0, nil
}

// matches checks if a record matches the filter criteria
func (r *HistoryRepository) matches(record history.InstallationRecord, filter history.RecordFilter) bool {
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
