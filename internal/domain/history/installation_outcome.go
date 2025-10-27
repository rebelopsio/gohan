package history

// InstallationOutcome represents the final result of an installation
type InstallationOutcome string

const (
	OutcomeSuccess    InstallationOutcome = "success"
	OutcomeFailed     InstallationOutcome = "failed"
	OutcomeRolledBack InstallationOutcome = "rolled_back"
)

// NewInstallationOutcome creates an outcome value object
func NewInstallationOutcome(outcome string) (InstallationOutcome, error) {
	switch InstallationOutcome(outcome) {
	case OutcomeSuccess, OutcomeFailed, OutcomeRolledBack:
		return InstallationOutcome(outcome), nil
	default:
		return "", ErrInvalidOutcome
	}
}

// String returns the string representation
func (o InstallationOutcome) String() string {
	return string(o)
}

// IsSuccessful returns true if installation succeeded
func (o InstallationOutcome) IsSuccessful() bool {
	return o == OutcomeSuccess
}

// IsFailed returns true if installation failed
func (o InstallationOutcome) IsFailed() bool {
	return o == OutcomeFailed
}

// IsRolledBack returns true if installation was rolled back
func (o InstallationOutcome) IsRolledBack() bool {
	return o == OutcomeRolledBack
}

// Equals checks if two outcomes are equal
func (o InstallationOutcome) Equals(other InstallationOutcome) bool {
	return o == other
}
