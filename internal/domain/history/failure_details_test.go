package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFailureDetails(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		reason    string
		failedAt  time.Time
		phase     string
		errorCode string
		wantErr   error
	}{
		{
			name:      "valid failure details with all fields",
			reason:    "Package dependency conflict",
			failedAt:  now,
			phase:     "package_installation",
			errorCode: "DEP_CONFLICT",
			wantErr:   nil,
		},
		{
			name:      "valid failure details without error code",
			reason:    "Network timeout",
			failedAt:  now,
			phase:     "download",
			errorCode: "",
			wantErr:   nil,
		},
		{
			name:      "with whitespace trimmed",
			reason:    "  Disk full  ",
			failedAt:  now,
			phase:     "  installation  ",
			errorCode: "  DISK_FULL  ",
			wantErr:   nil,
		},
		{
			name:      "empty reason",
			reason:    "",
			failedAt:  now,
			phase:     "installation",
			errorCode: "ERR",
			wantErr:   history.ErrInvalidFailureReason,
		},
		{
			name:      "whitespace-only reason",
			reason:    "   ",
			failedAt:  now,
			phase:     "installation",
			errorCode: "ERR",
			wantErr:   history.ErrInvalidFailureReason,
		},
		{
			name:      "zero timestamp",
			reason:    "Some error",
			failedAt:  time.Time{},
			phase:     "installation",
			errorCode: "ERR",
			wantErr:   history.ErrInvalidTimestamp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details, err := history.NewFailureDetails(
				tt.reason,
				tt.failedAt,
				tt.phase,
				tt.errorCode,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, details.Reason())
			}
		})
	}
}

func TestFailureDetails_Accessors(t *testing.T) {
	failedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)

	details, err := history.NewFailureDetails(
		"Package conflict",
		failedAt,
		"installation",
		"DEP_001",
	)
	require.NoError(t, err)

	assert.Equal(t, "Package conflict", details.Reason())
	assert.Equal(t, failedAt, details.FailedAt())
	assert.Equal(t, "installation", details.Phase())
	assert.Equal(t, "DEP_001", details.ErrorCode())
}

func TestFailureDetails_HasErrorCode(t *testing.T) {
	t.Run("with error code", func(t *testing.T) {
		details, err := history.NewFailureDetails(
			"Error occurred",
			time.Now(),
			"phase",
			"ERR_123",
		)
		require.NoError(t, err)

		assert.True(t, details.HasErrorCode())
	})

	t.Run("without error code", func(t *testing.T) {
		details, err := history.NewFailureDetails(
			"Error occurred",
			time.Now(),
			"phase",
			"",
		)
		require.NoError(t, err)

		assert.False(t, details.HasErrorCode())
	})
}

func TestFailureDetails_TrimsWhitespace(t *testing.T) {
	details, err := history.NewFailureDetails(
		"  Failure reason  ",
		time.Now(),
		"  installation  ",
		"  ERR_001  ",
	)
	require.NoError(t, err)

	assert.Equal(t, "Failure reason", details.Reason())
	assert.Equal(t, "installation", details.Phase())
	assert.Equal(t, "ERR_001", details.ErrorCode())
}

func TestFailureDetails_AllowsEmptyOptionalFields(t *testing.T) {
	details, err := history.NewFailureDetails(
		"Failure reason",
		time.Now(),
		"",
		"",
	)
	require.NoError(t, err)

	assert.Equal(t, "Failure reason", details.Reason())
	assert.Empty(t, details.Phase())
	assert.Empty(t, details.ErrorCode())
	assert.False(t, details.HasErrorCode())
}
