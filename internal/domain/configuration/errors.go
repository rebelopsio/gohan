package configuration

import "errors"

// Domain errors for configuration management
var (
	// Configuration validation errors
	ErrInvalidConfigurationName = errors.New("configuration name is invalid")
	ErrConfigurationNameTooLong = errors.New("configuration name exceeds maximum length")
	ErrInvalidDescription       = errors.New("description is invalid")
	ErrDescriptionTooLong       = errors.New("description exceeds maximum length")
	ErrTooManyTags              = errors.New("too many tags (maximum 10)")
	ErrInvalidTag               = errors.New("tag is invalid")

	// Manifest errors
	ErrInvalidManifest          = errors.New("configuration manifest is invalid")
	ErrNoComponents             = errors.New("manifest must have at least one component")
	ErrMissingCoreComponent     = errors.New("manifest must include core Hyprland component")
	ErrInvalidVersionConstraint = errors.New("version constraint is invalid")

	// Template errors
	ErrTemplateNotFound     = errors.New("configuration template not found")
	ErrDuplicateTemplate    = errors.New("configuration template with this name already exists")
	ErrInvalidTemplateState = errors.New("configuration template is in invalid state")

	// Composition errors
	ErrConflictingComponents = errors.New("configurations have conflicting components")
	ErrIncompatibleConfigs   = errors.New("configurations are incompatible for composition")
)

// Constants for validation
const (
	MaxConfigurationNameLength = 100
	MaxDescriptionLength       = 1000
	MaxTagCount                = 10
	MaxTagLength               = 50
)
