package verification

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CheckResult represents the outcome of a single verification check
type CheckResult struct {
	id          string
	component   ComponentName
	status      CheckStatus
	severity    CheckSeverity
	message     string
	details     []string
	suggestions []string
	checkedAt   time.Time
}

// NewCheckResult creates a new check result
func NewCheckResult(
	component ComponentName,
	status CheckStatus,
	severity CheckSeverity,
	message string,
	details []string,
	suggestions []string,
) CheckResult {
	// Ensure non-nil slices
	if details == nil {
		details = []string{}
	}
	if suggestions == nil {
		suggestions = []string{}
	}

	return CheckResult{
		id:          uuid.New().String(),
		component:   component,
		status:      status,
		severity:    severity,
		message:     message,
		details:     details,
		suggestions: suggestions,
		checkedAt:   time.Now(),
	}
}

// ID returns the result identifier
func (r CheckResult) ID() string {
	return r.id
}

// Component returns the component that was checked
func (r CheckResult) Component() ComponentName {
	return r.component
}

// Status returns the check status
func (r CheckResult) Status() CheckStatus {
	return r.status
}

// Severity returns the severity level
func (r CheckResult) Severity() CheckSeverity {
	return r.severity
}

// Message returns the main status message
func (r CheckResult) Message() string {
	return r.message
}

// Details returns additional details
func (r CheckResult) Details() []string {
	if r.details == nil {
		return nil
	}
	result := make([]string, len(r.details))
	copy(result, r.details)
	return result
}

// Suggestions returns fix suggestions
func (r CheckResult) Suggestions() []string {
	if r.suggestions == nil {
		return nil
	}
	result := make([]string, len(r.suggestions))
	copy(result, r.suggestions)
	return result
}

// CheckedAt returns when this check was performed
func (r CheckResult) CheckedAt() time.Time {
	return r.checkedAt
}

// IsPassing returns true if check passed
func (r CheckResult) IsPassing() bool {
	return r.status == StatusPass
}

// IsCritical returns true if this is a critical failure
func (r CheckResult) IsCritical() bool {
	return r.severity == SeverityCritical && r.status != StatusPass
}

// IsWarning returns true if this is a warning
func (r CheckResult) IsWarning() bool {
	return r.status == StatusWarning
}

// FormatMessage returns a formatted message with details
func (r CheckResult) FormatMessage() string {
	msg := fmt.Sprintf("%s: %s", r.component, r.message)

	if len(r.details) > 0 {
		msg += "\nDetails:"
		for _, detail := range r.details {
			msg += fmt.Sprintf("\n  - %s", detail)
		}
	}

	if len(r.suggestions) > 0 {
		msg += "\nSuggestions:"
		for _, suggestion := range r.suggestions {
			msg += fmt.Sprintf("\n  - %s", suggestion)
		}
	}

	return msg
}
