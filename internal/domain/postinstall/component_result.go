package postinstall

import (
	"time"

	"github.com/google/uuid"
)

// ComponentResult represents the result of setting up a component
type ComponentResult struct {
	id          string
	component   ComponentType
	status      SetupStatus
	message     string
	details     []string
	error       error
	startedAt   time.Time
	completedAt time.Time
}

// NewComponentResult creates a new component result
func NewComponentResult(component ComponentType, status SetupStatus, message string) ComponentResult {
	return ComponentResult{
		id:          uuid.New().String(),
		component:   component,
		status:      status,
		message:     message,
		details:     []string{},
		startedAt:   time.Now(),
		completedAt: time.Time{},
	}
}

// NewComponentResultWithError creates a failed component result with error
func NewComponentResultWithError(component ComponentType, message string, err error) ComponentResult {
	return ComponentResult{
		id:          uuid.New().String(),
		component:   component,
		status:      StatusFailed,
		message:     message,
		details:     []string{},
		error:       err,
		startedAt:   time.Now(),
		completedAt: time.Now(),
	}
}

// ID returns the result ID
func (r ComponentResult) ID() string {
	return r.id
}

// Component returns the component type
func (r ComponentResult) Component() ComponentType {
	return r.component
}

// Status returns the setup status
func (r ComponentResult) Status() SetupStatus {
	return r.status
}

// Message returns the result message
func (r ComponentResult) Message() string {
	return r.message
}

// Details returns additional details
func (r ComponentResult) Details() []string {
	result := make([]string, len(r.details))
	copy(result, r.details)
	return result
}

// Error returns any error that occurred
func (r ComponentResult) Error() error {
	return r.error
}

// StartedAt returns when the setup started
func (r ComponentResult) StartedAt() time.Time {
	return r.startedAt
}

// CompletedAt returns when the setup completed
func (r ComponentResult) CompletedAt() time.Time {
	return r.completedAt
}

// Duration returns how long the setup took
func (r ComponentResult) Duration() time.Duration {
	if r.completedAt.IsZero() {
		return time.Since(r.startedAt)
	}
	return r.completedAt.Sub(r.startedAt)
}

// IsSuccess returns true if setup succeeded
func (r ComponentResult) IsSuccess() bool {
	return r.status.IsSuccess()
}

// IsFailure returns true if setup failed
func (r ComponentResult) IsFailure() bool {
	return r.status.IsFailure()
}

// WithDetails adds details to the result (creates new instance)
func (r ComponentResult) WithDetails(details ...string) ComponentResult {
	r.details = append(r.details, details...)
	return r
}

// Complete marks the result as completed
func (r ComponentResult) Complete(status SetupStatus) ComponentResult {
	r.status = status
	r.completedAt = time.Now()
	return r
}
