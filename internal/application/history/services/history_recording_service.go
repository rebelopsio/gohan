package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// HistoryRecordingService records installation sessions to history
type HistoryRecordingService struct {
	historyRepo history.Repository
}

// NewHistoryRecordingService creates a new history recording service
func NewHistoryRecordingService(historyRepo history.Repository) *HistoryRecordingService {
	return &HistoryRecordingService{
		historyRepo: historyRepo,
	}
}

// RecordInstallation creates a history record from a completed installation session
func (s *HistoryRecordingService) RecordInstallation(
	ctx context.Context,
	session *installation.InstallationSession,
) (history.RecordID, error) {
	// Validate session is in terminal state
	if !session.IsCompleted() && !session.IsFailed() {
		return history.RecordID{}, fmt.Errorf("cannot record incomplete session: status is %s", session.Status())
	}

	// Determine outcome
	var outcome history.InstallationOutcome
	var err error
	if session.IsCompleted() {
		outcome, err = history.NewInstallationOutcome("success")
	} else if session.IsFailed() {
		outcome, err = history.NewInstallationOutcome("failed")
	} else {
		return history.RecordID{}, fmt.Errorf("unexpected session state")
	}
	if err != nil {
		return history.RecordID{}, fmt.Errorf("failed to create outcome: %w", err)
	}

	// Build installed packages list
	var installedPackages []history.InstalledPackage
	for _, comp := range session.InstalledComponents() {
		// Determine size
		var sizeBytes uint64 = 1024 // Default 1KB if not available
		if comp.PackageInfo() != nil {
			sizeBytes = comp.PackageInfo().SizeBytes()
		}

		pkg, err := history.NewInstalledPackage(
			string(comp.Component()),
			comp.Version(),
			sizeBytes,
		)
		if err != nil {
			return history.RecordID{}, fmt.Errorf("failed to create installed package: %w", err)
		}
		installedPackages = append(installedPackages, pkg)
	}

	// Get installation times
	installedAt := session.StartedAt()
	completedAt := session.CompletedAt()

	// Determine target package and version from first component
	config := session.Configuration()
	components := config.Components()
	if len(components) == 0 {
		return history.RecordID{}, fmt.Errorf("session has no components")
	}

	// Use first component as the "target" for history
	targetPackage := string(components[0].Component())
	targetVersion := components[0].Version()

	// Build metadata
	metadata, err := history.NewInstallationMetadata(
		targetPackage,
		targetVersion,
		installedAt,
		completedAt,
		installedPackages,
	)
	if err != nil {
		return history.RecordID{}, fmt.Errorf("failed to create metadata: %w", err)
	}

	// Capture system context
	systemContext, err := s.captureSystemContext()
	if err != nil {
		return history.RecordID{}, fmt.Errorf("failed to capture system context: %w", err)
	}

	// Build failure details if failed
	var failureDetails *history.FailureDetails
	if session.IsFailed() && session.FailureReason() != "" {
		fd, err := history.NewFailureDetails(
			session.FailureReason(),
			completedAt,
			session.Status().String(),
			"", // No error code available from session
		)
		if err != nil {
			return history.RecordID{}, fmt.Errorf("failed to create failure details: %w", err)
		}
		failureDetails = &fd
	}

	// Create installation record
	record, err := history.NewInstallationRecord(
		session.ID(),
		outcome,
		metadata,
		systemContext,
		failureDetails,
		time.Now(),
	)
	if err != nil {
		return history.RecordID{}, fmt.Errorf("failed to create installation record: %w", err)
	}

	// Save to repository
	if err := s.historyRepo.Save(ctx, record); err != nil {
		return history.RecordID{}, fmt.Errorf("failed to save installation record: %w", err)
	}

	return record.ID(), nil
}

// captureSystemContext captures current system information
func (s *HistoryRecordingService) captureSystemContext() (history.SystemContext, error) {
	// Get OS version from /etc/os-release or similar
	osVersion := s.detectOSVersion()

	// Get kernel version
	kernelVersion := s.detectKernelVersion()

	// Get hostname
	hostname, _ := os.Hostname()

	// Create system context
	return history.NewSystemContext(
		osVersion,
		kernelVersion,
		"1.0.0", // Gohan version - could be injected
		hostname,
	)
}

// detectOSVersion attempts to detect OS version
func (s *HistoryRecordingService) detectOSVersion() string {
	// Try to read /etc/os-release
	// For now, return a default value
	// In production, this would read from the filesystem
	return "Debian GNU/Linux"
}

// detectKernelVersion attempts to detect kernel version
func (s *HistoryRecordingService) detectKernelVersion() string {
	// Try to get kernel version from uname
	// For now, return empty (optional field)
	// In production, this would exec `uname -r`
	return ""
}
