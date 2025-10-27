package configuration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ConfigurationTemplate is the aggregate root for saved/reusable configurations
// It represents a template that can be used to install a specific set of components
type ConfigurationTemplate struct {
	id        string
	metadata  ConfigurationMetadata
	manifest  ConfigurationManifest
	createdAt time.Time
	version   int // For optimistic locking
}

// NewConfigurationTemplate creates a new configuration template aggregate root
// Generates a unique ID and sets creation timestamp
func NewConfigurationTemplate(
	metadata ConfigurationMetadata,
	manifest ConfigurationManifest,
) (*ConfigurationTemplate, error) {
	id := uuid.New().String()

	return &ConfigurationTemplate{
		id:        id,
		metadata:  metadata,
		manifest:  manifest,
		createdAt: time.Now(),
		version:   1,
	}, nil
}

// ReconstructConfigurationTemplate reconstructs a configuration template from persistent storage
// This is a domain factory method that allows repositories to rebuild the aggregate
// while maintaining encapsulation and enforcing invariants
func ReconstructConfigurationTemplate(
	id string,
	metadata ConfigurationMetadata,
	manifest ConfigurationManifest,
	createdAt time.Time,
	version int,
) (*ConfigurationTemplate, error) {
	// Validate reconstruction parameters
	if id == "" {
		return nil, fmt.Errorf("template ID cannot be empty")
	}

	if createdAt.IsZero() {
		return nil, fmt.Errorf("created time cannot be zero")
	}

	if version < 1 {
		return nil, fmt.Errorf("version must be positive")
	}

	return &ConfigurationTemplate{
		id:        id,
		metadata:  metadata,
		manifest:  manifest,
		createdAt: createdAt,
		version:   version,
	}, nil
}

// ID returns the unique identifier for this template
func (t *ConfigurationTemplate) ID() string {
	return t.id
}

// Metadata returns the configuration metadata
func (t *ConfigurationTemplate) Metadata() ConfigurationMetadata {
	return t.metadata
}

// Manifest returns the configuration manifest
func (t *ConfigurationTemplate) Manifest() ConfigurationManifest {
	return t.manifest
}

// CreatedAt returns when the template was created
func (t *ConfigurationTemplate) CreatedAt() time.Time {
	return t.createdAt
}

// Version returns the current version for optimistic locking
func (t *ConfigurationTemplate) Version() int {
	return t.version
}

// Age returns how long ago the template was created
func (t *ConfigurationTemplate) Age() time.Duration {
	return time.Since(t.createdAt)
}

// String returns a human-readable representation
func (t *ConfigurationTemplate) String() string {
	return fmt.Sprintf("Configuration '%s' (%s, %d components)",
		t.metadata.Name().String(),
		t.metadata.Category().String(),
		t.manifest.ComponentCount(),
	)
}
