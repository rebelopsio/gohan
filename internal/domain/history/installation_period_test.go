package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationPeriod(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr error
	}{
		{
			name:    "valid period",
			start:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:     time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "same start and end time",
			start:   time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
			end:     time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		{
			name:    "end before start",
			start:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			end:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: history.ErrInvalidTimeRange,
		},
		{
			name:    "zero start time",
			start:   time.Time{},
			end:     time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			wantErr: history.ErrInvalidPeriod,
		},
		{
			name:    "zero end time",
			start:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:     time.Time{},
			wantErr: history.ErrInvalidPeriod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period, err := history.NewInstallationPeriod(tt.start, tt.end)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.start, period.Start())
				assert.Equal(t, tt.end, period.End())
			}
		})
	}
}

func TestNewPeriodFromDaysAgo(t *testing.T) {
	t.Run("valid days ago", func(t *testing.T) {
		period, err := history.NewPeriodFromDaysAgo(30)
		require.NoError(t, err)

		// Period should end approximately now
		assert.WithinDuration(t, time.Now(), period.End(), time.Second)

		// Period should start approximately 30 days ago
		expectedStart := time.Now().AddDate(0, 0, -30)
		assert.WithinDuration(t, expectedStart, period.Start(), time.Second)
	})

	t.Run("zero days", func(t *testing.T) {
		period, err := history.NewPeriodFromDaysAgo(0)
		require.NoError(t, err)

		// Both start and end should be approximately now
		assert.WithinDuration(t, time.Now(), period.Start(), time.Second)
		assert.WithinDuration(t, time.Now(), period.End(), time.Second)
	})

	t.Run("negative days", func(t *testing.T) {
		_, err := history.NewPeriodFromDaysAgo(-1)
		assert.ErrorIs(t, err, history.ErrInvalidPeriod)
	})
}

func TestInstallationPeriod_Contains(t *testing.T) {
	start := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 9, 30, 23, 59, 59, 0, time.UTC)
	period, err := history.NewInstallationPeriod(start, end)
	require.NoError(t, err)

	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{
			name:     "time at start boundary",
			time:     start,
			expected: true,
		},
		{
			name:     "time at end boundary",
			time:     end,
			expected: true,
		},
		{
			name:     "time in middle of period",
			time:     time.Date(2025, 9, 15, 12, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "time before start",
			time:     time.Date(2025, 8, 31, 23, 59, 59, 0, time.UTC),
			expected: false,
		},
		{
			name:     "time after end",
			time:     time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, period.Contains(tt.time))
		})
	}
}

func TestInstallationPeriod_Duration(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected time.Duration
	}{
		{
			name:     "one day",
			start:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: 24 * time.Hour,
		},
		{
			name:     "one hour",
			start:    time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC),
			expected: 1 * time.Hour,
		},
		{
			name:     "zero duration",
			start:    time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period, err := history.NewInstallationPeriod(tt.start, tt.end)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, period.Duration())
		})
	}
}
