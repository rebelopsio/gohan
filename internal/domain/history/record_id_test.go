package history_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRecordID(t *testing.T) {
	t.Run("generates valid UUID-based ID", func(t *testing.T) {
		id, err := history.NewRecordID()

		require.NoError(t, err)
		assert.NotEmpty(t, id.String())
		assert.True(t, id.IsValid())
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		id1, err := history.NewRecordID()
		require.NoError(t, err)

		id2, err := history.NewRecordID()
		require.NoError(t, err)

		assert.NotEqual(t, id1.String(), id2.String())
		assert.False(t, id1.Equals(id2))
	})
}

func TestParseRecordID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid UUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: nil,
		},
		{
			name:    "valid UUID with whitespace",
			input:   "  550e8400-e29b-41d4-a716-446655440000  ",
			wantErr: nil,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: history.ErrInvalidRecordID,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: history.ErrInvalidRecordID,
		},
		{
			name:    "invalid UUID format",
			input:   "not-a-uuid",
			wantErr: history.ErrInvalidRecordID,
		},
		{
			name:    "invalid UUID format - incomplete",
			input:   "550e8400-e29b-41d4-a716",
			wantErr: history.ErrInvalidRecordID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := history.ParseRecordID(tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, id.IsValid())
			} else {
				assert.NoError(t, err)
				assert.True(t, id.IsValid())
			}
		})
	}
}

func TestRecordID_String(t *testing.T) {
	id, err := history.ParseRecordID("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
}

func TestRecordID_ShortString(t *testing.T) {
	t.Run("returns first 8 characters for full UUID", func(t *testing.T) {
		id, err := history.ParseRecordID("550e8400-e29b-41d4-a716-446655440000")
		require.NoError(t, err)

		assert.Equal(t, "550e8400", id.ShortString())
	})
}

func TestRecordID_Equals(t *testing.T) {
	id1, err := history.ParseRecordID("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	id2, err := history.ParseRecordID("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	id3, err := history.ParseRecordID("660e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}

func TestRecordID_IsValid(t *testing.T) {
	t.Run("valid ID", func(t *testing.T) {
		id, err := history.NewRecordID()
		require.NoError(t, err)

		assert.True(t, id.IsValid())
	})

	t.Run("zero value ID is invalid", func(t *testing.T) {
		var id history.RecordID

		assert.False(t, id.IsValid())
	})
}
