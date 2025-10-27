package configuration

import "strings"

// ConfigurationName is a value object representing a validated configuration name
type ConfigurationName struct {
	value string
}

// NewConfigurationName creates a new configuration name value object
// Name is required, trimmed, and validated for length
func NewConfigurationName(name string) (ConfigurationName, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Validate not empty
	if name == "" {
		return ConfigurationName{}, ErrInvalidConfigurationName
	}

	// Validate length
	if len(name) > MaxConfigurationNameLength {
		return ConfigurationName{}, ErrConfigurationNameTooLong
	}

	return ConfigurationName{value: name}, nil
}

// String returns the configuration name as a string
func (n ConfigurationName) String() string {
	return n.value
}

// Equals checks if two configuration names are equal
func (n ConfigurationName) Equals(other ConfigurationName) bool {
	return n.value == other.value
}
