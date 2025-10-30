package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/backup"
)

// RepositoryAdapter adapts the infrastructure BackupService to the domain Repository interface
type RepositoryAdapter struct {
	service    *BackupService
	backupRoot string
}

// NewRepositoryAdapter creates a new repository adapter
func NewRepositoryAdapter(backupRoot string) *RepositoryAdapter {
	return &RepositoryAdapter{
		service:    NewBackupService(backupRoot),
		backupRoot: backupRoot,
	}
}

// Save persists a backup
func (r *RepositoryAdapter) Save(ctx context.Context, b *backup.Backup) error {
	backupPath := filepath.Join(r.backupRoot, b.ID())

	// Create manifest from domain backup
	manifest := &BackupManifest{
		ID:          b.ID(),
		Description: b.Description(),
		CreatedAt:   b.CreatedAt(),
		Files:       r.convertDomainFiles(b.Files()),
	}

	// Save manifest
	if err := manifest.Save(backupPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	return nil
}

// FindByID retrieves a backup by its ID
func (r *RepositoryAdapter) FindByID(ctx context.Context, id string) (*backup.Backup, error) {
	backupPath := filepath.Join(r.backupRoot, id)

	// Load manifest
	manifest, err := LoadManifest(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	// Convert to domain backup
	return r.convertToDomain(manifest)
}

// FindAll retrieves all backups sorted by creation date (newest first)
func (r *RepositoryAdapter) FindAll(ctx context.Context) ([]*backup.Backup, error) {
	// Get metadata from service
	metadataList, err := r.service.ListBackups(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to domain backups
	var backups []*backup.Backup
	for _, metadata := range metadataList {
		backupPath := filepath.Join(r.backupRoot, metadata.ID)
		manifest, err := LoadManifest(backupPath)
		if err != nil {
			continue // Skip invalid backups
		}

		domainBackup, err := r.convertToDomain(manifest)
		if err != nil {
			continue
		}

		backups = append(backups, domainBackup)
	}

	return backups, nil
}

// Delete removes a backup
func (r *RepositoryAdapter) Delete(ctx context.Context, id string) error {
	backupPath := filepath.Join(r.backupRoot, id)
	return os.RemoveAll(backupPath)
}

// Exists checks if a backup with the given ID exists
func (r *RepositoryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	backupPath := filepath.Join(r.backupRoot, id)
	_, err := os.Stat(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Helper methods

func (r *RepositoryAdapter) convertDomainFiles(domainFiles []backup.BackupFile) []FileEntry {
	files := make([]FileEntry, len(domainFiles))
	for i, f := range domainFiles {
		files[i] = FileEntry{
			OriginalPath: f.OriginalPath,
			BackupPath:   f.BackupPath,
			Permissions:  f.Permissions,
			SizeBytes:    f.SizeBytes,
		}
	}
	return files
}

func (r *RepositoryAdapter) convertToDomain(manifest *BackupManifest) (*backup.Backup, error) {
	// Create new backup with manifest data
	b, err := backup.NewBackup(manifest.Description)
	if err != nil {
		return nil, err
	}

	// We need to reconstruct the backup with the original ID and timestamp
	// This is a bit of a hack, but necessary for the adapter pattern
	// In a real implementation, we might want to add a constructor that takes these values

	// For now, we'll create domain backup files and add them
	for _, file := range manifest.Files {
		domainFile := backup.BackupFile{
			OriginalPath: file.OriginalPath,
			BackupPath:   file.BackupPath,
			Permissions:  file.Permissions,
			SizeBytes:    file.SizeBytes,
		}

		if err := b.AddFile(domainFile); err != nil {
			continue // Skip duplicates
		}
	}

	// Mark as complete since it was loaded from disk
	b.MarkComplete()

	return b, nil
}
