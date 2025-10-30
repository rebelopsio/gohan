package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/backup"
)

// CleanupBackupsRequest contains parameters for cleaning up old backups
type CleanupBackupsRequest struct {
	RetentionDays int  // Keep backups newer than this many days
	KeepMinimum   int  // Always keep at least this many backups
	DryRun        bool // If true, don't actually delete
}

// CleanupBackupsResponse contains the result of cleanup
type CleanupBackupsResponse struct {
	RemovedCount   int      // Number of backups removed
	RemovedIDs     []string // IDs of removed backups
	FreedBytes     int64    // Bytes freed
	RemainingCount int      // Number of backups remaining
}

// CleanupBackupsUseCase handles cleanup of old backups
type CleanupBackupsUseCase struct {
	repository backup.Repository
}

// NewCleanupBackupsUseCase creates a new use case instance
func NewCleanupBackupsUseCase(repository backup.Repository) *CleanupBackupsUseCase {
	return &CleanupBackupsUseCase{
		repository: repository,
	}
}

// Execute removes old backups based on retention policy
func (uc *CleanupBackupsUseCase) Execute(ctx context.Context, req CleanupBackupsRequest) (*CleanupBackupsResponse, error) {
	// Validate request
	if req.RetentionDays < 0 {
		return nil, fmt.Errorf("retention days cannot be negative")
	}

	if req.KeepMinimum < 0 {
		return nil, fmt.Errorf("keep minimum cannot be negative")
	}

	// Get all backups
	allBackups, err := uc.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -req.RetentionDays)

	var toRemove []*backup.Backup
	var toKeep []*backup.Backup

	// Determine which backups to remove
	for _, b := range allBackups {
		if b.CreatedAt().Before(cutoffDate) {
			toRemove = append(toRemove, b)
		} else {
			toKeep = append(toKeep, b)
		}
	}

	// Ensure we keep minimum number of backups
	if len(toKeep) < req.KeepMinimum {
		// Move some from toRemove to toKeep to maintain minimum
		needed := req.KeepMinimum - len(toKeep)
		if needed <= len(toRemove) {
			// Keep the newest N from toRemove
			toKeep = append(toKeep, toRemove[:needed]...)
			toRemove = toRemove[needed:]
		} else {
			// Keep all
			toKeep = append(toKeep, toRemove...)
			toRemove = nil
		}
	}

	response := &CleanupBackupsResponse{
		RemovedCount:   0,
		RemovedIDs:     []string{},
		FreedBytes:     0,
		RemainingCount: len(allBackups) - len(toRemove),
	}

	// Remove backups
	for _, b := range toRemove {
		if req.DryRun {
			// Dry run - just count
			response.RemovedCount++
			response.RemovedIDs = append(response.RemovedIDs, b.ID())
			response.FreedBytes += b.TotalSize()
			continue
		}

		// Actually delete
		if err := uc.repository.Delete(ctx, b.ID()); err != nil {
			// Log error but continue
			continue
		}

		response.RemovedCount++
		response.RemovedIDs = append(response.RemovedIDs, b.ID())
		response.FreedBytes += b.TotalSize()
	}

	return response, nil
}
