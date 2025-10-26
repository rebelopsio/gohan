.PHONY: build test test-unit test-integration test-all test-race test-wizard clean install dev-install lint

# Build configuration
BINARY_NAME=gohan
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/gohan

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

# Run all tests (unit + integration)
test-all:
	@echo "Running all tests..."
	@go test -v ./internal/...
	@go test -v -tags=integration ./tests/integration/...

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

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run

# Run the application
run:
	@go run ./cmd/gohan

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Display help
help:
	@echo "Gohan - Makefile commands:"
	@echo ""
	@echo "Build & Install:"
	@echo "  make build           Build the binary"
	@echo "  make install         Install to /usr/local/bin (requires sudo)"
	@echo "  make dev-install     Build for local development (no sudo)"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run all tests (unit + integration)"
	@echo "  make test-unit       Run unit tests only"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-all        Run both unit and integration tests"
	@echo "  make test-race       Run tests with race detector"
	@echo "  make test-coverage   Generate test coverage report"
	@echo "  make test-wizard     Run the preflight validation TUI wizard"
	@echo ""
	@echo "Development:"
	@echo "  make run             Run the application"
	@echo "  make fmt             Format code"
	@echo "  make tidy            Tidy dependencies"
	@echo "  make lint            Run linter (requires golangci-lint)"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean           Clean build artifacts"
