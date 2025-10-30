package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/backup"
)

// RestoreBackupRequest contains parameters for restoring a backup
type RestoreBackupRequest struct {
	BackupID string   // ID of the backup to restore
	Selective []string // Optional: specific files to restore (empty = all)
}

// RestoreBackupResponse contains the result of backup restoration
type RestoreBackupResponse struct {
	BackupID      string   // ID of the restored backup
	FilesRestored int      // Number of files restored
	RestoredPaths []string // Paths of restored files
}

// RestoreBackupUseCase handles restoring configuration backups
type RestoreBackupUseCase struct {
	repository backup.Repository
}

// NewRestoreBackupUseCase creates a new use case instance
func NewRestoreBackupUseCase(repository backup.Repository) *RestoreBackupUseCase {
	return &RestoreBackupUseCase{
		repository: repository,
	}
}

// Execute restores files from a backup
func (uc *RestoreBackupUseCase) Execute(ctx context.Context, req RestoreBackupRequest) (*RestoreBackupResponse, error) {
	// Validate request
	if req.BackupID == "" {
		return nil, fmt.Errorf("backup ID is required")
	}

	// Find backup
	b, err := uc.repository.FindByID(ctx, req.BackupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find backup: %w", err)
	}

	// Check backup is complete
	if !b.IsComplete() {
		return nil, fmt.Errorf("backup is not complete (status: %s)", b.Status())
	}

	var restoredPaths []string

	// Restore files
	for _, file := range b.Files() {
		// Skip if selective restore and file not in list
		if len(req.Selective) > 0 && !contains(req.Selective, file.OriginalPath) {
			continue
		}

		// Restore file
		if err := uc.restoreFile(file); err != nil {
			// Log error but continue
			continue
		}

		restoredPaths = append(restoredPaths, file.OriginalPath)
	}

	return &RestoreBackupResponse{
		BackupID:      b.ID(),
		FilesRestored: len(restoredPaths),
		RestoredPaths: restoredPaths,
	}, nil
}

func (uc *RestoreBackupUseCase) restoreFile(file backup.BackupFile) error {
	// Read backup file
	content, err := os.ReadFile(file.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Ensure target directory exists
	targetDir := filepath.Dir(file.OriginalPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write to original location
	if err := os.WriteFile(file.OriginalPath, content, file.Permissions); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
