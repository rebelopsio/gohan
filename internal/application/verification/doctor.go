package verification

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/verification"
)

// DoctorRequest contains parameters for system verification
type DoctorRequest struct {
	ShowProgress bool
	QuickCheck   bool // If true, only run critical checks
}

// DoctorResponse contains verification results
type DoctorResponse struct {
	ReportID        string
	OverallStatus   string
	PassedChecks    int
	WarningChecks   int
	FailedChecks    int
	CriticalIssues  int
	TotalChecks     int
	Results         []CheckResultDTO
	Recommendations []string
	DurationMs      int64
}

// CheckResultDTO represents a check result for display
type CheckResultDTO struct {
	Component   string
	Status      string
	Severity    string
	Message     string
	Details     []string
	Suggestions []string
}

// ProgressCallback is called for each check
type ProgressCallback func(checkerName string, result CheckResultDTO)

// Checkers aggregates all verification checkers
type Checkers struct {
	HyprlandChecker     verification.VerificationChecker
	ThemeChecker        verification.VerificationChecker
	ConfigChecker       verification.VerificationChecker
	// Additional checkers can be added here
}

// DoctorUseCase coordinates system verification
type DoctorUseCase struct {
	checkers Checkers
}

// NewDoctorUseCase creates a new use case instance
func NewDoctorUseCase(checkers Checkers) *DoctorUseCase {
	return &DoctorUseCase{
		checkers: checkers,
	}
}

// Execute runs system verification
func (uc *DoctorUseCase) Execute(ctx context.Context, req DoctorRequest) (*DoctorResponse, error) {
	// Create list of checkers to run
	checkerList := uc.getCheckers(req.QuickCheck)

	// Create orchestrator
	orchestrator := verification.NewVerificationOrchestrator(checkerList)

	// Execute verification
	var report *verification.VerificationReport
	if req.ShowProgress {
		report = orchestrator.RunVerificationWithProgress(ctx, func(name string, result verification.CheckResult) {
			// Progress handled by CLI layer
		})
	} else {
		report = orchestrator.RunVerification(ctx)
	}

	// Convert to response
	return uc.buildResponse(report), nil
}

// ExecuteWithProgress runs verification with progress callbacks
func (uc *DoctorUseCase) ExecuteWithProgress(
	ctx context.Context,
	req DoctorRequest,
	progressFn ProgressCallback,
) (*DoctorResponse, error) {
	// Create list of checkers to run
	checkerList := uc.getCheckers(req.QuickCheck)

	// Create orchestrator
	orchestrator := verification.NewVerificationOrchestrator(checkerList)

	// Execute with progress
	report := orchestrator.RunVerificationWithProgress(ctx, func(name string, result verification.CheckResult) {
		if progressFn != nil {
			progressFn(name, uc.convertResult(result))
		}
	})

	return uc.buildResponse(report), nil
}

func (uc *DoctorUseCase) getCheckers(quickCheck bool) []verification.VerificationChecker {
	checkers := []verification.VerificationChecker{}

	// Always include critical checkers
	if uc.checkers.HyprlandChecker != nil {
		checkers = append(checkers, uc.checkers.HyprlandChecker)
	}
	if uc.checkers.ConfigChecker != nil {
		checkers = append(checkers, uc.checkers.ConfigChecker)
	}

	// Include additional checkers for full check
	if !quickCheck {
		if uc.checkers.ThemeChecker != nil {
			checkers = append(checkers, uc.checkers.ThemeChecker)
		}
	}

	return checkers
}

func (uc *DoctorUseCase) buildResponse(report *verification.VerificationReport) *DoctorResponse {
	results := report.Results()

	response := &DoctorResponse{
		ReportID:       report.ID(),
		OverallStatus:  string(report.OverallStatus()),
		PassedChecks:   report.PassedCount(),
		WarningChecks:  report.WarningCount(),
		FailedChecks:   report.FailedCount(),
		CriticalIssues: report.CriticalCount(),
		TotalChecks:    len(results),
		Results:        make([]CheckResultDTO, 0, len(results)),
		Recommendations: []string{},
		DurationMs:     report.Duration().Milliseconds(),
	}

	for _, result := range results {
		response.Results = append(response.Results, uc.convertResult(result))

		// Collect suggestions as recommendations
		if !result.IsPassing() {
			response.Recommendations = append(response.Recommendations, result.Suggestions()...)
		}
	}

	return response
}

func (uc *DoctorUseCase) convertResult(result verification.CheckResult) CheckResultDTO {
	return CheckResultDTO{
		Component:   result.Component().String(),
		Status:      result.Status().String(),
		Severity:    result.Severity().String(),
		Message:     result.Message(),
		Details:     result.Details(),
		Suggestions: result.Suggestions(),
	}
}
