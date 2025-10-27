package configuration

import "context"

// Repository defines the interface for configuration template persistence
// Implementations will be provided by the infrastructure layer
type Repository interface {
	// Save persists a configuration template (create or update)
	Save(ctx context.Context, template *ConfigurationTemplate) error

	// FindByID retrieves a template by its unique identifier
	// Returns ErrTemplateNotFound if the template doesn't exist
	FindByID(ctx context.Context, id string) (*ConfigurationTemplate, error)

	// FindByName retrieves a template by its name
	// Returns ErrTemplateNotFound if the template doesn't exist
	FindByName(ctx context.Context, name string) (*ConfigurationTemplate, error)

	// ExistsByName checks if a template with the given name exists
	// Used for uniqueness validation before creating new templates
	ExistsByName(ctx context.Context, name string) (bool, error)

	// List retrieves all configuration templates
	// Results are ordered by creation date (most recent first)
	List(ctx context.Context) ([]*ConfigurationTemplate, error)

	// ListByCategory retrieves templates filtered by category
	// Results are ordered by creation date (most recent first)
	ListByCategory(ctx context.Context, category ConfigurationCategory) ([]*ConfigurationTemplate, error)

	// ListByTag retrieves templates that have the specified tag
	// Results are ordered by creation date (most recent first)
	ListByTag(ctx context.Context, tag string) ([]*ConfigurationTemplate, error)

	// Delete removes a configuration template by ID
	// Returns ErrTemplateNotFound if the template doesn't exist
	Delete(ctx context.Context, id string) error
}
