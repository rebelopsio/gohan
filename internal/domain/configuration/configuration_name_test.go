package configuration_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   error
		wantValue string
	}{
		{
			name:      "valid name",
			input:     "Development Stack",
			wantErr:   nil,
			wantValue: "Development Stack",
		},
		{
			name:      "name with special characters",
			input:     "My-Config_v2.1",
			wantErr:   nil,
			wantValue: "My-Config_v2.1",
		},
		{
			name:      "name with spaces trimmed",
			input:     "  Trimmed Name  ",
			wantErr:   nil,
			wantValue: "Trimmed Name",
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: configuration.ErrInvalidConfigurationName,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: configuration.ErrInvalidConfigurationName,
		},
		{
			name:    "name too long",
			input:   string(make([]byte, 101)),
			wantErr: configuration.ErrConfigurationNameTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, err := configuration.NewConfigurationName(tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantValue, name.String())
		})
	}
}

func TestConfigurationName_Equals(t *testing.T) {
	name1, _ := configuration.NewConfigurationName("Test Config")
	name2, _ := configuration.NewConfigurationName("Test Config")
	name3, _ := configuration.NewConfigurationName("Different Config")

	assert.True(t, name1.Equals(name2))
	assert.False(t, name1.Equals(name3))
}
