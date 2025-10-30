package postinstall

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// SetupSession is the aggregate root for post-installation setup
type SetupSession struct {
	mu          sync.RWMutex
	id          string
	startedAt   time.Time
	completedAt time.Time
	results     map[ComponentType]ComponentResult
	rollbackLog []RollbackAction
}

// RollbackAction represents an action that can be rolled back
type RollbackAction struct {
	Component   ComponentType
	Description string
	UndoFunc    func() error
	Timestamp   time.Time
}

// NewSetupSession creates a new setup session
func NewSetupSession() *SetupSession {
	return &SetupSession{
		id:          uuid.New().String(),
		startedAt:   time.Now(),
		results:     make(map[ComponentType]ComponentResult),
		rollbackLog: []RollbackAction{},
	}
}

// ID returns the session ID
func (s *SetupSession) ID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.id
}

// AddResult adds a component result to the session
func (s *SetupSession) AddResult(result ComponentResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.Component()] = result
}

// AddRollbackAction adds a rollback action
func (s *SetupSession) AddRollbackAction(action RollbackAction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rollbackLog = append(s.rollbackLog, action)
}

// Complete marks the session as completed
func (s *SetupSession) Complete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.completedAt = time.Now()
}

// IsComplete returns true if the session is completed
func (s *SetupSession) IsComplete() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.completedAt.IsZero()
}

// Duration returns how long the session took
func (s *SetupSession) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.completedAt.IsZero() {
		return time.Since(s.startedAt)
	}
	return s.completedAt.Sub(s.startedAt)
}

// Results returns all component results
func (s *SetupSession) Results() map[ComponentType]ComponentResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create defensive copy
	results := make(map[ComponentType]ComponentResult, len(s.results))
	for k, v := range s.results {
		results[k] = v
	}
	return results
}

// GetResult returns the result for a specific component
func (s *SetupSession) GetResult(component ComponentType) (ComponentResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.results[component]
	return result, ok
}

// HasFailures returns true if any component failed
func (s *SetupSession) HasFailures() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, result := range s.results {
		if result.IsFailure() {
			return true
		}
	}
	return false
}

// SuccessCount returns the number of successful components
func (s *SetupSession) SuccessCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, result := range s.results {
		if result.IsSuccess() {
			count++
		}
	}
	return count
}

// FailureCount returns the number of failed components
func (s *SetupSession) FailureCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, result := range s.results {
		if result.IsFailure() {
			count++
		}
	}
	return count
}

// TotalComponents returns the total number of components processed
func (s *SetupSession) TotalComponents() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.results)
}

// RollbackActions returns all rollback actions
func (s *SetupSession) RollbackActions() []RollbackAction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create defensive copy
	actions := make([]RollbackAction, len(s.rollbackLog))
	copy(actions, s.rollbackLog)
	return actions
}

// CanRollback returns true if there are actions to roll back
func (s *SetupSession) CanRollback() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.rollbackLog) > 0
}

// Rollback executes all rollback actions in reverse order
func (s *SetupSession) Rollback() []error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var errors []error

	// Execute rollback actions in reverse order
	for i := len(s.rollbackLog) - 1; i >= 0; i-- {
		action := s.rollbackLog[i]
		if action.UndoFunc != nil {
			if err := action.UndoFunc(); err != nil {
				errors = append(errors, err)
			}
		}
	}

	return errors
}

// OverallStatus returns the overall status of the session
func (s *SetupSession) OverallStatus() SetupStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.results) == 0 {
		return StatusPending
	}

	hasFailure := false
	allCompleted := true

	for _, result := range s.results {
		if result.IsFailure() {
			hasFailure = true
		}
		if result.Status() != StatusCompleted && result.Status() != StatusSkipped {
			allCompleted = false
		}
	}

	if hasFailure {
		return StatusFailed
	}
	if allCompleted {
		return StatusCompleted
	}
	return StatusInProgress
}
