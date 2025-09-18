# DataCrunch SDK for Go - Makefile

# Variables
GO_VERSION := 1.24
BINARY_NAME := datacrunch-sdk-example
BUILD_DIR := ./build
EXAMPLES_DIR := ./examples
COVERAGE_DIR := ./coverage

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Git related variables
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# LDFLAGS for build info
LDFLAGS := -X main.GitCommit=$(GIT_COMMIT) -X main.GitBranch=$(GIT_BRANCH) -X main.BuildTime=$(BUILD_TIME)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help
help: ## Display this help screen
	@echo "DataCrunch Go SDK - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

##@ Code Generation

generate-service: ## Generate service with all files (usage: make generate-service SERVICE=myservice [CLASS=MyService] [NAME="My Service"] [DRY_RUN=true])
	@if [ -z "$(SERVICE)" ]; then \
		echo "$(RED)❌ SERVICE parameter is required$(NC)"; \
		echo "$(YELLOW)Usage: make generate-service SERVICE=myservice$(NC)"; \
		echo "$(YELLOW)Optional: CLASS=MyService NAME=\"My Service\" DRY_RUN=true$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Generating service: $(SERVICE)...$(NC)"
	@ARGS="-service $(SERVICE)"; \
	if [ -n "$(CLASS)" ]; then \
		ARGS="$$ARGS -class $(CLASS)"; \
	fi; \
	if [ -n "$(NAME)" ]; then \
		ARGS="$$ARGS -name \"$(NAME)\""; \
	fi; \
	if [ "$(DRY_RUN)" = "true" ]; then \
		echo "$(YELLOW)🏃 Dry run mode enabled$(NC)"; \
		ARGS="$$ARGS -dry-run"; \
	fi; \
	go run tools/cmd/svc_codegen/main.go $$ARGS
	@if [ "$(DRY_RUN)" != "true" ]; then \
		echo "$(GREEN)✅ Service $(SERVICE) generated!$(NC)"; \
	fi

generate: generate-service ## Alias for generate-service

##@ Build & Test

.PHONY: all
all: clean deps lint test build ## Run all checks and build

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify

.PHONY: test
test: deps ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

.PHONY: integration-test
integration-test: deps ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v -tags="integration" ./service/...

.PHONY: test-unit
test-unit: deps ## Run unit tests only (fast, no external dependencies)
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test -v -tags="unit" ./...

.PHONY: test-coverage
test-coverage: deps ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: test-race
test-race: deps ## Run tests with race detection
	@echo "$(BLUE)Running tests with race detection...$(NC)"
	@go test -v -race ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		GOVCS=off golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		GOVCS=off golangci-lint run; \
	fi

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(BLUE)Running golangci-lint with auto-fix...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		GOVCS=off golangci-lint run --fix; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		GOVCS=off golangci-lint run --fix; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .

.PHONY: vet
vet: deps ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

.PHONY: check
check: fmt vet lint ## Run all code quality checks
	@echo "$(GREEN)All checks passed!$(NC)"

.PHONY: static-analysis
static-analysis: ## Run comprehensive static analysis (golangci-lint with all checks)
	@echo "$(BLUE)Running comprehensive static analysis...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		GOVCS=off GOFLAGS="-buildvcs=false" golangci-lint run --timeout=10m; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		GOVCS=off GOFLAGS="-buildvcs=false" golangci-lint run --timeout=10m; \
	fi

.PHONY: pre-commit
pre-commit: deps ## Run essential pre-commit checks (format, basic lint, test, build)
	@echo "$(BLUE)Running pre-commit checks...$(NC)"
	@echo "$(BLUE)Step 1/4: Dependencies and modules...$(NC)"
	@go mod tidy
	@go mod verify
	@echo "$(BLUE)Step 2/4: Code formatting...$(NC)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi
	@echo "$(BLUE)Step 3/4: Running tests with race detection...$(NC)"
	@go test -race -v ./...
	@echo "$(BLUE)Step 4/4: Build verification...$(NC)"
	@GOFLAGS="-buildvcs=false" go build ./...
	@GOFLAGS="-buildvcs=false" go build ./examples/basic/
	@echo "$(GREEN)✅ All pre-commit checks passed! Ready to commit.$(NC)"
	@echo "$(YELLOW)💡 Run 'make static-analysis' for comprehensive linting when needed.$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean -cache
	@go clean -testcache
	@go clean -modcache

.PHONY: install-tools
install-tools: ## Install required development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

.PHONY: install-hooks
install-hooks: ## Install git pre-commit hooks (lightweight version)
	@echo "$(BLUE)Installing git pre-commit hooks...$(NC)"
	@mkdir -p .git/hooks
	@echo '#!/bin/bash' > .git/hooks/pre-commit
	@echo 'set -e' >> .git/hooks/pre-commit
	@echo 'echo "Running basic pre-commit formatting..."' >> .git/hooks/pre-commit
	@echo 'go fmt ./...' >> .git/hooks/pre-commit
	@echo 'if command -v goimports >/dev/null 2>&1; then' >> .git/hooks/pre-commit
	@echo '    goimports -w .' >> .git/hooks/pre-commit
	@echo 'fi' >> .git/hooks/pre-commit
	@echo 'echo "✅ Basic formatting completed. Run '\''make pre-commit'\'' for full checks."' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Git pre-commit hook installed!$(NC)"
	@echo "$(YELLOW)The hook only runs basic formatting. Use 'make pre-commit' manually for full checks.$(NC)"

.PHONY: uninstall-hooks
uninstall-hooks: ## Remove git pre-commit hooks
	@echo "$(BLUE)Removing git pre-commit hooks...$(NC)"
	@rm -f .git/hooks/pre-commit
	@echo "$(GREEN)Git pre-commit hook removed!$(NC)"

.PHONY: mod-update
mod-update: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy

.PHONY: security
security: ## Run security checks
	@echo "$(BLUE)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)gosec not found. Installing...$(NC)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

.PHONY: version
version: ## Display version information
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Git branch: $(GIT_BRANCH)"
	@echo "Build time: $(BUILD_TIME)"

.PHONY: example
example: build-example ## Build and run example
	@echo "$(BLUE)Running example...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

# Development workflow targets
.PHONY: setup
setup: install-tools ## Setup development environment (tools only)
	@echo "$(GREEN)Development environment setup complete!$(NC)"
	@echo "$(YELLOW)💡 Optional: Run 'make install-hooks' for lightweight git formatting hooks$(NC)"
	@echo "$(YELLOW)💡 Use 'make pre-commit' manually before important commits$(NC)"

.PHONY: dev
dev: clean deps fmt vet test build ## Complete development workflow (without heavy static analysis)

.PHONY: dev-full  
dev-full: clean deps static-analysis test build ## Complete development workflow with full static analysis

.PHONY: ci
ci: deps fmt vet test-coverage build ## CI/CD workflow

# Default target
.DEFAULT_GOAL := help