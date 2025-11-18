.PHONY: all build test fmt vet lint check-lint install-lint run deps tidy

# Default target
all: tidy fmt vet lint test build

build:
	@echo "Building diploma.exe..."
	go build -o diploma.exe .

test:
	@echo "Running tests..."
	go test ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

lint: check-lint
	@echo "Running linters..."
	golangci-lint run --fix --config .golangci.yml

check-lint:
	@echo "Checking golangci-lint installation..."
	@which golangci-lint > /dev/null || (echo "ERROR: golangci-lint is not installed. Run 'make install-lint' to install it." && exit 1)
	@echo "Checking golangci-lint version..."
	@INSTALLED_VERSION=$$(golangci-lint --version 2>/dev/null | grep -oP 'version \K\d+\.\d+\.\d+' || echo "unknown"); \
	LATEST_VERSION=$$(curl -s https://api.github.com/repos/golangci/golangci-lint/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/' || echo "unknown"); \
	if [ "$$INSTALLED_VERSION" = "unknown" ] || [ "$$LATEST_VERSION" = "unknown" ]; then \
		echo "WARNING: Could not verify version. Installed: $$INSTALLED_VERSION"; \
	elif [ "$$INSTALLED_VERSION" != "$$LATEST_VERSION" ]; then \
		echo "ERROR: golangci-lint is outdated!"; \
		echo "  Installed: v$$INSTALLED_VERSION"; \
		echo "  Latest:    v$$LATEST_VERSION"; \
		echo "  Run 'make install-lint' to update."; \
		exit 1; \
	else \
		echo "golangci-lint is up-to-date (v$$INSTALLED_VERSION)"; \
	fi

install-lint:
	@echo "Installing latest golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest

run: build
	@echo "Running diploma.exe..."
	./diploma.exe

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
