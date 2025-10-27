package usecases

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// GetInstallationStatusUseCase retrieves the status of an installation session
type GetInstallationStatusUseCase struct {
	sessionRepo installation.InstallationSessionRepository
}

// NewGetInstallationStatusUseCase creates a new GetInstallationStatusUseCase
func NewGetInstallationStatusUseCase(sessionRepo installation.InstallationSessionRepository) *GetInstallationStatusUseCase {
	return &GetInstallationStatusUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute retrieves the installation status for a given session ID
func (u *GetInstallationStatusUseCase) Execute(ctx context.Context, sessionID string) (*dto.InstallationProgressResponse, error) {
	// Retrieve session from repository
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	// Build progress response
	config := session.Configuration()
	installedComponents := session.InstalledComponents()
	componentsTotal := len(config.Components())

	// Calculate progress percentage
	percentComplete := 0
	if componentsTotal > 0 {
		percentComplete = (len(installedComponents) * 100) / componentsTotal
	}

	// Determine current phase based on status
	currentPhase := "pending"
	switch session.Status() {
	case installation.StatusPending:
		currentPhase = "pending"
	case installation.StatusPreparation:
		currentPhase = "preparation"
	case installation.StatusInstalling:
		currentPhase = "installing"
	case installation.StatusConfiguring:
		currentPhase = "configuring"
	case installation.StatusVerifying:
		currentPhase = "verifying"
	case installation.StatusCompleted:
		currentPhase = "completed"
		percentComplete = 100
	case installation.StatusFailed:
		currentPhase = "failed"
	}

	// Get message (failure reason or empty)
	message := session.FailureReason()

	// Estimate remaining time (simplified)
	estimatedRemaining := "0s"
	if session.IsInProgress() && !session.StartedAt().IsZero() {
		elapsed := session.Duration()
		if percentComplete > 0 && percentComplete < 100 {
			// Calculate total estimated time based on current progress
			elapsedNs := int64(elapsed)
			totalEstimatedNs := elapsedNs * 100 / int64(percentComplete)
			remainingNs := totalEstimatedNs - elapsedNs
			estimatedRemaining = fmt.Sprintf("%v", remainingNs)
		}
	}

	return &dto.InstallationProgressResponse{
		SessionID:            session.ID(),
		Status:               string(session.Status()),
		CurrentPhase:         currentPhase,
		PercentComplete:      percentComplete,
		ComponentsInstalled:  len(installedComponents),
		ComponentsTotal:      componentsTotal,
		EstimatedRemaining:   estimatedRemaining,
		Message:              message,
	}, nil
}
