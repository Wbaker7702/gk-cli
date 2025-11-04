.PHONY: build install clean test run help

# Build variables
BINARY_NAME=gk
VERSION?=dev
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v ./...

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 -v ./...

build-darwin: ## Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 -v ./...
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 -v ./...

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe -v ./...

build-all: build-linux build-darwin build-windows ## Build for all platforms

install: build ## Install the binary to $GOPATH/bin
	$(GOCMD) install -ldflags "$(LDFLAGS)" ./...

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

test: ## Run tests
	$(GOTEST) -v ./...

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

run: build ## Build and run the CLI
	./$(BINARY_NAME)

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

lint: fmt vet ## Run linters
