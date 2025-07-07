# CLI Package Makefile

.PHONY: test test-unit test-integration build clean fmt vet lint help

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires Docker)"
	@echo "  build         - Build the CLI binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format Go code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint (if available)"
	@echo "  help          - Show this help"

# Run all tests
test: fmt vet
	@echo "Running all tests..."
	go test -v ./...

# Run unit tests only (short mode)
test-unit: fmt vet
	@echo "Running unit tests..."
	go test -v -short ./...

# Run integration tests (requires Docker)
test-integration:
	@echo "Running integration tests..."
	go test -v -run "Integration" ./...

# Build the CLI binary
build:
	@echo "Building CLI binary..."
	go build -o codequest .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f codequest
	go clean

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint if available
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

# Test coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. ./...