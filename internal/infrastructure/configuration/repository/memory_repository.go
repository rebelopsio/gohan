package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
)

// MemoryRepository is an in-memory implementation of the configuration repository
// Suitable for development, testing, and single-instance deployments
// Thread-safe using sync.RWMutex
type MemoryRepository struct {
	mu        sync.RWMutex
	templates map[string]*configuration.ConfigurationTemplate
	nameIndex map[string]string // name -> id mapping for fast name lookups
}

// NewMemoryRepository creates a new in-memory configuration repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		templates: make(map[string]*configuration.ConfigurationTemplate),
		nameIndex: make(map[string]string),
	}
}

// Save persists a configuration template (create or update)
func (r *MemoryRepository) Save(ctx context.Context, template *configuration.ConfigurationTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update name index
	r.nameIndex[template.Metadata().Name().String()] = template.ID()

	// Store template
	r.templates[template.ID()] = template

	return nil
}

// FindByID retrieves a template by its unique identifier
func (r *MemoryRepository) FindByID(ctx context.Context, id string) (*configuration.ConfigurationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[id]
	if !exists {
		return nil, configuration.ErrTemplateNotFound
	}

	return template, nil
}

// FindByName retrieves a template by its name
func (r *MemoryRepository) FindByName(ctx context.Context, name string) (*configuration.ConfigurationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.nameIndex[name]
	if !exists {
		return nil, configuration.ErrTemplateNotFound
	}

	template, exists := r.templates[id]
	if !exists {
		return nil, configuration.ErrTemplateNotFound
	}

	return template, nil
}

// ExistsByName checks if a template with the given name exists
func (r *MemoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.nameIndex[name]
	return exists, nil
}

// List retrieves all configuration templates ordered by creation date (most recent first)
func (r *MemoryRepository) List(ctx context.Context) ([]*configuration.ConfigurationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]*configuration.ConfigurationTemplate, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}

	// Sort by creation date, most recent first
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].CreatedAt().After(templates[j].CreatedAt())
	})

	return templates, nil
}

// ListByCategory retrieves templates filtered by category
func (r *MemoryRepository) ListByCategory(ctx context.Context, category configuration.ConfigurationCategory) ([]*configuration.ConfigurationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*configuration.ConfigurationTemplate
	for _, template := range r.templates {
		if template.Metadata().Category() == category {
			filtered = append(filtered, template)
		}
	}

	// Sort by creation date, most recent first
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt().After(filtered[j].CreatedAt())
	})

	return filtered, nil
}

// ListByTag retrieves templates that have the specified tag
func (r *MemoryRepository) ListByTag(ctx context.Context, tag string) ([]*configuration.ConfigurationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*configuration.ConfigurationTemplate
	for _, template := range r.templates {
		// Check if template has this tag
		for _, t := range template.Metadata().Tags() {
			if t == tag {
				filtered = append(filtered, template)
				break
			}
		}
	}

	// Sort by creation date, most recent first
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt().After(filtered[j].CreatedAt())
	})

	return filtered, nil
}

// Delete removes a configuration template by ID
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	template, exists := r.templates[id]
	if !exists {
		return fmt.Errorf("template %s: %w", id, configuration.ErrTemplateNotFound)
	}

	// Remove from name index
	delete(r.nameIndex, template.Metadata().Name().String())

	// Remove from templates map
	delete(r.templates, id)

	return nil
}

// Clear removes all templates (useful for testing)
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.templates = make(map[string]*configuration.ConfigurationTemplate)
	r.nameIndex = make(map[string]string)
}

// Count returns the number of templates (useful for testing)
func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.templates)
}
