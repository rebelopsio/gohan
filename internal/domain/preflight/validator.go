package preflight

import "context"

// Validator defines the interface for validation checks
type Validator interface {
	// Name returns the validator name for display
	Name() string

	// Validate executes the validation and returns a result
	Validate(ctx context.Context) ValidationResult

	// RequirementName returns the requirement being validated
	RequirementName() RequirementName
}

// ValidationOrchestrator coordinates all validation checks
type ValidationOrchestrator struct {
	validators []Validator
}

// NewValidationOrchestrator creates a new orchestrator
func NewValidationOrchestrator(validators []Validator) *ValidationOrchestrator {
	return &ValidationOrchestrator{
		validators: validators,
	}
}

// ExecuteValidations runs all validators and returns a session
func (o *ValidationOrchestrator) ExecuteValidations(ctx context.Context) *ValidationSession {
	session := NewValidationSession()

	for _, validator := range o.validators {
		result := validator.Validate(ctx)
		session.AddResult(result)
	}

	session.Complete()
	return session
}

// ExecuteValidationsWithProgress runs validators with progress updates
func (o *ValidationOrchestrator) ExecuteValidationsWithProgress(
	ctx context.Context,
	progressFn func(validator string, result ValidationResult),
) *ValidationSession {
	session := NewValidationSession()

	for _, validator := range o.validators {
		result := validator.Validate(ctx)
		session.AddResult(result)

		if progressFn != nil {
			progressFn(validator.Name(), result)
		}
	}

	session.Complete()
	return session
}
