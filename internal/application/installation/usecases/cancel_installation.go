package usecases

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// CancelInstallationUseCase cancels an in-progress installation
type CancelInstallationUseCase struct {
	sessionRepo installation.InstallationSessionRepository
}

// NewCancelInstallationUseCase creates a new CancelInstallationUseCase
func NewCancelInstallationUseCase(sessionRepo installation.InstallationSessionRepository) *CancelInstallationUseCase {
	return &CancelInstallationUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute cancels the installation session with the given ID
func (u *CancelInstallationUseCase) Execute(ctx context.Context, sessionID string) error {
	// Retrieve session from repository
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	// Cancel the session (mark as failed with cancellation reason)
	if err := session.Fail("installation cancelled by user"); err != nil {
		return fmt.Errorf("failed to cancel installation: %w", err)
	}

	// Save the updated session
	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return fmt.Errorf("failed to save cancelled session: %w", err)
	}

	return nil
}
