package configuration

import "strings"

// ConfigurationCategory represents the type/purpose of a configuration
type ConfigurationCategory string

const (
	CategoryDevelopment ConfigurationCategory = "development"
	CategoryProduction  ConfigurationCategory = "production"
	CategoryTesting     ConfigurationCategory = "testing"
	CategoryCustom      ConfigurationCategory = "custom"
)

// String returns the string representation of the category
func (c ConfigurationCategory) String() string {
	switch c {
	case CategoryDevelopment:
		return "Development"
	case CategoryProduction:
		return "Production"
	case CategoryTesting:
		return "Testing"
	case CategoryCustom:
		return "Custom"
	default:
		return "Custom"
	}
}

// ConfigurationMetadata is a value object containing descriptive information about a configuration
type ConfigurationMetadata struct {
	name        ConfigurationName
	description string
	author      string
	tags        []string
	category    ConfigurationCategory
}

// NewConfigurationMetadata creates a new configuration metadata value object
// Name is required, other fields are optional
// Tags are validated, deduplicated, and defensively copied
func NewConfigurationMetadata(
	name string,
	description string,
	author string,
	tags []string,
	category ConfigurationCategory,
) (ConfigurationMetadata, error) {
	// Validate and create configuration name
	configName, err := NewConfigurationName(name)
	if err != nil {
		return ConfigurationMetadata{}, err
	}

	// Trim description
	description = strings.TrimSpace(description)

	// Validate description length
	if len(description) > MaxDescriptionLength {
		return ConfigurationMetadata{}, ErrDescriptionTooLong
	}

	// Trim author
	author = strings.TrimSpace(author)

	// Process tags: validate, deduplicate, and copy
	processedTags, err := processTags(tags)
	if err != nil {
		return ConfigurationMetadata{}, err
	}

	return ConfigurationMetadata{
		name:        configName,
		description: description,
		author:      author,
		tags:        processedTags,
		category:    category,
	}, nil
}

// processTags validates, deduplicates, and copies tags
func processTags(tags []string) ([]string, error) {
	if tags == nil {
		return nil, nil
	}

	// Validate count
	if len(tags) > MaxTagCount {
		return nil, ErrTooManyTags
	}

	// Use map for deduplication
	tagSet := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		// Validate tag
		if tag == "" {
			return nil, ErrInvalidTag
		}

		if len(tag) > MaxTagLength {
			return nil, ErrInvalidTag
		}

		// Add if not duplicate
		if !tagSet[tag] {
			tagSet[tag] = true
			result = append(result, tag)
		}
	}

	return result, nil
}

// Name returns the configuration name
func (m ConfigurationMetadata) Name() ConfigurationName {
	return m.name
}

// Description returns the description
func (m ConfigurationMetadata) Description() string {
	return m.description
}

// Author returns the author
func (m ConfigurationMetadata) Author() string {
	return m.author
}

// Tags returns a defensive copy of the tags
func (m ConfigurationMetadata) Tags() []string {
	if m.tags == nil {
		return nil
	}
	tags := make([]string, len(m.tags))
	copy(tags, m.tags)
	return tags
}

// HasTags returns true if the metadata has tags
func (m ConfigurationMetadata) HasTags() bool {
	return len(m.tags) > 0
}

// TagCount returns the number of tags
func (m ConfigurationMetadata) TagCount() int {
	return len(m.tags)
}

// Category returns the configuration category
func (m ConfigurationMetadata) Category() ConfigurationCategory {
	return m.category
}
