package history_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationOutcome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  error
		expected history.InstallationOutcome
	}{
		{
			name:     "success outcome",
			input:    "success",
			wantErr:  nil,
			expected: history.OutcomeSuccess,
		},
		{
			name:     "failed outcome",
			input:    "failed",
			wantErr:  nil,
			expected: history.OutcomeFailed,
		},
		{
			name:     "rolled_back outcome",
			input:    "rolled_back",
			wantErr:  nil,
			expected: history.OutcomeRolledBack,
		},
		{
			name:    "invalid outcome",
			input:   "unknown",
			wantErr: history.ErrInvalidOutcome,
		},
		{
			name:    "empty outcome",
			input:   "",
			wantErr: history.ErrInvalidOutcome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outcome, err := history.NewInstallationOutcome(tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, outcome)
			}
		})
	}
}

func TestInstallationOutcome_String(t *testing.T) {
	tests := []struct {
		name     string
		outcome  history.InstallationOutcome
		expected string
	}{
		{
			name:     "success outcome",
			outcome:  history.OutcomeSuccess,
			expected: "success",
		},
		{
			name:     "failed outcome",
			outcome:  history.OutcomeFailed,
			expected: "failed",
		},
		{
			name:     "rolled_back outcome",
			outcome:  history.OutcomeRolledBack,
			expected: "rolled_back",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.outcome.String())
		})
	}
}

func TestInstallationOutcome_IsSuccessful(t *testing.T) {
	assert.True(t, history.OutcomeSuccess.IsSuccessful())
	assert.False(t, history.OutcomeFailed.IsSuccessful())
	assert.False(t, history.OutcomeRolledBack.IsSuccessful())
}

func TestInstallationOutcome_IsFailed(t *testing.T) {
	assert.False(t, history.OutcomeSuccess.IsFailed())
	assert.True(t, history.OutcomeFailed.IsFailed())
	assert.False(t, history.OutcomeRolledBack.IsFailed())
}

func TestInstallationOutcome_IsRolledBack(t *testing.T) {
	assert.False(t, history.OutcomeSuccess.IsRolledBack())
	assert.False(t, history.OutcomeFailed.IsRolledBack())
	assert.True(t, history.OutcomeRolledBack.IsRolledBack())
}

func TestInstallationOutcome_Equals(t *testing.T) {
	outcome1, err := history.NewInstallationOutcome("success")
	require.NoError(t, err)

	outcome2, err := history.NewInstallationOutcome("success")
	require.NoError(t, err)

	outcome3, err := history.NewInstallationOutcome("failed")
	require.NoError(t, err)

	assert.True(t, outcome1.Equals(outcome2))
	assert.False(t, outcome1.Equals(outcome3))
}
