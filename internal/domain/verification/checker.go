package verification

import "context"

// VerificationChecker defines the interface for system verification checks
type VerificationChecker interface {
	// Name returns a human-readable name for this checker
	Name() string

	// Component returns the component being checked
	Component() ComponentName

	// Check performs the verification and returns a result
	Check(ctx context.Context) CheckResult
}

// VerificationOrchestrator coordinates multiple verification checkers
type VerificationOrchestrator struct {
	checkers []VerificationChecker
}

// NewVerificationOrchestrator creates a new orchestrator
func NewVerificationOrchestrator(checkers []VerificationChecker) *VerificationOrchestrator {
	return &VerificationOrchestrator{
		checkers: checkers,
	}
}

// RunVerification executes all checkers and returns a report
func (o *VerificationOrchestrator) RunVerification(ctx context.Context) *VerificationReport {
	report := NewVerificationReport()

	for _, checker := range o.checkers {
		result := checker.Check(ctx)
		report.AddResult(result)
	}

	report.Complete()
	return report
}

// RunVerificationWithProgress executes checkers with progress reporting
func (o *VerificationOrchestrator) RunVerificationWithProgress(
	ctx context.Context,
	progressFn func(checkerName string, result CheckResult),
) *VerificationReport {
	report := NewVerificationReport()

	for _, checker := range o.checkers {
		result := checker.Check(ctx)
		report.AddResult(result)

		if progressFn != nil {
			progressFn(checker.Name(), result)
		}
	}

	report.Complete()
	return report
}
