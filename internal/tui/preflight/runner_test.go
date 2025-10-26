package preflight

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidationRunner(t *testing.T) {
	runner := NewValidationRunner()

	assert.NotNil(t, runner)
	assert.NotNil(t, runner.session)
	assert.NotNil(t, runner.debianDetector)
	assert.NotNil(t, runner.gpuDetector)
	assert.NotNil(t, runner.diskSpaceDetector)
	assert.NotNil(t, runner.connectivityChecker)
	assert.NotNil(t, runner.sourceRepoChecker)
	assert.NotNil(t, runner.progressChan)
}

func TestValidationRunner_Session(t *testing.T) {
	runner := NewValidationRunner()

	session := runner.Session()
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.ID())
}

func TestValidationRunner_Progress(t *testing.T) {
	runner := NewValidationRunner()

	progressChan := runner.Progress()
	assert.NotNil(t, progressChan)
}

func TestValidationRunner_Run_CompletesAllValidations(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run validation
	err := runner.Run(ctx)
	assert.NoError(t, err, "Run should not return error even if validations fail")

	// Verify session was completed
	session := runner.Session()
	assert.False(t, session.CompletedAt().IsZero(), "Session should be marked complete")
	assert.NotEmpty(t, session.Results(), "Session should have results")

	// Should have exactly 5 validation results (one for each check)
	results := session.Results()
	assert.Len(t, results, 5, "Should have 5 validation results")
}

func TestValidationRunner_Run_ProgressUpdates(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Collect progress updates
	updates := make([]ProgressUpdate, 0)
	done := make(chan struct{})

	go func() {
		for update := range runner.Progress() {
			updates = append(updates, update)
		}
		close(done)
	}()

	// Run validation
	err := runner.Run(ctx)
	require.NoError(t, err)

	// Wait for progress channel to close
	<-done

	// Verify we received progress updates
	assert.NotEmpty(t, updates, "Should receive progress updates")

	// Should have at least 5 updates (one for each validation)
	assert.GreaterOrEqual(t, len(updates), 5, "Should have at least 5 progress updates")

	// Verify all requirements were checked
	requirements := make(map[preflight.RequirementName]bool)
	for _, update := range updates {
		requirements[update.RequirementName] = true
	}

	assert.True(t, requirements[preflight.RequirementDebianVersion], "Should check Debian version")
	assert.True(t, requirements[preflight.RequirementGPUSupport], "Should check GPU")
	assert.True(t, requirements[preflight.RequirementDiskSpace], "Should check disk space")
	assert.True(t, requirements[preflight.RequirementInternet], "Should check connectivity")
	assert.True(t, requirements[preflight.RequirementSourceRepos], "Should check source repos")
}

func TestValidationRunner_Run_ContextCancellation(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context immediately
	cancel()

	// Run validation with canceled context
	err := runner.Run(ctx)

	// Should complete without panicking even with canceled context
	assert.NoError(t, err, "Should handle canceled context gracefully")
}

func TestValidationRunner_Run_HandlesDetectorErrors(t *testing.T) {
	runner := NewValidationRunner()
	ctx := context.Background()

	// Run validation - some detectors may error on this system
	err := runner.Run(ctx)
	require.NoError(t, err)

	// Verify session still has results even if some checks failed
	session := runner.Session()
	results := session.Results()

	assert.NotEmpty(t, results, "Should have results even if some checks failed")
	assert.Len(t, results, 5, "Should attempt all 5 validations")
}

func TestValidationRunner_ValidationResults_HaveGuidance(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := runner.Run(ctx)
	require.NoError(t, err)

	session := runner.Session()
	results := session.Results()

	// Check that failed validations have guidance
	for _, result := range results {
		if result.Status() == preflight.StatusFail || result.Status() == preflight.StatusWarning {
			guidance := result.Guidance()
			assert.NotEmpty(t, guidance.Message(), "Failed/warning results should have guidance message")

			// Most failures should have actionable steps
			if result.IsBlocking() {
				assert.True(t, guidance.HasSteps(), "Blocking failures should have actionable steps")
			}
		}
	}
}

func TestValidationRunner_SessionOutcome(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := runner.Run(ctx)
	require.NoError(t, err)

	session := runner.Session()
	outcome := session.OverallResult()

	// Outcome should be one of the valid outcomes
	assert.Contains(t, []preflight.ValidationOutcome{
		preflight.OutcomeSuccess,
		preflight.OutcomeWarnings,
		preflight.OutcomeBlocked,
		preflight.OutcomePartialSuccess,
	}, outcome, "Should have a valid outcome")
}

func TestValidationRunner_RequirementNames(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := runner.Run(ctx)
	require.NoError(t, err)

	session := runner.Session()
	results := session.Results()

	// Map results by requirement name
	resultsByName := make(map[preflight.RequirementName]preflight.ValidationResult)
	for _, result := range results {
		resultsByName[result.RequirementName()] = result
	}

	// Verify all expected requirements are present
	_, hasDebian := resultsByName[preflight.RequirementDebianVersion]
	_, hasGPU := resultsByName[preflight.RequirementGPUSupport]
	_, hasDisk := resultsByName[preflight.RequirementDiskSpace]
	_, hasInternet := resultsByName[preflight.RequirementInternet]
	_, hasSourceRepos := resultsByName[preflight.RequirementSourceRepos]

	assert.True(t, hasDebian, "Should have Debian version result")
	assert.True(t, hasGPU, "Should have GPU result")
	assert.True(t, hasDisk, "Should have disk space result")
	assert.True(t, hasInternet, "Should have internet result")
	assert.True(t, hasSourceRepos, "Should have source repos result")
}

func TestValidationRunner_Duration(t *testing.T) {
	runner := NewValidationRunner()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	err := runner.Run(ctx)
	elapsed := time.Since(start)

	require.NoError(t, err)

	session := runner.Session()
	duration := session.Duration()

	// Duration should be positive and reasonable
	assert.Greater(t, duration, time.Duration(0), "Duration should be positive")
	assert.LessOrEqual(t, duration, elapsed+100*time.Millisecond, "Duration should be close to actual elapsed time")
}

func TestProgressUpdate_Structure(t *testing.T) {
	update := ProgressUpdate{
		RequirementName: preflight.RequirementDebianVersion,
		Status:          preflight.StatusPass,
		Message:         "Test message",
		Result:          nil,
	}

	assert.Equal(t, preflight.RequirementDebianVersion, update.RequirementName)
	assert.Equal(t, preflight.StatusPass, update.Status)
	assert.Equal(t, "Test message", update.Message)
	assert.Nil(t, update.Result)
}
