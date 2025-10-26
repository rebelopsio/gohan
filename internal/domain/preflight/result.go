package preflight

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ValidationResult represents the outcome of a single validation check
type ValidationResult struct {
	id              string
	requirementName RequirementName
	status          ValidationStatus
	severity        Severity
	actualValue     interface{}
	expectedValue   interface{}
	guidance        UserGuidance
	detectedAt      time.Time
}

// NewValidationResult creates a new validation result
func NewValidationResult(
	requirementName RequirementName,
	status ValidationStatus,
	severity Severity,
	actualValue interface{},
	expectedValue interface{},
	guidance UserGuidance,
) ValidationResult {
	return ValidationResult{
		id:              uuid.New().String(),
		requirementName: requirementName,
		status:          status,
		severity:        severity,
		actualValue:     actualValue,
		expectedValue:   expectedValue,
		guidance:        guidance,
		detectedAt:      time.Now(),
	}
}

// ID returns the result identifier
func (r ValidationResult) ID() string {
	return r.id
}

// RequirementName returns the requirement being validated
func (r ValidationResult) RequirementName() RequirementName {
	return r.requirementName
}

// Status returns the validation status
func (r ValidationResult) Status() ValidationStatus {
	return r.status
}

// Severity returns the severity level
func (r ValidationResult) Severity() Severity {
	return r.severity
}

// ActualValue returns what was detected
func (r ValidationResult) ActualValue() interface{} {
	return r.actualValue
}

// ExpectedValue returns what was expected
func (r ValidationResult) ExpectedValue() interface{} {
	return r.expectedValue
}

// Guidance returns user guidance for failures
func (r ValidationResult) Guidance() UserGuidance {
	return r.guidance
}

// DetectedAt returns when this was validated
func (r ValidationResult) DetectedAt() time.Time {
	return r.detectedAt
}

// IsBlocking returns true if this failure blocks installation
func (r ValidationResult) IsBlocking() bool {
	return r.status == StatusFail &&
		(r.severity == SeverityCritical || r.severity == SeverityHigh)
}

// IsWarning returns true if this is a warning-level result
func (r ValidationResult) IsWarning() bool {
	return r.status == StatusWarning ||
		(r.status == StatusFail && (r.severity == SeverityMedium || r.severity == SeverityLow))
}

// IsPassing returns true if validation passed
func (r ValidationResult) IsPassing() bool {
	return r.status == StatusPass
}

// FormatMessage returns a human-readable message
func (r ValidationResult) FormatMessage() string {
	switch r.status {
	case StatusPass:
		return fmt.Sprintf("✓ %s: Valid", r.requirementName)
	case StatusFail:
		return fmt.Sprintf("✗ %s: %s", r.requirementName, r.guidance.Message())
	case StatusWarning:
		return fmt.Sprintf("⚠ %s: %s", r.requirementName, r.guidance.Message())
	default:
		return fmt.Sprintf("? %s: Unknown status", r.requirementName)
	}
}
