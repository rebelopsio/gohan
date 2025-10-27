package history

import (
	"strings"
	"time"
)

// FailureDetails is a value object containing information about installation failures
type FailureDetails struct {
	reason    string
	failedAt  time.Time
	phase     string
	errorCode string
}

// NewFailureDetails creates failure details value object
func NewFailureDetails(
	reason string,
	failedAt time.Time,
	phase string,
	errorCode string,
) (FailureDetails, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return FailureDetails{}, ErrInvalidFailureReason
	}

	if failedAt.IsZero() {
		return FailureDetails{}, ErrInvalidTimestamp
	}

	phase = strings.TrimSpace(phase)
	errorCode = strings.TrimSpace(errorCode)

	return FailureDetails{
		reason:    reason,
		failedAt:  failedAt,
		phase:     phase,
		errorCode: errorCode,
	}, nil
}

// Reason returns the failure reason
func (f FailureDetails) Reason() string {
	return f.reason
}

// FailedAt returns when the failure occurred
func (f FailureDetails) FailedAt() time.Time {
	return f.failedAt
}

// Phase returns which phase failed
func (f FailureDetails) Phase() string {
	return f.phase
}

// ErrorCode returns the error code if available
func (f FailureDetails) ErrorCode() string {
	return f.errorCode
}

// HasErrorCode returns true if an error code is set
func (f FailureDetails) HasErrorCode() bool {
	return f.errorCode != ""
}
