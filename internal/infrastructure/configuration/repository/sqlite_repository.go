package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// SQLiteRepository is a SQLite implementation of the configuration repository
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite configuration repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	_, err = db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return repo, nil
}

// initialize creates the necessary tables
func (r *SQLiteRepository) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS configuration_templates (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		category TEXT NOT NULL,
		data TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		version INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_templates_name ON configuration_templates(name);
	CREATE INDEX IF NOT EXISTS idx_templates_category ON configuration_templates(category);
	CREATE INDEX IF NOT EXISTS idx_templates_created_at ON configuration_templates(created_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// templateRecord represents a row in the database
type templateRecord struct {
	ID        string
	Name      string
	Category  string
	Data      string
	CreatedAt time.Time
	Version   int
}

// templateStorageModel is a serializable representation of a configuration template
type templateStorageModel struct {
	ID        string       `json:"id"`
	Metadata  metadataDTO  `json:"metadata"`
	Manifest  manifestDTO  `json:"manifest"`
	CreatedAt time.Time    `json:"created_at"`
	Version   int          `json:"version"`
}

// metadataDTO is a serializable version of ConfigurationMetadata
type metadataDTO struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category"`
}

// manifestDTO is a serializable version of ConfigurationManifest
type manifestDTO struct {
	Components        []componentDTO `json:"components"`
	DiskRequiredBytes uint64         `json:"disk_required_bytes"`
	GPURequired       bool           `json:"gpu_required"`
}

// componentDTO is a serializable version of ComponentSelection
type componentDTO struct {
	Component    string   `json:"component"`
	Version      string   `json:"version"`
	PackageName  string   `json:"package_name,omitempty"`
	SizeBytes    uint64   `json:"size_bytes,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// toStorageModel converts a domain template to a storage model
func toStorageModel(template *configuration.ConfigurationTemplate) *templateStorageModel {
	// Convert metadata
	metadata := template.Metadata()
	metadataDTO := metadataDTO{
		Name:        metadata.Name().String(),
		Description: metadata.Description(),
		Author:      metadata.Author(),
		Tags:        metadata.Tags(),
		Category:    string(metadata.Category()),
	}

	// Convert manifest components
	manifest := template.Manifest()
	components := make([]componentDTO, 0, manifest.ComponentCount())
	for _, comp := range manifest.Components() {
		compDTO := componentDTO{
			Component: string(comp.Component()),
			Version:   comp.Version(),
		}
		if comp.HasPackageInfo() {
			compDTO.PackageName = comp.PackageInfo().Name()
			compDTO.SizeBytes = comp.PackageInfo().SizeBytes()
			compDTO.Dependencies = comp.PackageInfo().Dependencies()
		}
		components = append(components, compDTO)
	}

	manifestDTO := manifestDTO{
		Components:        components,
		DiskRequiredBytes: manifest.DiskRequiredBytes(),
		GPURequired:       manifest.GPURequired(),
	}

	return &templateStorageModel{
		ID:        template.ID(),
		Metadata:  metadataDTO,
		Manifest:  manifestDTO,
		CreatedAt: template.CreatedAt(),
		Version:   template.Version(),
	}
}

// fromStorageModel converts a storage DTO back to domain objects
func fromStorageModel(model *templateStorageModel) (*configuration.ConfigurationTemplate, error) {
	// Reconstruct component selections
	components := make([]installation.ComponentSelection, 0, len(model.Manifest.Components))
	for _, compDTO := range model.Manifest.Components {
		// Create package info if present
		var packageInfo *installation.PackageInfo
		if compDTO.PackageName != "" {
			pkgInfo, err := installation.NewPackageInfo(
				compDTO.PackageName,
				compDTO.Version,
				compDTO.SizeBytes,
				compDTO.Dependencies,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create package info: %w", err)
			}
			packageInfo = &pkgInfo
		}

		// Create component selection
		compSel, err := installation.NewComponentSelection(
			installation.ComponentName(compDTO.Component),
			compDTO.Version,
			packageInfo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create component selection: %w", err)
		}
		components = append(components, compSel)
	}

	// Reconstruct manifest
	manifest, err := configuration.NewConfigurationManifest(
		components,
		model.Manifest.DiskRequiredBytes,
		model.Manifest.GPURequired,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest: %w", err)
	}

	// Reconstruct metadata
	metadata, err := configuration.NewConfigurationMetadata(
		model.Metadata.Name,
		model.Metadata.Description,
		model.Metadata.Author,
		model.Metadata.Tags,
		configuration.ConfigurationCategory(model.Metadata.Category),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata: %w", err)
	}

	// Reconstruct template using domain factory
	template, err := configuration.ReconstructConfigurationTemplate(
		model.ID,
		metadata,
		manifest,
		model.CreatedAt,
		model.Version,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct template: %w", err)
	}

	return template, nil
}

// Save persists a configuration template (create or update)
func (r *SQLiteRepository) Save(ctx context.Context, template *configuration.ConfigurationTemplate) error {
	// Convert to storage model
	model := toStorageModel(template)

	// Serialize to JSON
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	// Upsert query
	query := `
	INSERT INTO configuration_templates (id, name, category, data, created_at, version)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		category = excluded.category,
		data = excluded.data,
		version = excluded.version
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		template.ID(),
		template.Metadata().Name().String(),
		string(template.Metadata().Category()),
		string(data),
		template.CreatedAt(),
		template.Version(),
	)

	return err
}

