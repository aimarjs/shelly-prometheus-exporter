.PHONY: build test lint fmt clean docker-build docker-run help

# Variables
BINARY_NAME=shelly-exporter
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build directory
BUILD_DIR=build
BIN_DIR=bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/shelly-exporter

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/shelly-exporter
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/shelly-exporter
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/shelly-exporter
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/shelly-exporter
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/shelly-exporter

test: ## Run tests
	@echo "Running tests..."
	CGO_ENABLED=0 $(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	CGO_ENABLED=0 $(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-coverage-lcov: ## Run tests with coverage in LCOV format
	@echo "Running tests with LCOV coverage..."
	CGO_ENABLED=0 $(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -func=coverage.out

test-coverage-relative: ## Generate coverage with relative paths for Qlty
	@echo "Generating coverage with relative paths..."
	CGO_ENABLED=0 $(GOTEST) -coverprofile=coverage.out -covermode=atomic ./...
	sed 's|github.com/aimar/shelly-prometheus-exporter/||g' coverage.out > coverage_relative.out
	@echo "Coverage file with relative paths: coverage_relative.out"

test-race: ## Run tests with race detector (Linux only)
	@echo "Running tests with race detector..."
	@echo "Note: Race detector requires CGO and may not work on macOS"
	CGO_ENABLED=1 $(GOTEST) -race -v ./...

lint: ## Run linter
	@echo "Running linter..."
	$(GOLINT) run

qlty-metrics: ## Run Qlty code quality metrics
	@echo "Running Qlty code quality metrics..."
	~/.qlty/bin/qlty metrics --all --exclude-tests

qlty-coverage-dry-run: test-coverage-relative ## Test Qlty coverage upload (dry run)
	@echo "Testing Qlty coverage upload (dry run)..."
	~/.qlty/bin/qlty coverage publish --dry-run --override-commit-sha=$$(git rev-parse HEAD) --override-branch=$$(git branch --show-current) coverage_relative.out

qlty-coverage-upload: test-coverage-relative ## Upload coverage to Qlty (requires QLTY_COVERAGE_TOKEN)
	@echo "Uploading coverage to Qlty..."
	@if [ -z "$$QLTY_COVERAGE_TOKEN" ]; then echo "Error: QLTY_COVERAGE_TOKEN environment variable is required"; exit 1; fi
	~/.qlty/bin/qlty coverage publish --override-commit-sha=$$(git rev-parse HEAD) --override-branch=$$(git branch --show-current) coverage_relative.out

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(BIN_DIR)
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -d --name $(BINARY_NAME) -p 8080:8080 $(BINARY_NAME):$(VERSION)

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	cp $(BIN_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BIN_DIR)/$(BINARY_NAME)

dev: ## Run in development mode with hot reload
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular run..."; \
		make run; \
	fi

release: ## Create a release (requires goreleaser)
	@echo "Creating release..."
	goreleaser release --clean

release-snapshot: ## Create a snapshot release
	@echo "Creating snapshot release..."
	goreleaser release --snapshot --clean

version-patch: ## Bump patch version and create release
	@echo "Bumping patch version..."
	@git add -A
	@git commit -m "chore: bump patch version" || true
	@git push origin main
	@echo "Triggering patch release workflow..."
	@gh workflow run release.yml -f version_type=patch

version-minor: ## Bump minor version and create release
	@echo "Bumping minor version..."
	@git add -A
	@git commit -m "chore: bump minor version" || true
	@git push origin main
	@echo "Triggering minor release workflow..."
	@gh workflow run release.yml -f version_type=minor

version-major: ## Bump major version and create release
	@echo "Bumping major version..."
	@git add -A
	@git commit -m "chore: bump major version" || true
	@git push origin main
	@echo "Triggering major release workflow..."
	@gh workflow run release.yml -f version_type=major

version-prerelease: ## Create prerelease version
	@echo "Creating prerelease..."
	@git add -A
	@git commit -m "chore: create prerelease" || true
	@git push origin main
	@echo "Triggering prerelease workflow..."
	@gh workflow run release.yml -f version_type=prerelease

current-version: ## Show current version
	@echo "Current version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build time: $(BUILD_TIME)"

next-version: ## Show what the next version would be
	@echo "Next version would be determined by semantic-release based on commit messages"
	@echo "Use 'make version-patch', 'make version-minor', or 'make version-major' to trigger releases"

check: fmt lint test ## Run all checks (format, lint, test)

ci: deps check build ## Run CI pipeline locally
