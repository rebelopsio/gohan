package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	domainRepo "github.com/rebelopsio/gohan/internal/domain/repository"
)

// CheckRepositoryRequest contains parameters for checking repository configuration
type CheckRepositoryRequest struct {
	SourcesListPath string // Path to sources.list file
}

// CheckRepositoryResponse contains the current repository configuration status
type CheckRepositoryResponse struct {
	HasMain              bool
	HasContrib           bool
	HasNonFree           bool
	HasNonFreeFirmware   bool
	HasDebSrc            bool
	TotalEntries         int
	DebEntries           int
	DebSrcEntries        int
	SourcesListContent   string
}

// EnableNonFreeRequest contains parameters for enabling non-free repositories
type EnableNonFreeRequest struct {
	SourcesListPath string
	BackupFirst     bool
}

// EnableNonFreeResponse contains the result of enabling non-free
type EnableNonFreeResponse struct {
	Modified         bool
	BackupPath       string
	ComponentsAdded  []string
}

// EnableDebSrcRequest contains parameters for enabling deb-src
type EnableDebSrcRequest struct {
	SourcesListPath string
	BackupFirst     bool
}

// EnableDebSrcResponse contains the result of enabling deb-src
type EnableDebSrcResponse struct {
	Modified       bool
	BackupPath     string
	EntriesAdded   int
}

// BackupSourcesListRequest contains parameters for backing up sources.list
type BackupSourcesListRequest struct {
	SourcesListPath string
	BackupDir       string
}

// BackupSourcesListResponse contains information about the backup
type BackupSourcesListResponse struct {
	BackupPath string
	Timestamp  time.Time
	Size       int64
}

// SourcesListManager is the interface for managing sources.list files
type SourcesListManager interface {
	ReadConfig(path string) (*domainRepo.RepositoryConfig, error)
	WriteConfig(path string, config *domainRepo.RepositoryConfig) error
	Backup(path string, backupDir string) (string, error)
	Exists(path string) (bool, error)
}

// CheckRepositoryUseCase handles checking repository configuration
type CheckRepositoryUseCase struct {
	manager SourcesListManager
}

// NewCheckRepositoryUseCase creates a new use case instance
func NewCheckRepositoryUseCase(manager SourcesListManager) *CheckRepositoryUseCase {
	return &CheckRepositoryUseCase{
		manager: manager,
	}
}

// Execute checks the current repository configuration
func (uc *CheckRepositoryUseCase) Execute(ctx context.Context, req CheckRepositoryRequest) (*CheckRepositoryResponse, error) {
	// Read repository config
	config, err := uc.manager.ReadConfig(req.SourcesListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository config: %w", err)
	}

	// Get entries
	entries := config.Entries()

	// Count entry types
	debCount := 0
	debSrcCount := 0
	for _, entry := range entries {
		if entry.Type == "deb" {
			debCount++
		} else if entry.Type == "deb-src" {
			debSrcCount++
		}
	}

	// Build response
	return &CheckRepositoryResponse{
		HasMain:            config.HasComponent("main"),
		HasContrib:         config.HasComponent("contrib"),
		HasNonFree:         config.HasComponent("non-free"),
		HasNonFreeFirmware: config.HasComponent("non-free-firmware"),
		HasDebSrc:          config.HasDebSrc(),
		TotalEntries:       len(entries),
		DebEntries:         debCount,
		DebSrcEntries:      debSrcCount,
		SourcesListContent: config.String(),
	}, nil
}

// EnableNonFreeUseCase handles enabling non-free repositories
type EnableNonFreeUseCase struct {
	manager SourcesListManager
}

// NewEnableNonFreeUseCase creates a new use case instance
func NewEnableNonFreeUseCase(manager SourcesListManager) *EnableNonFreeUseCase {
	return &EnableNonFreeUseCase{
		manager: manager,
	}
}

