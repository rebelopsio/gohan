package usecases

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// ListInstallationsUseCase retrieves all installation sessions
type ListInstallationsUseCase struct {
	sessionRepo installation.InstallationSessionRepository
}

// NewListInstallationsUseCase creates a new ListInstallationsUseCase
func NewListInstallationsUseCase(sessionRepo installation.InstallationSessionRepository) *ListInstallationsUseCase {
	return &ListInstallationsUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute retrieves all installation sessions and returns their summaries
func (u *ListInstallationsUseCase) Execute(ctx context.Context) (*dto.ListInstallationsResponse, error) {
	// Retrieve all sessions from repository
	sessions, err := u.sessionRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// Build response
	summaries := make([]dto.InstallationSessionSummary, 0, len(sessions))
	for _, session := range sessions {
		summary := buildSessionSummary(session)
		summaries = append(summaries, summary)
	}

	return &dto.ListInstallationsResponse{
		Sessions:   summaries,
		TotalCount: len(summaries),
	}, nil
}

// buildSessionSummary creates a session summary from a domain session
func buildSessionSummary(session *installation.InstallationSession) dto.InstallationSessionSummary {
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

	// Format timestamps
	startedAt := ""
	if !session.StartedAt().IsZero() {
		startedAt = session.StartedAt().Format("2006-01-02T15:04:05Z07:00")
	}

	completedAt := ""
	if !session.CompletedAt().IsZero() {
		completedAt = session.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
	}

	return dto.InstallationSessionSummary{
		SessionID:           session.ID(),
		Status:              string(session.Status()),
		CurrentPhase:        currentPhase,
		PercentComplete:     percentComplete,
		ComponentsInstalled: len(installedComponents),
		ComponentsTotal:     componentsTotal,
		StartedAt:           startedAt,
		CompletedAt:         completedAt,
	}
}
