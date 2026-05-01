.PHONY: build run clean help

# Variáveis
APP_NAME=api
BIN_DIR=bin
MAIN_PATH=./cmd/

# Go build flags
GO_BUILD_FLAGS=-ldflags="-s -w"

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

deps: ## Download dependencies
	go mod download
	go mod tidy

build: deps ## Build the application
	@mkdir -p $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application
	go run $(MAIN_PATH)

clean: ## Clean build artifacts
	rm -rf $(BIN_DIR)

lint: ## Run linter
	@hash golangci-lint 2>/dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest

dev: ## Run with hot reload (requires air)
	@hash air 2>/dev/null || go install github.com/cosmtrek/air@latest
	air -c .air.toml

.DEFAULT_GOAL := help
