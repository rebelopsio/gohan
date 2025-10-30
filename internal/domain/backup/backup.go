package backup

import (
	"fmt"
	"os"
	"time"
)

// BackupStatus represents the state of a backup
type BackupStatus string

const (
	// StatusPending indicates backup is being created
	StatusPending BackupStatus = "pending"

	// StatusComplete indicates backup completed successfully
	StatusComplete BackupStatus = "complete"

	// StatusFailed indicates backup failed
	StatusFailed BackupStatus = "failed"
)

// Backup represents a configuration backup aggregate root
type Backup struct {
	id          string
	description string
	createdAt   time.Time
	files       []BackupFile
	status      BackupStatus
	err         string // Error message if status is Failed
}

// BackupFile represents a single file in a backup (value object)
type BackupFile struct {
	OriginalPath string      // Where the file came from
	BackupPath   string      // Where it's stored in backup
	Permissions  os.FileMode // File permissions
	SizeBytes    int64       // File size in bytes
}

// NewBackup creates a new backup with a timestamp-based ID
func NewBackup(description string) (*Backup, error) {
	now := time.Now()
	id := now.Format("2006-01-02_150405")

	return &Backup{
		id:          id,
		description: description,
		createdAt:   now,
		files:       []BackupFile{},
		status:      StatusPending,
	}, nil
}

// ID returns the backup's unique identifier
func (b *Backup) ID() string {
	return b.id
}

// Description returns the backup description
func (b *Backup) Description() string {
	return b.description
}

// CreatedAt returns when the backup was created
func (b *Backup) CreatedAt() time.Time {
	return b.createdAt
}

// Files returns all files in the backup
func (b *Backup) Files() []BackupFile {
	// Return copy to prevent external modification
	files := make([]BackupFile, len(b.files))
	copy(files, b.files)
	return files
}

// Status returns the current backup status
func (b *Backup) Status() BackupStatus {
	return b.status
}

// Error returns the error message if backup failed
func (b *Backup) Error() string {
	return b.err
}

// AddFile adds a file to the backup
func (b *Backup) AddFile(file BackupFile) error {
	// Validate file
	if err := file.Validate(); err != nil {
		return fmt.Errorf("invalid file: %w", err)
	}

	// Check for duplicates
	if b.Contains(file.OriginalPath) {
		return ErrDuplicateFile
	}

	b.files = append(b.files, file)
	return nil
}

// MarkComplete marks the backup as successfully completed
func (b *Backup) MarkComplete() {
	b.status = StatusComplete
}

// MarkFailed marks the backup as failed with an error message
func (b *Backup) MarkFailed(err string) {
	b.status = StatusFailed
	b.err = err
}

// TotalSize calculates the total size of all files in the backup
func (b *Backup) TotalSize() int64 {
	var total int64
	for _, file := range b.files {
		total += file.SizeBytes
	}
	return total
}

// Contains checks if the backup contains a file with the given original path
func (b *Backup) Contains(originalPath string) bool {
	for _, file := range b.files {
		if file.OriginalPath == originalPath {
			return true
		}
	}
	return false
}

// Age returns how long ago the backup was created
func (b *Backup) Age() time.Duration {
	return time.Since(b.createdAt)
}

// IsComplete returns true if backup completed successfully
func (b *Backup) IsComplete() bool {
	return b.status == StatusComplete
}

// IsFailed returns true if backup failed
func (b *Backup) IsFailed() bool {
	return b.status == StatusFailed
}

// FileCount returns the number of files in the backup
func (b *Backup) FileCount() int {
	return len(b.files)
}

// Validate validates the entire backup
func (b *Backup) Validate() error {
	if len(b.files) == 0 {
		return ErrEmptyBackup
	}

	for _, file := range b.files {
		if err := file.Validate(); err != nil {
			return fmt.Errorf("invalid file %s: %w", file.OriginalPath, err)
		}
	}

	return nil
}

// Validate validates a backup file
func (f BackupFile) Validate() error {
	if f.OriginalPath == "" {
		return fmt.Errorf("%w: original path is empty", ErrInvalidFile)
	}

	if f.BackupPath == "" {
		return fmt.Errorf("%w: backup path is empty", ErrInvalidFile)
	}

	if f.SizeBytes < 0 {
		return fmt.Errorf("%w: size cannot be negative", ErrInvalidFile)
	}

	return nil
}
