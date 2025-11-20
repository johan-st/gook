.PHONY: build install test test-verbose test-coverage clean publish help

# Variables
BINARY_NAME=gook
CMD_PATH=./cmd
MODULE_PATH=$(shell go list -m)
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.DEFAULT_GOAL := help

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(CMD_PATH)

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(CMD_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

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

# Publish: tag version and push (requires VERSION variable)
publish: test
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION must be set (e.g., make publish VERSION=v1.0.0)"; \
		exit 1; \
	fi
	@echo "Publishing version $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)" || true
	@git push origin $(VERSION) || true
	@GOPROXY=proxy.golang.org go list -m $(MODULE_PATH)@$(VERSION) || echo "Note: Module will be available after proxy sync"

# Help target
help:
	@echo "Available targets:"
	@echo "  make build          - Build the binary to bin/$(BINARY_NAME)"
	@echo "  make install        - Install the binary to GOPATH/bin"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests and generate coverage report"
	@echo "  make clean          - Remove build artifacts and coverage files"
	@echo "  make publish        - Tag and publish a version (requires VERSION=v1.0.0)"
	@echo "  make help           - Show this help message"

