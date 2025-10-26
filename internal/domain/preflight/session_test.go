package preflight_test

import (
	"sync"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidationSession(t *testing.T) {
	session := preflight.NewValidationSession()

	assert.NotEmpty(t, session.ID())
	assert.False(t, session.StartedAt().IsZero())
	assert.True(t, session.CompletedAt().IsZero())
	assert.Empty(t, session.Results())
	assert.False(t, session.HasBlockers())
	assert.False(t, session.HasWarnings())
}

func TestValidationSession_AddResult(t *testing.T) {
	session := preflight.NewValidationSession()

	result1 := createPassResult(preflight.RequirementDebianVersion)
	result2 := createPassResult(preflight.RequirementDiskSpace)

	session.AddResult(result1)
	assert.Len(t, session.Results(), 1)

	session.AddResult(result2)
	assert.Len(t, session.Results(), 2)
}

func TestValidationSession_Complete(t *testing.T) {
	session := preflight.NewValidationSession()

	assert.True(t, session.CompletedAt().IsZero())

	session.Complete()

	assert.False(t, session.CompletedAt().IsZero())
	assert.True(t, session.CompletedAt().After(session.StartedAt()))
}

func TestValidationSession_HasBlockers(t *testing.T) {
	tests := []struct {
		name         string
		results      []preflight.ValidationResult
		wantBlockers bool
	}{
		{
			name: "no blockers",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createPassResult(preflight.RequirementDiskSpace),
			},
			wantBlockers: false,
		},
		{
			name: "has critical blocker",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityCritical),
			},
			wantBlockers: true,
		},
		{
			name: "has high severity blocker",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createBlockingResult(preflight.RequirementInternet, preflight.SeverityHigh),
			},
			wantBlockers: true,
		},
		{
			name: "has warning but no blockers",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createWarningResult(preflight.RequirementGPUSupport),
			},
			wantBlockers: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := preflight.NewValidationSession()
			for _, result := range tt.results {
				session.AddResult(result)
			}

			assert.Equal(t, tt.wantBlockers, session.HasBlockers())
		})
	}
}

func TestValidationSession_HasWarnings(t *testing.T) {
	tests := []struct {
		name         string
		results      []preflight.ValidationResult
		wantWarnings bool
	}{
		{
			name: "no warnings",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createPassResult(preflight.RequirementDiskSpace),
			},
			wantWarnings: false,
		},
		{
			name: "has warning",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createWarningResult(preflight.RequirementGPUSupport),
			},
			wantWarnings: true,
		},
		{
			name: "has medium severity failure (counts as warning)",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createMediumFailureResult(preflight.RequirementSourceRepos),
			},
			wantWarnings: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := preflight.NewValidationSession()
			for _, result := range tt.results {
				session.AddResult(result)
			}

			assert.Equal(t, tt.wantWarnings, session.HasWarnings())
		})
	}
}

func TestValidationSession_BlockingResults(t *testing.T) {
	session := preflight.NewValidationSession()

	// Add mix of results
	session.AddResult(createPassResult(preflight.RequirementDebianVersion))
	session.AddResult(createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityCritical))
	session.AddResult(createWarningResult(preflight.RequirementGPUSupport))
	session.AddResult(createBlockingResult(preflight.RequirementInternet, preflight.SeverityHigh))

	blockers := session.BlockingResults()

	assert.Len(t, blockers, 2)
	assert.Equal(t, preflight.RequirementDiskSpace, blockers[0].RequirementName())
	assert.Equal(t, preflight.RequirementInternet, blockers[1].RequirementName())
}

func TestValidationSession_WarningResults(t *testing.T) {
	session := preflight.NewValidationSession()

	// Add mix of results
	session.AddResult(createPassResult(preflight.RequirementDebianVersion))
	session.AddResult(createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityCritical))
	session.AddResult(createWarningResult(preflight.RequirementGPUSupport))
	session.AddResult(createMediumFailureResult(preflight.RequirementSourceRepos))

	warnings := session.WarningResults()

	assert.Len(t, warnings, 2)
	assert.Equal(t, preflight.RequirementGPUSupport, warnings[0].RequirementName())
	assert.Equal(t, preflight.RequirementSourceRepos, warnings[1].RequirementName())
}

