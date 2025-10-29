.PHONY: build build-all build-server test test-unit test-integration test-e2e test-all test-race test-wizard clean install dev-install lint
.PHONY: server dev-server watch docker-build docker-run ci pre-commit tools bench mod-verify install-hooks
.PHONY: act-ci act-lint act-test act-release
.PHONY: docs-venv docs-install docs-build docs-serve docs-deploy docs-clean

# Build configuration
BINARY_NAME=gohan
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Python venv configuration
VENV_DIR=.venv
PYTHON=$(VENV_DIR)/bin/python3
PIP=$(VENV_DIR)/bin/pip
MKDOCS=$(VENV_DIR)/bin/mkdocs

# Default target
.DEFAULT_GOAL := help

# Build the CLI binary
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/gohan

# Build all binaries (CLI with server support)
build-all: build
	@echo "All binaries built successfully!"

# Alias for consistency
build-server: build

# Install to /usr/local/bin
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	@sudo cp bin/${BINARY_NAME} /usr/local/bin/
	@echo "Installed! Run with: ${BINARY_NAME}"

# Install for local development (no sudo)
dev-install: build
	@echo "Binary built for local development: ./bin/${BINARY_NAME}"
	@echo "Run with: ./bin/${BINARY_NAME} or sudo ./bin/${BINARY_NAME} init"

# Run all tests (alias for test-all)
test: test-all

# Run unit tests only (excluding integration)
test-unit:
	@echo "Running unit tests..."
	@go test -v ./internal/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	@go test -v -tags=e2e ./tests/e2e/...

# Run all tests (unit + integration + e2e)
test-all:
	@echo "Running all tests..."
	@go test -v ./internal/...
	@go test -v -tags=integration ./tests/integration/...
	@go test -v -tags=e2e ./tests/e2e/...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race ./internal/...

# Run the TUI wizard for testing
test-wizard: build
	@echo "Running preflight validation wizard..."
	@sudo ./bin/${BINARY_NAME} init

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -rf site/

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run

# Run the CLI application
run:
	@go run ./cmd/gohan

# Start the API server
server: build
	@echo "Starting Gohan server..."
	@./bin/${BINARY_NAME} server

# Start the API server in development mode (with logging)
dev-server: build
	@echo "Starting Gohan server in dev mode..."
	@./bin/${BINARY_NAME} server --host 0.0.0.0 --port 8080

# Watch for changes and rebuild (requires entr: brew install entr)
watch:
	@echo "Watching for changes..."
	@find . -name "*.go" | entr -r make build

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Verify dependencies
mod-verify:
	@echo "Verifying dependencies..."
	@go mod verify

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Install development tools
tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo "Tools installed!"

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t gohan:${VERSION} .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 gohan:${VERSION}

# CI target (used in CI/CD pipeline)
ci: lint test-all
	@echo "CI checks passed!"

# Pre-commit checks
pre-commit: fmt lint test-unit
	@echo "Pre-commit checks passed!"

# Install git hooks
install-hooks:
	@echo "Installing git hooks..."
	@./scripts/install-hooks.sh

# Test CI locally with act (requires act: https://github.com/nektos/act)
act-ci:
	@echo "Running CI workflow locally with act..."
	@act -j lint -j test -j test-race -j build --artifact-server-path /tmp/artifacts

# Test lint job locally with act
act-lint:
	@echo "Running lint job locally with act..."
	@act -j lint

# Test test jobs locally with act
act-test:
	@echo "Running test jobs locally with act..."
	@act -j test

# Test release workflow locally with act
act-release:
	@echo "Testing release workflow locally with act..."
	@echo "Note: This requires a version tag. Use: git tag v0.1.0-test"
	@act release -e .github/workflows/release-event.json --artifact-server-path /tmp/artifacts

# ============================================================================
# Documentation targets
# ============================================================================

# Create Python virtual environment
docs-venv:
	@echo "Creating Python virtual environment..."
	@python3 -m venv $(VENV_DIR)
	@echo "Virtual environment created at $(VENV_DIR)"

# Install documentation dependencies
docs-install: docs-venv
	@echo "Installing documentation dependencies..."
	@$(PIP) install --upgrade pip
	@$(PIP) install mkdocs-material mkdocs-minify-plugin
	@echo "Documentation dependencies installed!"

# Build documentation with strict mode
docs-build: docs-install
	@echo "Building documentation..."
	@$(MKDOCS) build --strict
	@echo "Documentation built successfully in site/"

# Serve documentation locally
docs-serve: docs-install
	@echo "Starting documentation server..."
	@echo "Open http://127.0.0.1:8000 in your browser"
	@$(MKDOCS) serve

# Deploy documentation to GitHub Pages
docs-deploy: docs-install
	@echo "Deploying documentation to GitHub Pages..."
	@$(MKDOCS) gh-deploy --force
	@echo "Documentation deployed!"

# Clean documentation build artifacts
docs-clean:
	@echo "Cleaning documentation artifacts..."
	@rm -rf site/
	@rm -rf $(VENV_DIR)
	@echo "Documentation artifacts cleaned!"

# Display help
help:
	@echo "Gohan - Makefile commands:"
	@echo ""
	@echo "Build & Install:"
	@echo "  make build           Build the CLI binary"
	@echo "  make build-all       Build all binaries"
	@echo "  make install         Install to /usr/local/bin (requires sudo)"
	@echo "  make dev-install     Build for local development (no sudo)"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run all tests (unit + integration + e2e)"
	@echo "  make test-unit       Run unit tests only"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-e2e        Run E2E tests"
	@echo "  make test-all        Run all tests (unit + integration + e2e)"
	@echo "  make test-race       Run tests with race detector"
	@echo "  make test-coverage   Generate test coverage report"
	@echo "  make test-wizard     Run the preflight validation TUI wizard"
	@echo "  make bench           Run benchmarks"
	@echo ""
	@echo "Server:"
	@echo "  make server          Start the API server"
	@echo "  make dev-server      Start the API server in dev mode"
	@echo ""
	@echo "Development:"
	@echo "  make run             Run the CLI application"
	@echo "  make watch           Watch for changes and rebuild (requires entr)"
	@echo "  make fmt             Format code"
	@echo "  make tidy            Tidy dependencies"
	@echo "  make mod-verify      Verify dependencies"
	@echo "  make lint            Run linter (requires golangci-lint)"
	@echo "  make tools           Install development tools"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build    Build Docker image"
	@echo "  make docker-run      Run Docker container"
	@echo ""
	@echo "CI/CD:"
	@echo "  make ci              Run CI checks (lint + all tests)"
	@echo "  make pre-commit      Run pre-commit checks (fmt + lint + unit tests)"
	@echo "  make install-hooks   Install git hooks for pre-commit/pre-push"
	@echo ""
	@echo "Local CI Testing (requires act):"
	@echo "  make act-ci          Run full CI workflow locally"
	@echo "  make act-lint        Run lint job locally"
	@echo "  make act-test        Run test jobs locally"
	@echo "  make act-release     Test release workflow locally"
	@echo ""
	@echo "Documentation:"
	@echo "  make docs-venv       Create Python virtual environment"
	@echo "  make docs-install    Install MkDocs and dependencies in venv"
	@echo "  make docs-build      Build documentation (strict mode)"
	@echo "  make docs-serve      Serve documentation locally at http://127.0.0.1:8000"
	@echo "  make docs-deploy     Deploy documentation to GitHub Pages"
	@echo "  make docs-clean      Clean documentation artifacts and venv"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean           Clean build artifacts"
