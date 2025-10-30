package backup

import "context"

// Repository defines the interface for backup persistence
type Repository interface {
	// Save persists a backup
	Save(ctx context.Context, backup *Backup) error

	// FindByID retrieves a backup by its ID
	FindByID(ctx context.Context, id string) (*Backup, error)

	// FindAll retrieves all backups sorted by creation date (newest first)
	FindAll(ctx context.Context) ([]*Backup, error)

	// Delete removes a backup
	Delete(ctx context.Context, id string) error

	// Exists checks if a backup with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)
}
