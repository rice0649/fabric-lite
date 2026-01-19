.PHONY: build build-forge build-all install install-forge install-all clean test lint run help

BINARY_NAME=fabric-lite
VERSION?=0.1.0
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -type f)

# Build fabric-lite
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

# Build forge CLI
build-forge:
	@echo "Building forge..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/forge ./cmd/forge
	@echo "Built: $(BUILD_DIR)/forge"

# Build all binaries
build-all: build build-forge
	@echo "All binaries built"

# Install fabric-lite to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/ 2>/dev/null || cp $(BUILD_DIR)/$(BINARY_NAME) ~/go/bin/
	@echo "Installed to ~/go/bin/$(BINARY_NAME)"

# Install forge to GOPATH/bin
install-forge: build-forge
	@echo "Installing forge..."
	cp $(BUILD_DIR)/forge $(GOPATH)/bin/ 2>/dev/null || cp $(BUILD_DIR)/forge ~/go/bin/
	@echo "Installed to ~/go/bin/forge"

# Install patterns to user config
install-patterns:
	@echo "Installing patterns..."
	@mkdir -p ~/.config/fabric-lite/patterns
	cp -r patterns/* ~/.config/fabric-lite/patterns/
	@echo "Patterns installed to ~/.config/fabric-lite/patterns/"

# Full install (all binaries + patterns)
install-all: install install-forge install-patterns
	@echo "Full installation complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run || echo "Install golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run the application (for development)
run:
	go run ./cmd/$(BINARY_NAME) $(ARGS)

# Show help
help:
	@echo "fabric-lite & forge Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build            - Build fabric-lite"
	@echo "  make build-forge      - Build forge CLI"
	@echo "  make build-all        - Build all binaries"
	@echo "  make install          - Install fabric-lite to ~/go/bin"
	@echo "  make install-forge    - Install forge to ~/go/bin"
	@echo "  make install-patterns - Install patterns to ~/.config/fabric-lite"
	@echo "  make install-all      - Full install (all binaries + patterns)"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make lint             - Run linter"
	@echo "  make fmt              - Format code"
	@echo "  make run ARGS='...'   - Run fabric-lite with arguments"
	@echo "  make help             - Show this help"

.DEFAULT_GOAL := help
