# GNUmakefile

# Ensure Make is run with bash shell as some commands in this Makefile require bash
SHELL := /bin/bash

# Use asdf to ensure correct tool versions
ASDF := $(shell command -v asdf 2> /dev/null)
ifdef ASDF
	# Set up asdf environment
	export ASDF_DIR := $(shell asdf where)
	export PATH := $(ASDF_DIR)/shims:$(PATH)
endif

# The name of the binary (default is current directory name)
TARGET := terraform-provider-peekaping
# These will be provided to the binary
VERSION ?= 0.0.1
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# These are the values we want to pass for Version and BuildTime
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# List of GOOS and GOARCH pairs for cross-compilation
XC_OS := darwin linux windows
XC_ARCH := amd64 arm64
XC_ARCH_ARM64 := arm64
XC_ARCH_AMD64 := amd64

# Default target
.PHONY: all
all: build

# Build the provider binary
.PHONY: build
build: asdf-install
	@echo "==> Building $(TARGET)..."
	@CGO_ENABLED=0 go build -o $(TARGET) $(LDFLAGS) .

# Install the provider binary
.PHONY: install
install: build
	@echo "==> Installing $(TARGET)..."
	@go install $(LDFLAGS) .

# Run tests
.PHONY: test
test: asdf-install
	@echo "==> Running tests..."
	@go test -v ./...

# Run acceptance tests
.PHONY: testacc
testacc: asdf-install
	@echo "==> Running acceptance tests..."
	@TF_ACC=1 go test -v ./...

# Run tests with race detection
.PHONY: testrace
testrace: asdf-install
	@echo "==> Running tests with race detection..."
	@go test -race -v ./...

# Run tests with coverage
.PHONY: testcover
testcover: asdf-install
	@echo "==> Running tests with coverage..."
	@go test -coverprofile=coverage.out -v ./...
	@go tool cover -html=coverage.out -o coverage.html

# Run linting
.PHONY: lint
lint: asdf-install
	@echo "==> Running linters..."
	@golangci-lint run

# Format code
.PHONY: fmt
fmt: asdf-install
	@echo "==> Formatting code..."
	@go fmt ./...

# Run go mod tidy
.PHONY: mod
mod: asdf-install
	@echo "==> Running go mod tidy..."
	@go mod tidy

# Generate documentation
.PHONY: docs
docs: asdf-install
	@echo "==> Generating documentation..."
	@tfplugindocs generate

# Ensure asdf tools are installed
.PHONY: asdf-install
asdf-install:
	@echo "==> Installing asdf tools..."
	@if command -v asdf >/dev/null 2>&1; then \
		asdf install; \
	else \
		echo "asdf not found. Please install asdf first: https://asdf-vm.com/"; \
		exit 1; \
	fi

# Install development tools
.PHONY: tools
tools: asdf-install
	@echo "==> Installing development tools..."
	@cd tools && go generate

# Run copywrite
.PHONY: copywrite
copywrite: asdf-install
	@echo "==> Running copywrite..."
	@copywrite headers --plan

# Fix copywrite
.PHONY: copywrite-fix
copywrite-fix: asdf-install
	@echo "==> Fixing copywrite headers..."
	@copywrite headers

# Build Docker image
.PHONY: docker
docker:
	@echo "==> Building Docker image..."
	@docker build -t terraform-provider-peekaping .

# Run Docker container
.PHONY: docker-run
docker-run:
	@echo "==> Running Docker container..."
	@docker run --rm terraform-provider-peekaping

# Clean build artifacts
.PHONY: clean
clean:
	@echo "==> Cleaning..."
	@rm -f $(TARGET)
	@rm -f coverage.out coverage.html
	@go clean

# Cross-compile for multiple platforms
.PHONY: xc
xc:
	@echo "==> Cross-compiling..."
	@for os in $(XC_OS); do \
		for arch in $(XC_ARCH); do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then \
				continue; \
			fi; \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -o $(TARGET)-$$os-$$arch $(LDFLAGS) .; \
		done; \
	done

# Validate Terraform configurations
.PHONY: validate
validate:
	@echo "==> Validating Terraform configurations..."
	@terraform validate

# Validate examples
.PHONY: validate-examples
validate-examples:
	@echo "==> Validating examples..."
	@for example in examples/*/main.tf; do \
		if [ -f "$$example" ]; then \
			echo "Validating $$example..."; \
			terraform fmt -check "$$example" > /dev/null 2>&1 || (echo "❌ $$example has syntax issues" && exit 1); \
			echo "✅ $$example syntax is valid"; \
		fi; \
	done

# Run all checks
.PHONY: check
check: fmt lint test validate validate-examples

# Run all checks including acceptance tests
.PHONY: checkacc
checkacc: fmt lint test testacc validate validate-examples

# Run comprehensive validation
.PHONY: validate-all
validate-all: check
	@echo "==> Running comprehensive validation..."
	@echo "✅ Provider builds successfully"
	@echo "✅ Go tests pass"
	@echo "✅ Linting passes"
	@echo "✅ Terraform configuration is valid"
	@echo "✅ Examples syntax is correct"
	@echo ""
	@echo "🎉 All validations completed successfully!"
	@echo "The Peekaping Terraform provider is ready for use."

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build the provider binary"
	@echo "  install         - Install the provider binary"
	@echo "  test            - Run tests"
	@echo "  testacc         - Run acceptance tests"
	@echo "  testrace        - Run tests with race detection"
	@echo "  testcover       - Run tests with coverage"
	@echo "  lint            - Run linters"
	@echo "  fmt             - Format code"
	@echo "  mod             - Run go mod tidy"
	@echo "  docs            - Generate documentation"
	@echo "  clean           - Clean build artifacts"
	@echo "  xc              - Cross-compile for multiple platforms"
	@echo "  validate        - Validate Terraform configuration"
	@echo "  validate-examples - Validate example configurations"
	@echo "  check           - Run all checks (fmt, lint, test, validate)"
	@echo "  checkacc        - Run all checks including acceptance tests"
	@echo "  validate-all    - Run comprehensive validation"
	@echo "  asdf-install    - Install asdf tool versions"
	@echo "  tools           - Install development tools via go generate"
	@echo "  copywrite       - Check copywrite headers"
	@echo "  copywrite-fix   - Fix copywrite headers"
	@echo "  docker          - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  help            - Show this help message"