func TestValidationSession_CanProceed(t *testing.T) {
	tests := []struct {
		name        string
		results     []preflight.ValidationResult
		canProceed  bool
	}{
		{
			name: "all pass - can proceed",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createPassResult(preflight.RequirementDiskSpace),
				createPassResult(preflight.RequirementInternet),
			},
			canProceed: true,
		},
		{
			name: "has warnings but no blockers - can proceed",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createWarningResult(preflight.RequirementGPUSupport),
			},
			canProceed: true,
		},
		{
			name: "has blocker - cannot proceed",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityCritical),
			},
			canProceed: false,
		},
		{
			name: "has multiple blockers - cannot proceed",
			results: []preflight.ValidationResult{
				createBlockingResult(preflight.RequirementDebianVersion, preflight.SeverityCritical),
				createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityHigh),
			},
			canProceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := preflight.NewValidationSession()
			for _, result := range tt.results {
				session.AddResult(result)
			}

			assert.Equal(t, tt.canProceed, session.CanProceed())
		})
	}
}

func TestValidationSession_OverallResult(t *testing.T) {
	tests := []struct {
		name           string
		results        []preflight.ValidationResult
		wantOutcome    preflight.ValidationOutcome
	}{
		{
			name: "all pass - success",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createPassResult(preflight.RequirementDiskSpace),
			},
			wantOutcome: preflight.OutcomeSuccess,
		},
		{
			name: "has blocker - blocked",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityCritical),
			},
			wantOutcome: preflight.OutcomeBlocked,
		},
		{
			name: "has warnings but no blockers - warnings outcome",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createWarningResult(preflight.RequirementGPUSupport),
			},
			wantOutcome: preflight.OutcomeWarnings,
		},
		{
			name: "has multiple warnings - warnings outcome",
			results: []preflight.ValidationResult{
				createPassResult(preflight.RequirementDebianVersion),
				createWarningResult(preflight.RequirementGPUSupport),
				createMediumFailureResult(preflight.RequirementSourceRepos),
			},
			wantOutcome: preflight.OutcomeWarnings,
		},
		{
			name: "has both blockers and warnings - blocked",
			results: []preflight.ValidationResult{
				createBlockingResult(preflight.RequirementDebianVersion, preflight.SeverityCritical),
				createWarningResult(preflight.RequirementGPUSupport),
			},
			wantOutcome: preflight.OutcomeBlocked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := preflight.NewValidationSession()
			for _, result := range tt.results {
				session.AddResult(result)
			}

			assert.Equal(t, tt.wantOutcome, session.OverallResult())
		})
	}
}

func TestValidationSession_Duration(t *testing.T) {
	session := preflight.NewValidationSession()

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Duration before completion
	durationBefore := session.Duration()
	assert.Greater(t, durationBefore, time.Duration(0))

	// Complete the session
	session.Complete()

	// Duration after completion should be fixed
	durationAfter := session.Duration()
	assert.Greater(t, durationAfter, time.Duration(0))

	// Sleep and check duration hasn't changed (it's completed)
	time.Sleep(10 * time.Millisecond)
	durationAfter2 := session.Duration()
	assert.Equal(t, durationAfter, durationAfter2)
}

func TestValidationSession_RealWorldScenario_AllPass(t *testing.T) {
	session := preflight.NewValidationSession()

	// All validations pass
	session.AddResult(createPassResult(preflight.RequirementDebianVersion))
	session.AddResult(createPassResult(preflight.RequirementGPUSupport))
	session.AddResult(createPassResult(preflight.RequirementDiskSpace))
	session.AddResult(createPassResult(preflight.RequirementInternet))
	session.AddResult(createPassResult(preflight.RequirementSourceRepos))

	session.Complete()

	assert.Equal(t, preflight.OutcomeSuccess, session.OverallResult())
	assert.True(t, session.CanProceed())
	assert.False(t, session.HasBlockers())
	assert.False(t, session.HasWarnings())
	assert.Len(t, session.Results(), 5)
}

func TestValidationSession_RealWorldScenario_WithWarnings(t *testing.T) {
	session := preflight.NewValidationSession()

	// Most pass, but NVIDIA warning
	session.AddResult(createPassResult(preflight.RequirementDebianVersion))
	session.AddResult(createWarningResult(preflight.RequirementGPUSupport)) // NVIDIA warning
	session.AddResult(createPassResult(preflight.RequirementDiskSpace))
	session.AddResult(createPassResult(preflight.RequirementInternet))
	session.AddResult(createPassResult(preflight.RequirementSourceRepos))

	session.Complete()

	assert.Equal(t, preflight.OutcomeWarnings, session.OverallResult())
	assert.True(t, session.CanProceed()) // Can still proceed with warnings
	assert.False(t, session.HasBlockers())
	assert.True(t, session.HasWarnings())
	assert.Len(t, session.WarningResults(), 1)
}

