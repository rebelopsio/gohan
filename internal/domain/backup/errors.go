package backup

import "errors"

var (
	// ErrDuplicateFile occurs when attempting to add a file that's already in the backup
	ErrDuplicateFile = errors.New("file already exists in backup")

	// ErrBackupNotFound occurs when a backup cannot be found
	ErrBackupNotFound = errors.New("backup not found")

	// ErrInvalidBackupID occurs when a backup ID is invalid
	ErrInvalidBackupID = errors.New("invalid backup ID")

	// ErrInvalidFile occurs when a backup file has invalid data
	ErrInvalidFile = errors.New("invalid backup file")

	// ErrBackupCorrupted occurs when a backup is corrupted or incomplete
	ErrBackupCorrupted = errors.New("backup is corrupted")

	// ErrInsufficientSpace occurs when there's not enough disk space for backup
	ErrInsufficientSpace = errors.New("insufficient disk space")

	// ErrEmptyBackup occurs when attempting to complete a backup with no files
	ErrEmptyBackup = errors.New("backup contains no files")
)
