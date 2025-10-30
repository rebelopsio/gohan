package verification

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// VerificationReport is the aggregate root for system verification
type VerificationReport struct {
	mu            sync.RWMutex
	id            string
	startedAt     time.Time
	completedAt   time.Time
	results       []CheckResult
	overallStatus CheckStatus
}

// NewVerificationReport creates a new verification report
func NewVerificationReport() *VerificationReport {
	return &VerificationReport{
		id:        uuid.New().String(),
		startedAt: time.Now(),
		results:   make([]CheckResult, 0),
	}
}

// ID returns the report identifier
func (r *VerificationReport) ID() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.id
}

// AddResult adds a check result to the report
func (r *VerificationReport) AddResult(result CheckResult) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results = append(r.results, result)
}

// Complete marks the report as complete
func (r *VerificationReport) Complete() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.completedAt = time.Now()
	r.calculateOverallStatus()
}

// Results returns all check results
func (r *VerificationReport) Results() []CheckResult {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]CheckResult, len(r.results))
	copy(results, r.results)
	return results
}

// OverallStatus returns the overall verification status
func (r *VerificationReport) OverallStatus() CheckStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.overallStatus
}

// HasCriticalFailures returns true if any critical checks failed
func (r *VerificationReport) HasCriticalFailures() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, result := range r.results {
		if result.IsCritical() {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any checks produced warnings
func (r *VerificationReport) HasWarnings() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, result := range r.results {
		if result.IsWarning() {
			return true
		}
	}
	return false
}

// PassedCount returns the number of passed checks
func (r *VerificationReport) PassedCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, result := range r.results {
		if result.IsPassing() {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warnings
func (r *VerificationReport) WarningCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, result := range r.results {
		if result.IsWarning() {
			count++
		}
	}
	return count
}

// FailedCount returns the number of failed checks
func (r *VerificationReport) FailedCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, result := range r.results {
		if result.Status() == StatusFail {
			count++
		}
	}
	return count
}

// CriticalCount returns the number of critical failures
func (r *VerificationReport) CriticalCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, result := range r.results {
		if result.IsCritical() {
			count++
		}
	}
	return count
}

// Duration returns how long the verification took
func (r *VerificationReport) Duration() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.completedAt.IsZero() {
		return time.Since(r.startedAt)
	}
	return r.completedAt.Sub(r.startedAt)
}

// calculateOverallStatus determines the overall status based on results
func (r *VerificationReport) calculateOverallStatus() {
	// Already locked by Complete()

	hasFailed := false
	hasWarning := false

	for _, result := range r.results {
		if result.Status() == StatusFail {
			hasFailed = true
		} else if result.Status() == StatusWarning {
			hasWarning = true
		}
	}

	if hasFailed {
		r.overallStatus = StatusFail
	} else if hasWarning {
		r.overallStatus = StatusWarning
	} else {
		r.overallStatus = StatusPass
	}
}
