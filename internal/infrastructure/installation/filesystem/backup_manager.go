package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles backup and restoration of files and directories
type BackupManager struct{}

// NewBackupManager creates a new backup manager
func NewBackupManager() *BackupManager {
	return &BackupManager{}
}

// CreateBackup creates a backup of a single file
// Returns the path to the backup file
func (b *BackupManager) CreateBackup(ctx context.Context, sourcePath, backupDir string) (string, error) {
	// Verify source exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("source file not found: %w", err)
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	baseFilename := filepath.Base(sourcePath)
	backupFilename := fmt.Sprintf("%s.%s.backup", baseFilename, timestamp)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Copy file
	if err := b.copyFile(sourcePath, backupPath, sourceInfo.Mode()); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return backupPath, nil
}

// CreateBackupDirectory creates a backup of an entire directory recursively
// Returns the path to the backup directory
func (b *BackupManager) CreateBackupDirectory(ctx context.Context, sourceDir, backupDir string) (string, error) {
	// Verify source directory exists
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		return "", fmt.Errorf("source directory not found: %w", err)
	}

	if !sourceInfo.IsDir() {
		return "", fmt.Errorf("source is not a directory: %s", sourceDir)
	}

	// Create backup directory
	timestamp := time.Now().Format("20060102-150405")
	baseDirname := filepath.Base(sourceDir)
	backupDirname := fmt.Sprintf("%s.%s.backup", baseDirname, timestamp)
	backupPath := filepath.Join(backupDir, backupDirname)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy directory contents recursively
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(backupPath, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return b.copyFile(path, destPath, info.Mode())
	})

	if err != nil {
		return "", fmt.Errorf("failed to copy directory: %w", err)
	}

	return backupPath, nil
}

// RestoreBackup restores a backup file to a specified location
func (b *BackupManager) RestoreBackup(ctx context.Context, backupPath, restorePath string) error {
	// Verify backup exists
	backupInfo, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Create parent directories if needed
	restoreDir := filepath.Dir(restorePath)
	if err := os.MkdirAll(restoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create restore directory: %w", err)
	}

	// Copy backup to restore location
	if err := b.copyFile(backupPath, restorePath, backupInfo.Mode()); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return nil
}

// VerifyBackup verifies that a backup file exists and is readable
func (b *BackupManager) VerifyBackup(ctx context.Context, backupPath string) (bool, error) {
	info, err := os.Stat(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Verify it's a regular file
	if !info.Mode().IsRegular() {
		return false, nil
	}

	// Verify it's readable
	file, err := os.Open(backupPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	return true, nil
}

// copyFile copies a file from source to destination with specified permissions
func (b *BackupManager) copyFile(src, dst string, mode os.FileMode) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file with same permissions
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Ensure all data is written
	return destFile.Sync()
}
