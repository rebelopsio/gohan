package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BackupService handles configuration backup and restore operations
type BackupService struct {
	backupRoot string // Root directory for all backups
}

// BackupMetadata contains information about a backup
type BackupMetadata struct {
	ID          string      `json:"id"`          // Timestamp-based ID
	Path        string      `json:"path"`        // Full path to backup directory
	Description string      `json:"description"` // User-provided description
	CreatedAt   time.Time   `json:"created_at"`  // When backup was created
	Files       []FileEntry `json:"files"`       // Files in this backup
	SizeBytes   int64       `json:"size_bytes"`  // Total backup size
}

// FileEntry represents a backed up file
type FileEntry struct {
	OriginalPath string      `json:"original_path"` // Where file came from
	BackupPath   string      `json:"backup_path"`   // Where it's stored in backup
	Permissions  os.FileMode `json:"permissions"`   // File permissions
	SizeBytes    int64       `json:"size_bytes"`    // File size
}

// BackupManifest is stored in each backup directory
type BackupManifest struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Files       []FileEntry `json:"files"`
}

// NewBackupService creates a new backup service
func NewBackupService(backupRoot string) *BackupService {
	return &BackupService{
		backupRoot: backupRoot,
	}
}

// BackupFile backs up a single file
func (s *BackupService) BackupFile(ctx context.Context, filePath string) (string, error) {
	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// Generate backup path
	timestamp := time.Now().Format("2006-01-02_150405")
	backupDir := filepath.Join(s.backupRoot, timestamp)

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Determine backup file path
	fileName := filepath.Base(filePath)
	backupPath := filepath.Join(backupDir, fileName)

	// Copy file preserving permissions
	if err := copyFile(filePath, backupPath); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Preserve permissions
	if err := os.Chmod(backupPath, info.Mode().Perm()); err != nil {
		return "", fmt.Errorf("failed to set permissions: %w", err)
	}

	return backupPath, nil
}

// BackupDirectory backs up an entire directory
func (s *BackupService) BackupDirectory(ctx context.Context, dirPath string) (string, error) {
	// Verify directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat directory %s: %w", dirPath, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", dirPath)
	}

	// Generate backup path
	timestamp := time.Now().Format("2006-01-02_150405")
	dirName := filepath.Base(dirPath)
	backupPath := filepath.Join(s.backupRoot, timestamp, dirName)

	// Copy directory recursively
	if err := copyDir(dirPath, backupPath); err != nil {
		return "", fmt.Errorf("failed to copy directory: %w", err)
	}

	return backupPath, nil
}

// CreateBackup creates a complete backup of multiple files
func (s *BackupService) CreateBackup(ctx context.Context, filePaths []string, description string) (*BackupMetadata, error) {
	// Generate backup ID (timestamp)
	timestamp := time.Now()
	backupID := timestamp.Format("2006-01-02_150405")
	backupPath := filepath.Join(s.backupRoot, backupID)

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create manifest
	manifest := NewBackupManifest(backupID, description)

	// Backup each file
	var totalSize int64
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			continue // Skip files that don't exist
		}

		// Determine relative backup path
		fileName := filepath.Base(filePath)
		backupFilePath := filepath.Join(backupPath, fileName)

		// Copy file
		if err := copyFile(filePath, backupFilePath); err != nil {
			continue
		}

		// Preserve permissions
		os.Chmod(backupFilePath, info.Mode().Perm())

		// Add to manifest
		manifest.AddFile(filePath, backupFilePath, info.Mode().Perm())
		totalSize += info.Size()
	}

	// Save manifest
	if err := manifest.Save(backupPath); err != nil {
		return nil, fmt.Errorf("failed to save manifest: %w", err)
	}

	// Create metadata
	metadata := &BackupMetadata{
		ID:          backupID,
		Path:        backupPath,
		Description: description,
		CreatedAt:   timestamp,
		Files:       manifest.Files,
		SizeBytes:   totalSize,
	}

	return metadata, nil
}