func TestValidationSession_RealWorldScenario_Blocked(t *testing.T) {
	session := preflight.NewValidationSession()

	// Bookworm detected (blocks) and insufficient disk space (blocks)
	session.AddResult(createBlockingResult(preflight.RequirementDebianVersion, preflight.SeverityCritical))
	session.AddResult(createPassResult(preflight.RequirementGPUSupport))
	session.AddResult(createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityHigh))
	session.AddResult(createWarningResult(preflight.RequirementSourceRepos))

	session.Complete()

	assert.Equal(t, preflight.OutcomeBlocked, session.OverallResult())
	assert.False(t, session.CanProceed())
	assert.True(t, session.HasBlockers())
	assert.True(t, session.HasWarnings())
	assert.Len(t, session.BlockingResults(), 2)
	assert.Len(t, session.WarningResults(), 1)
}

func TestValidationSession_ConcurrentAddResult(t *testing.T) {
	session := preflight.NewValidationSession()
	var wg sync.WaitGroup

	// Add 100 results concurrently
	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			session.AddResult(createPassResult(preflight.RequirementDebianVersion))
		}()
	}

	wg.Wait()

	results := session.Results()
	assert.Len(t, results, numGoroutines, "All concurrent AddResult calls should succeed")
}

func TestValidationSession_Results_DefensiveCopy(t *testing.T) {
	session := preflight.NewValidationSession()
	result := createPassResult(preflight.RequirementDebianVersion)
	session.AddResult(result)

	// Get results slice
	results1 := session.Results()
	require.Len(t, results1, 1)

	// Try to modify returned slice by appending
	results1 = append(results1, createPassResult(preflight.RequirementDiskSpace))

	// Get results again - should not be affected by external modification
	results2 := session.Results()
	assert.Len(t, results2, 1, "External modification should not affect internal session state")
}

func TestValidationSession_BlockingResults_DefensiveCopy(t *testing.T) {
	session := preflight.NewValidationSession()
	session.AddResult(createBlockingResult(preflight.RequirementDebianVersion, preflight.SeverityCritical))

	// Get blocking results
	blockers1 := session.BlockingResults()
	require.Len(t, blockers1, 1)

	// Try to modify returned slice
	blockers1 = append(blockers1, createBlockingResult(preflight.RequirementDiskSpace, preflight.SeverityHigh))

	// Get blocking results again
	blockers2 := session.BlockingResults()
	assert.Len(t, blockers2, 1, "External modification should not affect internal state")
}

func TestValidationSession_WarningResults_DefensiveCopy(t *testing.T) {
	session := preflight.NewValidationSession()
	session.AddResult(createWarningResult(preflight.RequirementGPUSupport))

	// Get warning results
	warnings1 := session.WarningResults()
	require.Len(t, warnings1, 1)

	// Try to modify returned slice
	warnings1 = append(warnings1, createWarningResult(preflight.RequirementSourceRepos))

	// Get warning results again
	warnings2 := session.WarningResults()
	assert.Len(t, warnings2, 1, "External modification should not affect internal state")
}

func TestValidationSession_Complete_Idempotent(t *testing.T) {
	session := preflight.NewValidationSession()

	// Complete once
	session.Complete()
	firstCompletedAt := session.CompletedAt()
	assert.False(t, firstCompletedAt.IsZero())

	// Wait and complete again
	time.Sleep(10 * time.Millisecond)
	session.Complete()
	secondCompletedAt := session.CompletedAt()

	// CompletedAt should not change
	assert.Equal(t, firstCompletedAt, secondCompletedAt,
		"Complete should be idempotent - completedAt should not change on second call")
}

// Helper functions for creating test results

func createPassResult(requirement preflight.RequirementName) preflight.ValidationResult {
	return preflight.NewValidationResult(
		requirement,
		preflight.StatusPass,
		preflight.SeverityLow,
		"pass",
		"pass",
		preflight.NewUserGuidance("", "", nil, ""),
	)
}

func createBlockingResult(requirement preflight.RequirementName, severity preflight.Severity) preflight.ValidationResult {
	return preflight.NewValidationResult(
		requirement,
		preflight.StatusFail,
		severity,
		"fail",
		"expected",
		preflight.NewUserGuidance("Blocking failure", "", nil, ""),
	)
}

func createWarningResult(requirement preflight.RequirementName) preflight.ValidationResult {
	return preflight.NewValidationResult(
		requirement,
		preflight.StatusWarning,
		preflight.SeverityMedium,
		"warning",
		"expected",
		preflight.NewUserGuidance("Warning message", "", nil, ""),
	)
}

func createMediumFailureResult(requirement preflight.RequirementName) preflight.ValidationResult {
	return preflight.NewValidationResult(
		requirement,
		preflight.StatusFail,
		preflight.SeverityMedium,
		"fail",
		"expected",
		preflight.NewUserGuidance("Medium failure", "", nil, ""),
	)
}
