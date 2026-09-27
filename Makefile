.PHONY: help run build test tidy lint init docker-up docker-down migrate-up migrate-down migrate-create migrate-version migrate-force seed

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

migrate-up: ## Run all pending database migrations
	@go run ./cmd/migrate up

migrate-down: ## Rollback 1 migration step
	@go run ./cmd/migrate down

migrate-version: ## Check current database migration version
	@go run ./cmd/migrate version

migrate-create: ## Create new up/down migration files (usage: make migrate-create NAME=create_table_name)
	@if [ -z "$(NAME)" ]; then echo "❌ Please specify NAME. Example: make migrate-create NAME=create_users_table"; exit 1; fi
	@go run ./cmd/migrate create $(NAME)

migrate-force: ## Force set migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "❌ Please specify VERSION. Example: make migrate-force VERSION=1"; exit 1; fi
	@go run ./cmd/migrate force $(VERSION)

seed: ## Run database seeding scripts
	@go run ./cmd/seed

docker-up: ## Start docker containers (app & postgres)
	@docker-compose up -d

docker-down: ## Stop docker containers
	@docker-compose down

init: ## Initialize as a new project with a custom module name (usage: make init MODULE=github.com/user/project)
	@if [ -z "$(MODULE)" ]; then echo "❌ Please specify MODULE. Example: make init MODULE=github.com/my/project"; exit 1; fi
	@./scripts/init-project.sh $(MODULE)
