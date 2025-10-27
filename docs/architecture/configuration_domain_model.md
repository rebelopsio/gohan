# Configuration Domain Model Design

## Executive Summary

This document outlines the domain model design for saved/reusable installation configurations in the Gohan system. The design follows Domain-Driven Design (DDD) principles with Clean Architecture, using aggregates, value objects, and repository patterns to manage configuration templates that can be saved, shared, and reused across multiple installations.

## Core Concepts & Ubiquitous Language

### Key Domain Terms

- **Configuration Template**: A saved, reusable blueprint for installations that defines what components to install and how
- **Configuration Manifest**: The immutable core definition of what packages and settings constitute a configuration
- **Configuration Metadata**: Descriptive information about a configuration (name, description, tags, author)
- **Configuration Source**: Where a configuration originated from (user-created, exported, imported, system-default)
- **Configuration Validation**: Process of ensuring a configuration is well-formed and installable
- **Configuration Composition**: Combining multiple configurations into a single installation plan
- **Configuration Export**: Creating a reusable configuration from a completed installation
- **Configuration Preview**: Dry-run analysis showing what would be installed without making changes

## Aggregate Design

### Primary Aggregate: `ConfigurationTemplate`

The `ConfigurationTemplate` is the aggregate root that represents a saved, reusable installation configuration. It maintains consistency across the configuration's lifecycle and enforces business invariants.

```go
// ConfigurationTemplate is the aggregate root for saved configurations
type ConfigurationTemplate struct {
    id          string                    // Unique identifier (UUID)
    metadata    ConfigurationMetadata     // Name, description, tags, etc.
    manifest    ConfigurationManifest     // The actual configuration content
    source      ConfigurationSource       // How this configuration was created
    validation  *ConfigurationValidation  // Last validation result (cached)
    createdAt   time.Time
    updatedAt   time.Time
    version     int                       // For optimistic locking
}
```

### Aggregate Boundaries

The `ConfigurationTemplate` aggregate includes:
- **ConfigurationMetadata** (Value Object) - Descriptive information
- **ConfigurationManifest** (Value Object) - The immutable configuration definition
- **ConfigurationSource** (Value Object) - Origin tracking
- **ConfigurationValidation** (Value Object) - Cached validation results

What stays **outside** the aggregate:
- Installation sessions using this configuration (referenced by ID)
- Actual package repositories and versions (validated at install time)
- System requirements checking (done during preview/installation)

## Value Objects

### 1. ConfigurationMetadata

Encapsulates all descriptive information about a configuration.

```go
type ConfigurationMetadata struct {
    name        ConfigurationName  // Required, unique within a namespace
    description string            // Optional detailed description
    author      string            // Who created this configuration
    tags        []ConfigurationTag // Searchable tags (e.g., "development", "production")
    category    ConfigurationCategory // Type classification
}

type ConfigurationName struct {
    value string // Validated: non-empty, alphanumeric + dash/underscore, max 100 chars
}

type ConfigurationTag struct {
    value string // Validated: lowercase, alphanumeric, max 30 chars
}

type ConfigurationCategory string
const (
    CategoryDevelopment ConfigurationCategory = "development"
    CategoryProduction  ConfigurationCategory = "production"
    CategoryTesting     ConfigurationCategory = "testing"
    CategoryCustom      ConfigurationCategory = "custom"
)
```

### 2. ConfigurationManifest

The immutable core of what gets installed. This maps closely to the existing `InstallationConfiguration` but adds version constraints and requirements.

```go
type ConfigurationManifest struct {
    components          []ComponentSpecification  // What to install
    requirements        SystemRequirements       // Minimum system requirements
    serviceDirectives   []ServiceDirective       // Services to stop/start
    compatibilityRules  []CompatibilityRule     // Conditional installation rules
}

type ComponentSpecification struct {
    name              ComponentName      // From existing domain
    versionConstraint VersionConstraint  // e.g., ">=1.2.0", "~1.2.0", "latest"
    optional          bool              // Can be skipped if not available
    dependencies      []ComponentName   // Must be installed together
}

type VersionConstraint struct {
    operator string // ">=", "<=", "=", "~" (compatible with)
    version  string // Semantic version or "latest"
}

type SystemRequirements struct {
    minimumDiskSpaceMB   uint64
    minimumMemoryMB      uint64
    requiredGPUVendors   []string // "nvidia", "amd", "intel", or empty for any
    incompatiblePackages []string // Packages that must not be present
}

type ServiceDirective struct {
    serviceName string
    action      ServiceAction // stop_before, start_after, restart_after
    required    bool         // If false, failure is non-fatal
}
```

