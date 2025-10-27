package services

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/history"
)

// HistoryQueryService provides read-only operations for querying installation history
type HistoryQueryService struct {
	historyRepo history.Repository
}

// NewHistoryQueryService creates a new history query service
func NewHistoryQueryService(historyRepo history.Repository) *HistoryQueryService {
	return &HistoryQueryService{
		historyRepo: historyRepo,
	}
}

// ListRecords retrieves installation records matching the provided filter
func (s *HistoryQueryService) ListRecords(
	ctx context.Context,
	filter history.RecordFilter,
) ([]history.InstallationRecord, error) {
	return s.historyRepo.FindAll(ctx, filter)
}

// GetRecordByID retrieves a specific installation record by its ID
func (s *HistoryQueryService) GetRecordByID(
	ctx context.Context,
	id history.RecordID,
) (history.InstallationRecord, error) {
	return s.historyRepo.FindByID(ctx, id)
}

// ListRecent retrieves the most recent installation records up to the specified limit
func (s *HistoryQueryService) ListRecent(
	ctx context.Context,
	limit int,
) ([]history.InstallationRecord, error) {
	return s.historyRepo.FindRecent(ctx, limit)
}

// CountRecords returns the number of records matching the provided filter
func (s *HistoryQueryService) CountRecords(
	ctx context.Context,
	filter history.RecordFilter,
) (int, error) {
	return s.historyRepo.Count(ctx, filter)
}
