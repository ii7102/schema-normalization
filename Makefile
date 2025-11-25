.PHONY: all build test fmt vet lint check-lint install-lint run deps tidy

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RM := del /q
    MKDIR := mkdir
    WHICH := where
    DEVNULL := nul
    SEP := ;
    EXE := .exe
else
    DETECTED_OS := $(shell uname -s)
    RM := rm -f
    MKDIR := mkdir -p
    WHICH := which
    DEVNULL := /dev/null
    SEP := &&
    EXE :=
endif

# Default target
all: tidy fmt vet lint test build

build:
	@echo "Building diploma.exe..."
	go build -o diploma.exe .
	@$(RM) diploma.exe 2>$(DEVNULL) || true

test:
	@echo "Running tests..."
	go test -count=1 ./...

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
ifeq ($(OS),Windows_NT)
	@golangci-lint version >$(DEVNULL) 2>$(DEVNULL) || (echo ERROR: golangci-lint is not installed. Run 'make install-lint' to install it. && exit 1)
else
	@golangci-lint --version >$(DEVNULL) 2>$(DEVNULL) || (echo "ERROR: golangci-lint is not installed. Run 'make install-lint' to install it." && exit 1)
endif
	@echo "Checking golangci-lint version..."
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "$$INSTALLED = (golangci-lint version 2>$$null | Select-String -Pattern 'version (\d+\.\d+\.\d+)' | ForEach-Object { $$_.Matches.Groups[1].Value }); \
	if (-not $$INSTALLED) { $$INSTALLED = 'unknown' }; \
	try { $$LATEST = (Invoke-RestMethod -Uri 'https://api.github.com/repos/golangci/golangci-lint/releases/latest').tag_name -replace '^v', '' } catch { $$LATEST = 'unknown' }; \
	if ($$INSTALLED -eq 'unknown' -or $$LATEST -eq 'unknown') { \
		Write-Host \"WARNING: Could not verify version. Installed: $$INSTALLED\" \
	} elseif ($$INSTALLED -ne $$LATEST) { \
		Write-Host \"ERROR: golangci-lint is outdated!\"; \
		Write-Host \"  Installed: v$$INSTALLED\"; \
		Write-Host \"  Latest:    v$$LATEST\"; \
		Write-Host \"  Run 'make install-lint' to update.\"; \
		exit 1 \
	} else { \
		Write-Host \"golangci-lint is up-to-date (v$$INSTALLED)\" \
	}"
else
	@INSTALLED_VERSION=$$(golangci-lint --version 2>$(DEVNULL) | grep -oP 'version \K\d+\.\d+\.\d+' || echo "unknown"); \
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
endif

install-lint:
	@echo "Installing latest golangci-lint..."
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "$$GOPATH = (go env GOPATH); \
	$$GOPATH_BIN = Join-Path $$GOPATH 'bin'; \
	if (-not (Test-Path $$GOPATH_BIN)) { New-Item -ItemType Directory -Path $$GOPATH_BIN -Force | Out-Null }; \
	$$LATEST = (Invoke-RestMethod -Uri 'https://api.github.com/repos/golangci/golangci-lint/releases/latest').tag_name -replace '^v', ''; \
	$$ARCH = 'amd64'; $$OS = 'windows'; \
	$$URL = \"https://github.com/golangci/golangci-lint/releases/latest/download/golangci-lint-$$LATEST-$$OS-$$ARCH.zip\"; \
	$$ZIP = Join-Path $$env:TEMP 'golangci-lint.zip'; \
	Invoke-WebRequest -Uri $$URL -OutFile $$ZIP; \
	Expand-Archive -Path $$ZIP -DestinationPath $$env:TEMP -Force; \
	$$EXTRACTED = Join-Path $$env:TEMP \"golangci-lint-$$LATEST-$$OS-$$ARCH\"; \
	Copy-Item (Join-Path $$EXTRACTED 'golangci-lint.exe') -Destination $$GOPATH_BIN -Force; \
	Remove-Item $$ZIP -ErrorAction SilentlyContinue; \
	Remove-Item $$EXTRACTED -Recurse -Force -ErrorAction SilentlyContinue; \
	Write-Host \"Installed golangci-lint v$$LATEST to $$GOPATH_BIN\""
else
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
endif

run: build
	@echo "Running diploma.exe..."
ifeq ($(OS),Windows_NT)
	@if exist diploma.exe (diploma.exe) else (.\diploma.exe)
else
	@./diploma.exe
endif

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
