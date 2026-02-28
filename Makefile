COMPOSE := docker compose
APP_NAME := pulselounge

.PHONY: help dev-up dev-down ui-install ui-build api-build clean

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev-up: ## Start hot-reload dev stack (Go + Vite + Postgres)
	$(COMPOSE) -f docker-compose.dev.yml up

dev-down: ## Stop hot-reload dev stack
	$(COMPOSE) -f docker-compose.dev.yml down

ui-install: ## Install frontend dependencies
	cd frontend && npm install

ui-build: ## Build frontend bundle embedded by Go
	cd frontend && npm run build

api-build: ## Build Go server binary
	go build -o $(APP_NAME) ./cmd/server

clean: ## Remove local build artifacts
	rm -f ./$(APP_NAME) ./server