// FindByID retrieves a template by its unique identifier
func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (*configuration.ConfigurationTemplate, error) {
	query := `SELECT id, name, category, data, created_at, version FROM configuration_templates WHERE id = ?`

	var record templateRecord
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.Name,
		&record.Category,
		&record.Data,
		&record.CreatedAt,
		&record.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found: %s", id)
		}
		return nil, fmt.Errorf("failed to query template: %w", err)
	}

	// Deserialize storage model
	var model templateStorageModel
	if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
	}

	// Convert to domain using helper
	template, err := fromStorageModel(&model)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// FindByName retrieves a template by its name
func (r *SQLiteRepository) FindByName(ctx context.Context, name string) (*configuration.ConfigurationTemplate, error) {
	query := `SELECT id, name, category, data, created_at, version FROM configuration_templates WHERE name = ?`

	var record templateRecord
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&record.ID,
		&record.Name,
		&record.Category,
		&record.Data,
		&record.CreatedAt,
		&record.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found: %s", name)
		}
		return nil, fmt.Errorf("failed to query template: %w", err)
	}

	// Deserialize storage model
	var model templateStorageModel
	if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
	}

	// Convert to domain using helper
	template, err := fromStorageModel(&model)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// ExistsByName checks if a template with the given name exists
func (r *SQLiteRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT COUNT(*) FROM configuration_templates WHERE name = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check name existence: %w", err)
	}

	return count > 0, nil
}

// List retrieves all configuration templates ordered by creation date
func (r *SQLiteRepository) List(ctx context.Context) ([]*configuration.ConfigurationTemplate, error) {
	query := `SELECT id, name, category, data, created_at, version FROM configuration_templates ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*configuration.ConfigurationTemplate

	for rows.Next() {
		var record templateRecord
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Category,
			&record.Data,
			&record.CreatedAt,
			&record.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template row: %w", err)
		}

		// Deserialize storage model
		var model templateStorageModel
		if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
		}

		// Convert to domain using helper
		template, err := fromStorageModel(&model)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	return templates, nil
}

// ListByCategory retrieves templates filtered by category
func (r *SQLiteRepository) ListByCategory(ctx context.Context, category configuration.ConfigurationCategory) ([]*configuration.ConfigurationTemplate, error) {
	query := `SELECT id, name, category, data, created_at, version FROM configuration_templates WHERE category = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(category))
	if err != nil {
		return nil, fmt.Errorf("failed to query templates by category: %w", err)
	}
	defer rows.Close()

	var templates []*configuration.ConfigurationTemplate

	for rows.Next() {
		var record templateRecord
		err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Category,
			&record.Data,
			&record.CreatedAt,
			&record.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template row: %w", err)
		}

		// Deserialize storage model
		var model templateStorageModel
		if err := json.Unmarshal([]byte(record.Data), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template data: %w", err)
		}

		// Convert to domain using helper
		template, err := fromStorageModel(&model)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating templates: %w", err)
	}

	return templates, nil
}

// ListByTag retrieves templates that have the specified tag
func (r *SQLiteRepository) ListByTag(ctx context.Context, tag string) ([]*configuration.ConfigurationTemplate, error) {
	// Since tags are stored in JSON, we need to fetch all and filter in memory
	// For better performance with many templates, consider a separate tags table
	templates, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*configuration.ConfigurationTemplate
	for _, template := range templates {
		for _, t := range template.Metadata().Tags() {
			if t == tag {
				filtered = append(filtered, template)
				break
			}
		}
	}

	return filtered, nil
}

// Delete removes a configuration template by ID
func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM configuration_templates WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("template %s: %w", id, configuration.ErrTemplateNotFound)
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Count returns the number of templates (useful for testing)
func (r *SQLiteRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM configuration_templates").Scan(&count)
	return count, err
}

// Clear removes all templates (useful for testing)
func (r *SQLiteRepository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM configuration_templates")
	return err
}
