.PHONY: build test test-integration test-race clean install lint

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
	@sudo mv bin/${BINARY_NAME} /usr/local/bin/

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

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
	@echo "  make build           Build the binary"
	@echo "  make install         Install to /usr/local/bin"
	@echo "  make test            Run all tests"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-race       Run tests with race detector"
	@echo "  make test-coverage   Generate test coverage report"
	@echo "  make clean           Clean build artifacts"
	@echo "  make lint            Run linter"
	@echo "  make run             Run the application"
	@echo "  make fmt             Format code"
	@echo "  make tidy            Tidy dependencies"