// RestoreBackup restores files from a backup
func (s *BackupService) RestoreBackup(ctx context.Context, backupID string) error {
	backupPath := filepath.Join(s.backupRoot, backupID)

	// Load manifest
	manifest, err := LoadManifest(backupPath)
	if err != nil {
		return fmt.Errorf("failed to load backup manifest: %w", err)
	}

	// Restore each file
	for _, fileEntry := range manifest.Files {
		// Read backup file
		backupFilePath := fileEntry.BackupPath
		content, err := os.ReadFile(backupFilePath)
		if err != nil {
			return fmt.Errorf("failed to read backup file %s: %w", backupFilePath, err)
		}

		// Ensure target directory exists
		targetDir := filepath.Dir(fileEntry.OriginalPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		// Write to original location
		if err := os.WriteFile(fileEntry.OriginalPath, content, fileEntry.Permissions); err != nil {
			return fmt.Errorf("failed to restore file %s: %w", fileEntry.OriginalPath, err)
		}
	}

	return nil
}

// ListBackups lists all available backups sorted by date (newest first)
func (s *BackupService) ListBackups(ctx context.Context) ([]*BackupMetadata, error) {
	// Create backup root if it doesn't exist
	if err := os.MkdirAll(s.backupRoot, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup root: %w", err)
	}

	// Read backup directories
	entries, err := os.ReadDir(s.backupRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []*BackupMetadata

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		backupPath := filepath.Join(s.backupRoot, entry.Name())

		// Load manifest
		manifest, err := LoadManifest(backupPath)
		if err != nil {
			continue // Skip invalid backups
		}

		// Calculate total size
		var totalSize int64
		for _, file := range manifest.Files {
			totalSize += file.SizeBytes
		}

		metadata := &BackupMetadata{
			ID:          manifest.ID,
			Path:        backupPath,
			Description: manifest.Description,
			CreatedAt:   manifest.CreatedAt,
			Files:       manifest.Files,
			SizeBytes:   totalSize,
		}

		backups = append(backups, metadata)
	}

	// Sort by creation time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

// CleanupOldBackups removes backups older than retentionDays
func (s *BackupService) CleanupOldBackups(ctx context.Context, retentionDays int) (int, error) {
	if retentionDays <= 0 {
		return 0, nil // Don't remove anything if retention is 0 or negative
	}

	// Get all backups
	backups, err := s.ListBackups(ctx)
	if err != nil {
		return 0, err
	}

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	removed := 0

	for _, backup := range backups {
		if backup.CreatedAt.Before(cutoffDate) {
			// Remove backup directory
			if err := os.RemoveAll(backup.Path); err != nil {
				// Log error but continue
				continue
			}
			removed++
		}
	}

	return removed, nil
}

// GetBackupInfo retrieves metadata for a specific backup
func (s *BackupService) GetBackupInfo(ctx context.Context, backupID string) (*BackupMetadata, error) {
	backupPath := filepath.Join(s.backupRoot, backupID)

	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		return nil, fmt.Errorf("backup %s not found: %w", backupID, err)
	}

	// Load manifest
	manifest, err := LoadManifest(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	// Calculate size
	var totalSize int64
	for _, file := range manifest.Files {
		totalSize += file.SizeBytes
	}

	metadata := &BackupMetadata{
		ID:          manifest.ID,
		Path:        backupPath,
		Description: manifest.Description,
		CreatedAt:   manifest.CreatedAt,
		Files:       manifest.Files,
		SizeBytes:   totalSize,
	}

	return metadata, nil
}

// NewBackupManifest creates a new backup manifest
func NewBackupManifest(id, description string) *BackupManifest {
	return &BackupManifest{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
		Files:       []FileEntry{},
	}
}

// AddFile adds a file entry to the manifest
func (m *BackupManifest) AddFile(originalPath, backupPath string, permissions os.FileMode) {
	info, _ := os.Stat(backupPath)
	var size int64
	if info != nil {
		size = info.Size()
	}

	m.Files = append(m.Files, FileEntry{
		OriginalPath: originalPath,
		BackupPath:   backupPath,
		Permissions:  permissions,
		SizeBytes:    size,
	})
}

// Save writes the manifest to disk
func (m *BackupManifest) Save(backupPath string) error {
	manifestPath := filepath.Join(backupPath, "manifest.json")

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// LoadManifest loads a manifest from disk
func LoadManifest(backupPath string) (*BackupManifest, error) {
	manifestPath := filepath.Join(backupPath, "manifest.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest BackupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	return &manifest, nil
}

// Helper functions

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read directory contents
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}

			// Preserve permissions
			info, _ := os.Stat(srcPath)
			if info != nil {
				os.Chmod(dstPath, info.Mode())
			}
		}
	}

	return nil
}
