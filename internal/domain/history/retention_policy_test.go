package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRetentionPolicy(t *testing.T) {
	tests := []struct {
		name    string
		days    int
		wantErr error
	}{
		{
			name:    "valid 30 days",
			days:    30,
			wantErr: nil,
		},
		{
			name:    "valid 90 days (default)",
			days:    90,
			wantErr: nil,
		},
		{
			name:    "valid 1 day (minimum)",
			days:    1,
			wantErr: nil,
		},
		{
			name:    "valid max days",
			days:    history.MaxRetentionDays,
			wantErr: nil,
		},
		{
			name:    "zero days invalid",
			days:    0,
			wantErr: history.ErrInvalidRetentionPeriod,
		},
		{
			name:    "negative days invalid",
			days:    -1,
			wantErr: history.ErrInvalidRetentionPeriod,
		},
		{
			name:    "exceeds maximum",
			days:    history.MaxRetentionDays + 1,
			wantErr: history.ErrRetentionPeriodTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := history.NewRetentionPolicy(tt.days)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.days, policy.RetentionDays())
			}
		})
	}
}

func TestDefaultRetentionPolicy(t *testing.T) {
	policy := history.DefaultRetentionPolicy()

	assert.Equal(t, history.DefaultRetentionDays, policy.RetentionDays())
	assert.Equal(t, 90, policy.RetentionDays())
}

func TestRetentionPolicy_ShouldPurge(t *testing.T) {
	policy, err := history.NewRetentionPolicy(90)
	require.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name       string
		recordTime time.Time
		expected   bool
	}{
		{
			name:       "record from today - should not purge",
			recordTime: now,
			expected:   false,
		},
		{
			name:       "record from 30 days ago - should not purge",
			recordTime: now.AddDate(0, 0, -30),
			expected:   false,
		},
		{
			name:       "record from 89 days ago - should not purge",
			recordTime: now.AddDate(0, 0, -89),
			expected:   false,
		},
		{
			name:       "record from 90 days ago - borderline",
			recordTime: now.AddDate(0, 0, -90),
			expected:   true, // Due to timing, slightly over 90 days when checked
		},
		{
			name:       "record from 91 days ago - should purge",
			recordTime: now.AddDate(0, 0, -91),
			expected:   true,
		},
		{
			name:       "record from 100 days ago - should purge",
			recordTime: now.AddDate(0, 0, -100),
			expected:   true,
		},
		{
			name:       "record from 1 year ago - should purge",
			recordTime: now.AddDate(-1, 0, 0),
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, policy.ShouldPurge(tt.recordTime))
		})
	}
}

func TestRetentionPolicy_CutoffDate(t *testing.T) {
	policy, err := history.NewRetentionPolicy(90)
	require.NoError(t, err)

	cutoff := policy.CutoffDate()

	// Cutoff should be approximately 90 days ago
	expectedCutoff := time.Now().AddDate(0, 0, -90)
	assert.WithinDuration(t, expectedCutoff, cutoff, time.Second)

	// Records before cutoff should be purged
	recordBeforeCutoff := cutoff.Add(-24 * time.Hour)
	assert.True(t, policy.ShouldPurge(recordBeforeCutoff))

	// Records after cutoff should not be purged
	recordAfterCutoff := cutoff.Add(24 * time.Hour)
	assert.False(t, policy.ShouldPurge(recordAfterCutoff))
}

func TestRetentionPolicy_DifferentPeriods(t *testing.T) {
	tests := []struct {
		name           string
		retentionDays  int
		recordDaysAgo  int
		shouldBePurged bool
	}{
		{
			name:           "30-day policy, 20-day-old record",
			retentionDays:  30,
			recordDaysAgo:  20,
			shouldBePurged: false,
		},
		{
			name:           "30-day policy, 40-day-old record",
			retentionDays:  30,
			recordDaysAgo:  40,
			shouldBePurged: true,
		},
		{
			name:           "365-day policy, 200-day-old record",
			retentionDays:  365,
			recordDaysAgo:  200,
			shouldBePurged: false,
		},
		{
			name:           "365-day policy, 400-day-old record",
			retentionDays:  365,
			recordDaysAgo:  400,
			shouldBePurged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := history.NewRetentionPolicy(tt.retentionDays)
			require.NoError(t, err)

			recordTime := time.Now().AddDate(0, 0, -tt.recordDaysAgo)
			assert.Equal(t, tt.shouldBePurged, policy.ShouldPurge(recordTime))
		})
	}
}
