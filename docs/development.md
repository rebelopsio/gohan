# Development Guide

Welcome to the Gohan development guide! This document covers everything you need to know to contribute to Gohan.

## Table of Contents
- [Development Setup](#development-setup)
- [Architecture Overview](#architecture-overview)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Standards](#code-standards)
- [Contributing](#contributing)

## Development Setup

### Prerequisites

- **Go 1.21+**: https://go.dev/dl/
- **Git**: For version control
- **Make**: For build automation (optional)
- **golangci-lint**: For linting (optional)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/rebelopsio/gohan.git
cd gohan

# Download dependencies
go mod download

# Build the project
go build -o gohan ./cmd/gohan

# Run tests
go test ./...

# Run with race detector
go test -race ./...
```

### Project Structure

```
gohan/
├── cmd/
│   └── gohan/          # Main application entry point
├── internal/
│   ├── application/    # Use cases (application layer)
│   ├── domain/         # Domain models and business logic
│   ├── infrastructure/ # External dependencies (DB, HTTP, etc.)
│   ├── cli/            # CLI commands
│   ├── config/         # Configuration management
│   ├── container/      # Dependency injection
│   └── tui/            # Terminal UI components
├── templates/          # Theme templates
├── docs/
│   └── features/       # BDD feature files (Gherkin)
├── tests/
│   ├── acceptance/     # ATDD acceptance tests
│   └── integration/    # Integration tests
├── wiki/               # Documentation (for GitHub wiki)
└── go.mod              # Go module definition
```

## Architecture Overview

Gohan follows **Clean Architecture** principles with clear separation of concerns:

### Layers

```
┌─────────────────────────────────────────┐
│   CLI Layer (cmd/gohan, internal/cli)   │
│   - Cobra commands                       │
│   - User interface                       │
├─────────────────────────────────────────┤
│   Application Layer (internal/application)│
│   - Use cases                            │
│   - Application services                 │
├─────────────────────────────────────────┤
│   Domain Layer (internal/domain)         │
│   - Business logic                       │
│   - Domain models                        │
│   - Repository interfaces                │
├─────────────────────────────────────────┤
│   Infrastructure Layer (internal/infrastructure)│
│   - Repository implementations          │
│   - External services                    │
│   - Template processing                  │
└─────────────────────────────────────────┘
```

### Key Patterns

1. **Dependency Injection**: Using container pattern
2. **Interface-Based Design**: Depend on abstractions, not concretions
3. **Repository Pattern**: Abstract data access
4. **Use Case Pattern**: Encapsulate business operations

## Development Workflow

### BDD → ATDD → TDD

Gohan uses a strict **Behavior-Driven Development** workflow:

```
1. Write Gherkin Feature (BDD)
   ↓
2. Write Acceptance Tests (ATDD)
   ↓
3. Write Unit Tests (TDD - RED)
   ↓
4. Implement Code (GREEN)
   ↓
5. Refactor (Keep tests GREEN)
```

### Example: Adding a New Feature

Let's add a "theme export" feature:

#### Step 1: Write BDD Feature

Create `docs/features/theme-export.feature`:
```gherkin
Feature: Theme Export
  As a user
  I want to export my current theme
  So that I can share it with others

  Scenario: Export active theme
    Given the active theme is "mocha"
    When I export the theme
    Then I should receive a JSON file
    And the file should contain theme colors
```

#### Step 2: Write Acceptance Test

Create `tests/acceptance/theme_export_test.go`:
```go
func TestThemeExport_ExportActiveTheme(t *testing.T) {
    t.Run("exports active theme to JSON", func(t *testing.T) {
        ctx := context.Background()

        // Given
        registry := theme.NewThemeRegistry()
        err := theme.InitializeStandardThemes(registry)
        require.NoError(t, err)

        // When
        useCase := theme.NewExportThemeUseCase(registry)
        result, err := useCase.Execute(ctx)

        // Then
        require.NoError(t, err)
        assert.NotEmpty(t, result.JSON)
        assert.Contains(t, result.JSON, "mocha")
    })
}
```

#### Step 3: Write Unit Tests

Create `internal/application/theme/export_theme_test.go`:
```go
func TestExportThemeUseCase_Execute(t *testing.T) {
    tests := []struct {
        name        string
        activeTheme string
        wantErr     bool
    }{
        {
            name:        "export mocha theme",
            activeTheme: "mocha",
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Step 4: Implement

Create `internal/application/theme/export_theme.go`:
```go
type ExportThemeUseCase struct {
    registry theme.ThemeRegistry
}

func (uc *ExportThemeUseCase) Execute(ctx context.Context) (*ExportResult, error) {
    // Implementation
}
```

#### Step 5: Add CLI Command

Update `internal/cli/cmd/theme.go`:
```go
var themeExportCmd = &cobra.Command{
    Use:   "export",
    Short: "Export current theme",
    RunE:  runThemeExport,
}
```

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/domain/theme/

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With race detector
go test -race ./...

# Integration tests only
go test -tags=integration ./tests/integration/...

# Acceptance tests
go test ./tests/acceptance/...
```

### Test Organization

```go
// Table-driven tests
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "test",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Coverage Goals

- **Domain Layer**: 90%+ coverage
- **Application Layer**: 85%+ coverage
- **Infrastructure Layer**: 75%+ coverage
- **CLI Layer**: 60%+ coverage (harder to test)

## Code Standards

### Go Style

Follow official Go guidelines:
- https://go.dev/doc/effective_go
- https://github.com/golang/go/wiki/CodeReviewComments

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

### Code Formatting

```bash
# Format code
go fmt ./...

# Import organization
goimports -w .
```

### Naming Conventions

```go
// Interfaces: Noun or Adjective + "er"
type ThemeRegistry interface
type CommandExecutor interface

// Structs: Descriptive nouns
type ThemeApplierImpl struct
type FileThemeStateStore struct

// Use Cases: Verb + Noun + "UseCase"
type ApplyThemeUseCase struct
type RollbackThemeUseCase struct

// Functions: Verb-first
func ApplyTheme()
func LoadConfiguration()

// Test functions: Test + FunctionName
func TestApplyTheme()
func TestLoadConfiguration()
```

### Error Handling

```go
// Wrap errors with context
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use custom errors
var ErrThemeNotFound = errors.New("theme not found")

// Check specific errors
if errors.Is(err, ErrThemeNotFound) {
    // Handle specific error
}
```

### Documentation

```go
// Package documentation
// Package theme provides theme management functionality.
package theme

// Public types and functions need documentation
// NewThemeRegistry creates a new theme registry.
func NewThemeRegistry() *ThemeRegistry {
    // Implementation
}
```

## Contributing

### Reporting Issues

1. Check existing issues first
2. Provide clear description
3. Include steps to reproduce
4. Add relevant logs/errors
5. Mention OS and versions

### Submitting Pull Requests

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/theme-export
   ```

3. **Follow the workflow**:
   - Write BDD feature
   - Write acceptance tests
   - Write unit tests (RED)
   - Implement (GREEN)
   - Refactor

4. **Ensure tests pass**:
   ```bash
   go test ./...
   go test -race ./...
   ```

5. **Update documentation**:
   - Update wiki if needed
   - Add code comments
   - Update CHANGELOG

6. **Commit with clear messages**:
   ```bash
   git add .
   git commit -m "feat: add theme export functionality"
   ```

7. **Push and create PR**:
   ```bash
   git push origin feature/theme-export
   ```
   Then create PR on GitHub

### Commit Message Format

```
<type>: <description>

[optional body]

[optional footer]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples**:
```
feat: add theme export command

Implements theme export functionality with JSON output.
Includes unit tests and acceptance tests.

Closes #123
```

### Code Review Process

1. **Automated checks run**: Tests, linting
2. **Maintainer review**: Code quality, design
3. **Address feedback**: Make requested changes
4. **Approval**: PR is approved
5. **Merge**: Changes merged to main

## Building Releases

### Version Numbering

Gohan follows Semantic Versioning (semver):
- **Major**: Breaking changes (v2.0.0)
- **Minor**: New features (v1.1.0)
- **Patch**: Bug fixes (v1.0.1)

### Creating a Release

```bash
# Tag the release
git tag -a v1.2.0 -m "Release v1.2.0"

# Push tag
git push origin v1.2.0

# GitHub Actions will:
# - Run tests
# - Build binaries
# - Create GitHub release
# - Upload artifacts
```

## Resources

- **Go Documentation**: https://go.dev/doc/
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- **BDD with Go**: https://github.com/cucumber/godog
- **Testing Best Practices**: https://go.dev/doc/effective_go#testing

## Getting Help

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community
- **Code Review**: Ask in PRs for feedback

Thank you for contributing to Gohan! 🚀