### 3. ConfigurationSource

Tracks how a configuration was created for auditability and trust.

```go
type ConfigurationSource struct {
    sourceType    SourceType
    sourceID      string    // Installation ID if exported, URL if imported, etc.
    importedFrom  string    // Original location if imported
    exportedAt    time.Time // When exported from an installation
}

type SourceType string
const (
    SourceUserCreated    SourceType = "user_created"    // Manually created
    SourceExported       SourceType = "exported"        // From installation
    SourceImported       SourceType = "imported"        // From file/URL
    SourceSystemDefault  SourceType = "system_default"  // Bundled with system
)
```

### 4. ConfigurationValidation

Cached validation results to avoid re-validating unchanged configurations.

```go
type ConfigurationValidation struct {
    isValid          bool
    validatedAt      time.Time
    errors           []ValidationError
    warnings         []ValidationWarning
    previewSummary   *InstallationPreview  // What would be installed
}

type ValidationError struct {
    field   string
    message string
    code    string // For programmatic handling
}

type InstallationPreview struct {
    totalPackages      int
    totalSizeMB        float64
    estimatedDuration  time.Duration
    componentsToInstall []ComponentPreview
}
```

## Repository Interface

```go
type ConfigurationTemplateRepository interface {
    // Create saves a new configuration template
    Create(ctx context.Context, template *ConfigurationTemplate) error

    // Update modifies an existing configuration template
    Update(ctx context.Context, template *ConfigurationTemplate) error

    // FindByID retrieves a configuration by its unique ID
    FindByID(ctx context.Context, id string) (*ConfigurationTemplate, error)

    // FindByName retrieves a configuration by its unique name
    FindByName(ctx context.Context, name ConfigurationName) (*ConfigurationTemplate, error)

    // List retrieves configurations with optional filtering
    List(ctx context.Context, filter ConfigurationFilter) ([]*ConfigurationTemplate, error)

    // Delete removes a configuration template
    Delete(ctx context.Context, id string) error

    // ExistsByName checks if a configuration with the given name exists
    ExistsByName(ctx context.Context, name ConfigurationName) (bool, error)
}

type ConfigurationFilter struct {
    Category *ConfigurationCategory
    Tags     []ConfigurationTag
    Author   *string
    Limit    int
    Offset   int
}
```

## Domain Services

### 1. ConfigurationExportService

Exports completed installations as reusable configurations.

```go
type ConfigurationExportService struct {
    installationRepo InstallationSessionRepository
    configRepo      ConfigurationTemplateRepository
}

func (s *ConfigurationExportService) ExportFromInstallation(
    ctx context.Context,
    installationID string,
    metadata ConfigurationMetadata,
) (*ConfigurationTemplate, error) {
    // 1. Load the installation session
    // 2. Extract installed components and their versions
    // 3. Create ConfigurationManifest from installation
    // 4. Build ConfigurationTemplate
    // 5. Validate and save
}
```

### 2. ConfigurationCompositionService

Combines multiple configurations into a single installation plan.

```go
type ConfigurationCompositionService struct {
    configRepo ConfigurationTemplateRepository
}

func (s *ConfigurationCompositionService) Compose(
    ctx context.Context,
    configIDs []string,
) (InstallationConfiguration, error) {
    // 1. Load all configurations
    // 2. Check for conflicts
    // 3. Merge components (union)
    // 4. Combine requirements (max of each)
    // 5. Merge service directives
    // 6. Return combined InstallationConfiguration
}
```

### 3. ConfigurationValidationService

Validates configurations before use.

```go
type ConfigurationValidationService struct {
    packageResolver PackageResolver // External service to check package availability
}

func (s *ConfigurationValidationService) Validate(
    ctx context.Context,
    template *ConfigurationTemplate,
) (*ConfigurationValidation, error) {
    // 1. Validate manifest structure
    // 2. Check component availability
    // 3. Verify version constraints can be satisfied
    // 4. Check for conflicts
    // 5. Return validation result
}

func (s *ConfigurationValidationService) Preview(
    ctx context.Context,
    template *ConfigurationTemplate,
    systemSnapshot *SystemSnapshot,
) (*InstallationPreview, error) {
    // 1. Validate configuration
    // 2. Check against current system state
    // 3. Calculate what would be installed
    // 4. Estimate size and duration
    // 5. Return preview
}
```