// Execute enables non-free and non-free-firmware components
func (uc *EnableNonFreeUseCase) Execute(ctx context.Context, req EnableNonFreeRequest) (*EnableNonFreeResponse, error) {
	response := &EnableNonFreeResponse{
		ComponentsAdded: []string{},
	}

	// Backup if requested
	if req.BackupFirst {
		backupPath, err := uc.manager.Backup(req.SourcesListPath, filepath.Dir(req.SourcesListPath))
		if err != nil {
			return nil, fmt.Errorf("failed to backup sources.list: %w", err)
		}
		response.BackupPath = backupPath
	}

	// Read repository config
	config, err := uc.manager.ReadConfig(req.SourcesListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository config: %w", err)
	}

	// Check what's missing
	needsNonFree := !config.HasComponent("non-free")
	needsNonFreeFirmware := !config.HasComponent("non-free-firmware")

	if !needsNonFree && !needsNonFreeFirmware {
		// Already has everything
		return response, nil
	}

	// Add components
	if needsNonFree {
		config.AddComponent("non-free")
		response.ComponentsAdded = append(response.ComponentsAdded, "non-free")
		response.Modified = true
	}

	if needsNonFreeFirmware {
		config.AddComponent("non-free-firmware")
		response.ComponentsAdded = append(response.ComponentsAdded, "non-free-firmware")
		response.Modified = true
	}

	// Write updated sources.list
	if response.Modified {
		if err := uc.manager.WriteConfig(req.SourcesListPath, config); err != nil {
			return nil, fmt.Errorf("failed to write sources.list: %w", err)
		}
	}

	return response, nil
}

// EnableDebSrcUseCase handles enabling deb-src entries
type EnableDebSrcUseCase struct {
	manager SourcesListManager
}

// NewEnableDebSrcUseCase creates a new use case instance
func NewEnableDebSrcUseCase(manager SourcesListManager) *EnableDebSrcUseCase {
	return &EnableDebSrcUseCase{
		manager: manager,
	}
}

// Execute enables deb-src entries for all deb entries
func (uc *EnableDebSrcUseCase) Execute(ctx context.Context, req EnableDebSrcRequest) (*EnableDebSrcResponse, error) {
	response := &EnableDebSrcResponse{}

	// Backup if requested
	if req.BackupFirst {
		backupPath, err := uc.manager.Backup(req.SourcesListPath, filepath.Dir(req.SourcesListPath))
		if err != nil {
			return nil, fmt.Errorf("failed to backup sources.list: %w", err)
		}
		response.BackupPath = backupPath
	}

	// Read repository config
	config, err := uc.manager.ReadConfig(req.SourcesListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository config: %w", err)
	}

	// Check if deb-src already enabled
	if config.HasDebSrc() {
		return response, nil
	}

	// Count entries before
	entriesBefore := len(config.Entries())

	// Enable deb-src
	config.EnableDebSrc()

	// Count entries after
	entriesAfter := len(config.Entries())
	response.EntriesAdded = entriesAfter - entriesBefore
	response.Modified = response.EntriesAdded > 0

	// Write updated sources.list
	if response.Modified {
		if err := uc.manager.WriteConfig(req.SourcesListPath, config); err != nil {
			return nil, fmt.Errorf("failed to write sources.list: %w", err)
		}
	}

	return response, nil
}

// BackupSourcesListUseCase handles backing up sources.list
type BackupSourcesListUseCase struct {
	manager SourcesListManager
}

// NewBackupSourcesListUseCase creates a new use case instance
func NewBackupSourcesListUseCase(manager SourcesListManager) *BackupSourcesListUseCase {
	return &BackupSourcesListUseCase{
		manager: manager,
	}
}

// Execute creates a backup of sources.list
func (uc *BackupSourcesListUseCase) Execute(ctx context.Context, req BackupSourcesListRequest) (*BackupSourcesListResponse, error) {
	// Create backup
	backupPath, err := uc.manager.Backup(req.SourcesListPath, req.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Get backup file info
	info, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	return &BackupSourcesListResponse{
		BackupPath: backupPath,
		Timestamp:  info.ModTime(),
		Size:       info.Size(),
	}, nil
}
