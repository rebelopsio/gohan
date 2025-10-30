package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/backup"
)

// CreateBackupRequest contains parameters for creating a backup
type CreateBackupRequest struct {
	Description string   // Description of the backup
	FilePaths   []string // Paths of files/directories to back up
	BackupRoot  string   // Root directory where backups are stored
}

// CreateBackupResponse contains the result of backup creation
type CreateBackupResponse struct {
	BackupID    string // ID of the created backup
	FileCount   int    // Number of files backed up
	TotalSize   int64  // Total size in bytes
	BackupPath  string // Full path to the backup directory
}

// CreateBackupUseCase handles creating configuration backups
type CreateBackupUseCase struct {
	repository backup.Repository
}

// NewCreateBackupUseCase creates a new use case instance
func NewCreateBackupUseCase(repository backup.Repository) *CreateBackupUseCase {
	return &CreateBackupUseCase{
		repository: repository,
	}
}

// Execute creates a backup of the specified files
func (uc *CreateBackupUseCase) Execute(ctx context.Context, req CreateBackupRequest) (*CreateBackupResponse, error) {
	// Validate request
	if err := uc.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create backup entity
	b, err := backup.NewBackup(req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Determine backup directory path
	backupPath := filepath.Join(req.BackupRoot, b.ID())

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		b.MarkFailed(fmt.Sprintf("failed to create backup directory: %v", err))
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Back up each file
	for _, filePath := range req.FilePaths {
		// Skip if file doesn't exist
		if _, err := os.Stat(filePath); err != nil {
			continue // Skip missing files
		}

		// Handle both files and directories
		if err := uc.backupPath(ctx, b, filePath, backupPath); err != nil {
			// Continue on error but log it
			continue
		}
	}

	// Mark backup as complete
	b.MarkComplete()

	// Save backup through repository
	if err := uc.repository.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("failed to save backup: %w", err)
	}

	return &CreateBackupResponse{
		BackupID:   b.ID(),
		FileCount:  b.FileCount(),
		TotalSize:  b.TotalSize(),
		BackupPath: backupPath,
	}, nil
}

func (uc *CreateBackupUseCase) validateRequest(req CreateBackupRequest) error {
	if req.BackupRoot == "" {
		return fmt.Errorf("backup root is required")
	}

	if len(req.FilePaths) == 0 {
		return fmt.Errorf("at least one file path is required")
	}

	return nil
}

func (uc *CreateBackupUseCase) backupPath(ctx context.Context, b *backup.Backup, sourcePath, backupDir string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return uc.backupDirectory(ctx, b, sourcePath, backupDir)
	}

	return uc.backupFile(ctx, b, sourcePath, backupDir)
}

func (uc *CreateBackupUseCase) backupFile(ctx context.Context, b *backup.Backup, sourcePath, backupDir string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// Determine backup file path
	fileName := filepath.Base(sourcePath)
	backupFilePath := filepath.Join(backupDir, fileName)

	// Copy file
	if err := copyFile(sourcePath, backupFilePath); err != nil {
		return err
	}

	// Add to backup
	backupFile := backup.BackupFile{
		OriginalPath: sourcePath,
		BackupPath:   backupFilePath,
		Permissions:  info.Mode().Perm(),
		SizeBytes:    info.Size(),
	}

	return b.AddFile(backupFile)
}

func (uc *CreateBackupUseCase) backupDirectory(ctx context.Context, b *backup.Backup, sourcePath, backupDir string) error {
	// Walk directory tree
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip directories themselves
		if info.IsDir() {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return nil
		}

		// Determine backup path
		backupFilePath := filepath.Join(backupDir, filepath.Base(sourcePath), relPath)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(backupFilePath), 0755); err != nil {
			return nil
		}

		// Copy file
		if err := copyFile(path, backupFilePath); err != nil {
			return nil
		}

		// Add to backup
		backupFile := backup.BackupFile{
			OriginalPath: path,
			BackupPath:   backupFilePath,
			Permissions:  info.Mode().Perm(),
			SizeBytes:    info.Size(),
		}

		b.AddFile(backupFile)
		return nil
	})
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, info.Mode().Perm())
}