## Domain Invariants

The following business rules must be enforced:

1. **Configuration Name Uniqueness**: No two configurations can have the same name
2. **Manifest Immutability**: Once created, a manifest cannot be modified (create new version instead)
3. **Valid Component References**: All components in manifest must be valid ComponentName values
4. **Non-Empty Components**: A configuration must specify at least one component
5. **Core Component Requirement**: If any Hyprland component is included, the core must be included
6. **Version Constraint Format**: Version constraints must follow semantic versioning rules
7. **Service Directive Ordering**: Stop directives must come before start directives
8. **Tag Limit**: Maximum 10 tags per configuration
9. **Description Length**: Maximum 1000 characters for description

## Integration with Existing Domain

### Relationship to InstallationConfiguration

`ConfigurationTemplate` serves as a persistent, reusable template that can be converted to an `InstallationConfiguration`:

```go
// Factory method on ConfigurationTemplate
func (t *ConfigurationTemplate) ToInstallationConfiguration(
    gpuSupport *GPUSupport,
    diskSpace DiskSpace,
    mergeExisting bool,
) (InstallationConfiguration, error) {
    // Convert ComponentSpecifications to ComponentSelections
    // Apply version resolution
    // Return InstallationConfiguration
}
```

### Relationship to InstallationSession

Installation sessions track which configuration template was used:

```go
type InstallationSession struct {
    // ... existing fields ...
    configurationTemplateID *string // Optional reference to template used
}
```

### Event Integration

New domain events for configuration management:

```go
type ConfigurationTemplateCreated struct {
    TemplateID string
    Name       string
    CreatedBy  string
    OccurredAt time.Time
}

type ConfigurationTemplateDeleted struct {
    TemplateID string
    DeletedBy  string
    OccurredAt time.Time
}

type ConfigurationExportedFromInstallation struct {
    InstallationID string
    TemplateID     string
    ExportedBy     string
    OccurredAt     time.Time
}
```

## Usage Examples

### Creating a Configuration Template

```go
// Create metadata
metadata, err := NewConfigurationMetadata(
    "web-stack-v2",
    "Production web application stack",
    "admin@example.com",
    []string{"production", "web", "docker"},
    CategoryProduction,
)

// Create manifest
manifest := NewConfigurationManifest(
    []ComponentSpecification{
        {ComponentDocker, ">=20.0.0", false, nil},
        {ComponentPostgreSQL, "~14.0", false, nil},
        {ComponentRedis, "latest", false, nil},
    },
    SystemRequirements{
        MinimumDiskSpaceMB: 10000,
        MinimumMemoryMB: 4096,
    },
    []ServiceDirective{
        {"docker", StopBefore, true},
        {"docker", StartAfter, true},
    },
)

// Create template
template := NewConfigurationTemplate(metadata, manifest, SourceUserCreated)
```

### Exporting from Installation

```go
exportService := NewConfigurationExportService(installRepo, configRepo)
template, err := exportService.ExportFromInstallation(
    ctx,
    "installation-123",
    metadata,
)
```

### Composing Multiple Configurations

```go
composeService := NewConfigurationCompositionService(configRepo)
combined, err := composeService.Compose(ctx, []string{
    "base-config-id",
    "project-specific-id",
})
```

## Migration Path

To integrate this design with the existing codebase:

1. **Phase 1**: Implement value objects and aggregate
2. **Phase 2**: Add repository and basic CRUD operations
3. **Phase 3**: Implement export service for creating templates from installations
4. **Phase 4**: Add validation and preview services
5. **Phase 5**: Implement composition service for combining configurations
6. **Phase 6**: Update UI to support configuration management

## Summary

This domain model provides:

- **Reusability**: Save and share installation configurations
- **Composability**: Combine multiple configurations
- **Validation**: Ensure configurations are valid before use
- **Traceability**: Track configuration origin and usage
- **Flexibility**: Support version constraints and optional components
- **Safety**: Preview changes before installation
- **Integration**: Works seamlessly with existing `InstallationConfiguration` and `InstallationSession`

The design maintains clean boundaries, follows DDD principles, and provides a solid foundation for the configuration management feature.