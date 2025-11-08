.PHONY: all build test test-verbose lint fmt vet clean run help

# Default target
all: tidy fmt vet lint test build

# Build the application
build:
	@echo "Building diploma.exe..."
	go build -o diploma.exe .

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linters (requires golangci-lint to be installed)
lint:
	@echo "Running linters..."
	golangci-lint run --config .golangci.yml

# Run linters with verbose output (shows which linters are enabled)
lint-verbose:
	@echo "Running linters (verbose)..."
	golangci-lint linters
	golangci-lint run --config .golangci.yml --verbose

# Run the application
run: build
	@echo "Running diploma.exe..."
	./diploma.exe

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Run quick checks (format, vet, test)
check: fmt vet test
	@echo "All checks passed!"