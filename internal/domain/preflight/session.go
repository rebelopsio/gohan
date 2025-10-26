package preflight

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ValidationSession is the aggregate root for all preflight validations
type ValidationSession struct {
	mu            sync.RWMutex
	id            string
	startedAt     time.Time
	completedAt   time.Time
	overallResult ValidationOutcome
	results       []ValidationResult
}

// NewValidationSession creates a new validation session
func NewValidationSession() *ValidationSession {
	return &ValidationSession{
		id:        uuid.New().String(),
		startedAt: time.Now(),
		results:   make([]ValidationResult, 0),
	}
}

// ID returns the session identifier
func (s *ValidationSession) ID() string {
	return s.id
}

// StartedAt returns when validation began
func (s *ValidationSession) StartedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.startedAt
}

// CompletedAt returns when validation finished
func (s *ValidationSession) CompletedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.completedAt
}

// OverallResult returns the aggregate outcome
func (s *ValidationSession) OverallResult() ValidationOutcome {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.overallResult
}

// Results returns all validation results
func (s *ValidationSession) Results() []ValidationResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.results == nil {
		return nil
	}
	result := make([]ValidationResult, len(s.results))
	copy(result, s.results)
	return result
}

// AddResult adds a validation result to the session
func (s *ValidationSession) AddResult(result ValidationResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results = append(s.results, result)
	s.recalculateOutcome()
}

// Complete marks the session as finished
func (s *ValidationSession) Complete() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make Complete idempotent - only set completedAt if not already set
	if s.completedAt.IsZero() {
		s.completedAt = time.Now()
		s.recalculateOutcome()
	}
}

// HasBlockers returns true if any critical failures exist
func (s *ValidationSession) HasBlockers() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, result := range s.results {
		if result.IsBlocking() {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any warnings exist
func (s *ValidationSession) HasWarnings() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, result := range s.results {
		if result.IsWarning() {
			return true
		}
	}
	return false
}

// BlockingResults returns all blocking validation failures
func (s *ValidationSession) BlockingResults() []ValidationResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var blockers []ValidationResult
	for _, result := range s.results {
		if result.IsBlocking() {
			blockers = append(blockers, result)
		}
	}
	if blockers == nil {
		return nil
	}
	result := make([]ValidationResult, len(blockers))
	copy(result, blockers)
	return result
}

// WarningResults returns all warning-level results
func (s *ValidationSession) WarningResults() []ValidationResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var warnings []ValidationResult
	for _, result := range s.results {
		if result.IsWarning() {
			warnings = append(warnings, result)
		}
	}
	if warnings == nil {
		return nil
	}
	result := make([]ValidationResult, len(warnings))
	copy(result, warnings)
	return result
}

// CanProceed returns true if installation can continue
func (s *ValidationSession) CanProceed() bool {
	return !s.HasBlockers()
}

// recalculateOutcome determines the overall validation outcome
// NOTE: Caller must hold the lock (either read or write lock)
func (s *ValidationSession) recalculateOutcome() {
	// Check for blockers directly (don't call HasBlockers to avoid lock)
	hasBlockers := false
	hasWarnings := false

	for _, result := range s.results {
		if result.IsBlocking() {
			hasBlockers = true
			break
		}
	}

	if hasBlockers {
		s.overallResult = OutcomeBlocked
		return
	}

	// Check for warnings directly (don't call HasWarnings to avoid lock)
	for _, result := range s.results {
		if result.IsWarning() {
			hasWarnings = true
			break
		}
	}

	if hasWarnings {
		s.overallResult = OutcomeWarnings
		return
	}

	// Check if all validations passed
	allPassed := true
	for _, result := range s.results {
		if result.Status() != StatusPass {
			allPassed = false
			break
		}
	}

	if allPassed {
		s.overallResult = OutcomeSuccess
	} else {
		s.overallResult = OutcomePartialSuccess
	}
}

// Duration returns how long the validation took
func (s *ValidationSession) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.completedAt.IsZero() {
		return time.Since(s.startedAt)
	}
	return s.completedAt.Sub(s.startedAt)
}
