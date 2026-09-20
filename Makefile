.PHONY: help run build test tidy lint init docker-up docker-down

help: ## Display available make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run application locally
	@go run ./cmd/api

build: ## Build binary executable
	@mkdir -p bin
	@go build -o bin/server ./cmd/api

test: ## Run unit tests
	@go test -v ./...

tidy: ## Clean up and download dependencies
	@go mod tidy

docker-up: ## Start docker containers (app & postgres)
	@docker-compose up -d

docker-down: ## Stop docker containers
	@docker-compose down

init: ## Initialize as a new project with a custom module name (usage: make init MODULE=github.com/user/project)
	@if [ -z "$(MODULE)" ]; then echo "❌ Please specify MODULE. Example: make init MODULE=github.com/my/project"; exit 1; fi
	@./scripts/init-project.sh $(MODULE)
