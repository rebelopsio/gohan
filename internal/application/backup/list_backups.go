package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/backup"
)

// BackupSummary contains summary information about a backup
type BackupSummary struct {
	ID          string // Backup ID
	Description string // Backup description
	CreatedAt   string // Creation timestamp (formatted)
	Age         string // Age in human-readable format
	FileCount   int    // Number of files
	TotalSize   string // Total size in human-readable format
	Status      string // Backup status
}

// ListBackupsResponse contains the list of backups
type ListBackupsResponse struct {
	Backups []*BackupSummary
	Total   int
}

// ListBackupsUseCase handles listing available backups
type ListBackupsUseCase struct {
	repository backup.Repository
}

// NewListBackupsUseCase creates a new use case instance
func NewListBackupsUseCase(repository backup.Repository) *ListBackupsUseCase {
	return &ListBackupsUseCase{
		repository: repository,
	}
}

// Execute retrieves all backups
func (uc *ListBackupsUseCase) Execute(ctx context.Context) (*ListBackupsResponse, error) {
	// Get all backups
	backups, err := uc.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	// Convert to summaries
	summaries := make([]*BackupSummary, 0, len(backups))
	for _, b := range backups {
		summaries = append(summaries, &BackupSummary{
			ID:          b.ID(),
			Description: b.Description(),
			CreatedAt:   b.CreatedAt().Format("2006-01-02 15:04:05"),
			Age:         formatDuration(b.Age()),
			FileCount:   b.FileCount(),
			TotalSize:   formatBytes(b.TotalSize()),
			Status:      string(b.Status()),
		})
	}

	return &ListBackupsResponse{
		Backups: summaries,
		Total:   len(summaries),
	}, nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	if d.Hours() >= 1 {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if d.Minutes() >= 1 {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return "just now"
}
