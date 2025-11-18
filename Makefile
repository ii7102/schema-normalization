.PHONY: all build test fmt vet lint install-lint run deps tidy

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

lint: install-lint
	@echo "Running linters..."
	golangci-lint run --fix --config .golangci.yml

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
