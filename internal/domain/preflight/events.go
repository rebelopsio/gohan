package preflight

import (
	"time"
)

// ValidationStartedEvent fires when validation begins
type ValidationStartedEvent struct {
	SessionID  string
	occurredAt time.Time
}

// NewValidationStartedEvent creates a new validation started event
func NewValidationStartedEvent(sessionID string) ValidationStartedEvent {
	return ValidationStartedEvent{
		SessionID:  sessionID,
		occurredAt: time.Now(),
	}
}

// OccurredAt returns when the event occurred
func (e ValidationStartedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type identifier
func (e ValidationStartedEvent) EventType() string {
	return "validation.started"
}

// ValidationCompletedEvent fires when all validations finish
type ValidationCompletedEvent struct {
	SessionID  string
	Outcome    ValidationOutcome
	Duration   time.Duration
	occurredAt time.Time
}

// NewValidationCompletedEvent creates a new validation completed event
func NewValidationCompletedEvent(
	sessionID string,
	outcome ValidationOutcome,
	duration time.Duration,
) ValidationCompletedEvent {
	return ValidationCompletedEvent{
		SessionID:  sessionID,
		Outcome:    outcome,
		Duration:   duration,
		occurredAt: time.Now(),
	}
}

// OccurredAt returns when the event occurred
func (e ValidationCompletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type identifier
func (e ValidationCompletedEvent) EventType() string {
	return "validation.completed"
}

// ValidationBlockedEvent fires when blocking failures prevent installation
type ValidationBlockedEvent struct {
	SessionID      string
	BlockingCount  int
	BlockedReasons []string
	occurredAt     time.Time
}

// NewValidationBlockedEvent creates a new validation blocked event
func NewValidationBlockedEvent(
	sessionID string,
	blockingCount int,
	blockedReasons []string,
) ValidationBlockedEvent {
	return ValidationBlockedEvent{
		SessionID:      sessionID,
		BlockingCount:  blockingCount,
		BlockedReasons: blockedReasons,
		occurredAt:     time.Now(),
	}
}

// OccurredAt returns when the event occurred
func (e ValidationBlockedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type identifier
func (e ValidationBlockedEvent) EventType() string {
	return "validation.blocked"
}

// ValidationWarningEvent fires when warnings exist but installation can proceed
type ValidationWarningEvent struct {
	SessionID    string
	WarningCount int
	Warnings     []string
	occurredAt   time.Time
}

// NewValidationWarningEvent creates a new validation warning event
func NewValidationWarningEvent(
	sessionID string,
	warningCount int,
	warnings []string,
) ValidationWarningEvent {
	return ValidationWarningEvent{
		SessionID:    sessionID,
		WarningCount: warningCount,
		Warnings:     warnings,
		occurredAt:   time.Now(),
	}
}

// OccurredAt returns when the event occurred
func (e ValidationWarningEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type identifier
func (e ValidationWarningEvent) EventType() string {
	return "validation.warning"
}

// RequirementFailedEvent fires when a specific requirement fails
type RequirementFailedEvent struct {
	SessionID       string
	RequirementName RequirementName
	Reason          string
	IsBlocking      bool
	occurredAt      time.Time
}

// NewRequirementFailedEvent creates a new requirement failed event
func NewRequirementFailedEvent(
	sessionID string,
	requirementName RequirementName,
	reason string,
	isBlocking bool,
) RequirementFailedEvent {
	return RequirementFailedEvent{
		SessionID:       sessionID,
		RequirementName: requirementName,
		Reason:          reason,
		IsBlocking:      isBlocking,
		occurredAt:      time.Now(),
	}
}

// OccurredAt returns when the event occurred
func (e RequirementFailedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type identifier
func (e RequirementFailedEvent) EventType() string {
	return "requirement.failed"
}
