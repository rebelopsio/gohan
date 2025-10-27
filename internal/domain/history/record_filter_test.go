package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRecordFilter(t *testing.T) {
	t.Run("empty filter matches all records", func(t *testing.T) {
		filter := history.NewRecordFilter()

		assert.False(t, filter.HasPeriodFilter())
		assert.False(t, filter.HasOutcomeFilter())
		assert.False(t, filter.HasPackageFilter())
		assert.True(t, filter.IsEmpty())
	})
}

func TestRecordFilter_WithPeriod(t *testing.T) {
	period, err := history.NewPeriodFromDaysAgo(30)
	require.NoError(t, err)

	filter := history.NewRecordFilter().WithPeriod(period)

	assert.True(t, filter.HasPeriodFilter())
	assert.False(t, filter.IsEmpty())
	assert.Equal(t, period, filter.Period())
}

func TestRecordFilter_WithOutcome(t *testing.T) {
	outcome, err := history.NewInstallationOutcome("success")
	require.NoError(t, err)

	filter := history.NewRecordFilter().WithOutcome(outcome)

	assert.True(t, filter.HasOutcomeFilter())
	assert.False(t, filter.IsEmpty())
	assert.Equal(t, outcome, filter.Outcome())
}

func TestRecordFilter_WithPackageName(t *testing.T) {
	filter := history.NewRecordFilter().WithPackageName("hyprland")

	assert.True(t, filter.HasPackageFilter())
	assert.False(t, filter.IsEmpty())
	assert.Equal(t, "hyprland", filter.PackageName())
}

func TestRecordFilter_CombinedFilters(t *testing.T) {
	period, err := history.NewPeriodFromDaysAgo(30)
	require.NoError(t, err)

	outcome, err := history.NewInstallationOutcome("success")
	require.NoError(t, err)

	filter := history.NewRecordFilter().
		WithPeriod(period).
		WithOutcome(outcome).
		WithPackageName("hyprland")

	assert.True(t, filter.HasPeriodFilter())
	assert.True(t, filter.HasOutcomeFilter())
	assert.True(t, filter.HasPackageFilter())
	assert.False(t, filter.IsEmpty())
}

func TestRecordFilter_MatchesMetadata(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
	packages := []history.InstalledPackage{pkg1, pkg2}

	metadata, err := history.NewInstallationMetadata(
		"hyprland",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filter   history.RecordFilter
		metadata history.InstallationMetadata
		expected bool
	}{
		{
			name:     "empty filter matches any metadata",
			filter:   history.NewRecordFilter(),
			metadata: metadata,
			expected: true,
		},
		{
			name: "period filter matches when timestamp in range",
			filter: history.NewRecordFilter().WithPeriod(
				mustCreatePeriod(t, installedAt.AddDate(0, 0, -1), installedAt.AddDate(0, 0, 1)),
			),
			metadata: metadata,
			expected: true,
		},
		{
			name: "period filter does not match when timestamp out of range",
			filter: history.NewRecordFilter().WithPeriod(
				mustCreatePeriod(t, installedAt.AddDate(0, 0, -10), installedAt.AddDate(0, 0, -5)),
			),
			metadata: metadata,
			expected: false,
		},
		{
			name: "package filter matches when package exists",
			filter: history.NewRecordFilter().WithPackageName("hyprland"),
			metadata: metadata,
			expected: true,
		},
		{
			name: "package filter does not match when package doesn't exist",
			filter: history.NewRecordFilter().WithPackageName("kitty"),
			metadata: metadata,
			expected: false,
		},
		{
			name: "combined filters all match",
			filter: history.NewRecordFilter().
				WithPeriod(mustCreatePeriod(t, installedAt.AddDate(0, 0, -1), installedAt.AddDate(0, 0, 1))).
				WithPackageName("hyprland"),
			metadata: metadata,
			expected: true,
		},
		{
			name: "combined filters - one doesn't match",
			filter: history.NewRecordFilter().
				WithPeriod(mustCreatePeriod(t, installedAt.AddDate(0, 0, -1), installedAt.AddDate(0, 0, 1))).
				WithPackageName("kitty"),
			metadata: metadata,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.filter.MatchesMetadata(tt.metadata))
		})
	}
}

func TestRecordFilter_MatchesOutcome(t *testing.T) {
	successOutcome, _ := history.NewInstallationOutcome("success")
	failedOutcome, _ := history.NewInstallationOutcome("failed")

	tests := []struct {
		name     string
		filter   history.RecordFilter
		outcome  history.InstallationOutcome
		expected bool
	}{
		{
			name:     "empty filter matches any outcome",
			filter:   history.NewRecordFilter(),
			outcome:  successOutcome,
			expected: true,
		},
		{
			name:     "outcome filter matches exact outcome",
			filter:   history.NewRecordFilter().WithOutcome(successOutcome),
			outcome:  successOutcome,
			expected: true,
		},
		{
			name:     "outcome filter does not match different outcome",
			filter:   history.NewRecordFilter().WithOutcome(successOutcome),
			outcome:  failedOutcome,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.filter.MatchesOutcome(tt.outcome))
		})
	}
}

func TestRecordFilter_TrimsPackageName(t *testing.T) {
	filter := history.NewRecordFilter().WithPackageName("  hyprland  ")

	assert.Equal(t, "hyprland", filter.PackageName())
}

func TestRecordFilter_EmptyPackageNameIgnored(t *testing.T) {
	filter := history.NewRecordFilter().WithPackageName("")

	assert.False(t, filter.HasPackageFilter())
	assert.True(t, filter.IsEmpty())
}

func TestRecordFilter_WhitespacePackageNameIgnored(t *testing.T) {
	filter := history.NewRecordFilter().WithPackageName("   ")

	assert.False(t, filter.HasPackageFilter())
	assert.True(t, filter.IsEmpty())
}

func TestRecordFilter_Chaining(t *testing.T) {
	period, _ := history.NewPeriodFromDaysAgo(30)
	outcome, _ := history.NewInstallationOutcome("success")

	// Test method chaining works correctly
	filter := history.NewRecordFilter()
	filter = filter.WithPeriod(period)
	filter = filter.WithOutcome(outcome)
	filter = filter.WithPackageName("hyprland")

	assert.True(t, filter.HasPeriodFilter())
	assert.True(t, filter.HasOutcomeFilter())
	assert.True(t, filter.HasPackageFilter())
}

// Helper function to create period or fail test
func mustCreatePeriod(t *testing.T, start, end time.Time) history.InstallationPeriod {
	period, err := history.NewInstallationPeriod(start, end)
	require.NoError(t, err)
	return period
}
