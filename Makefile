.PHONY: help build test lint fmt clean install-tools install-hooks

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the project
	@echo "Building..."
	go build -v ./...

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w -local github.com/platonoff-dev/corekv .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -f coverage.out coverage.html

install-tools: ## Install development tools
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

install-hooks: ## Install pre-commit hooks
	@echo "Installing pre-commit hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
		echo "Pre-commit hooks installed successfully"; \
	else \
		echo "pre-commit is not installed. Install it with: pip install pre-commit"; \
		exit 1; \
	fi

mod-tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	go mod tidy

mod-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	go mod verify

all: fmt lint test build ## Run all checks and build
