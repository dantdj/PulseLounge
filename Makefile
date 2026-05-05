ifneq (,$(wildcard ./.env))
include .env
export
endif

COMPOSE := docker compose
APP_NAME := pulselounge
POSTGRES_USER ?= postgres
POSTGRES_DB ?= pulselounge
UI_DIST_DIR := frontend/dist
UI_EMBED_DIR := frontend/embed/generated

.PHONY: help dev-up dev-down db-reset-dev migrate-up migrate-status seed-dev seed-reset-dev ui-install ui-build api-build lint-go lint-frontend clean

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev-up: ## Start hot-reload dev stack (Go + Vite + Postgres)
	$(COMPOSE) -f docker-compose.dev.yml up -d postgres
	# migrate-up waits for Postgres to accept connections before running goose.
	$(MAKE) migrate-up
	${COMPOSE} -f docker-compose.dev.yml up -d azurite
	$(COMPOSE) -f docker-compose.dev.yml up api ui 

dev-down: ## Stop hot-reload dev stack
	$(COMPOSE) -f docker-compose.dev.yml down

db-reset-dev: ## Destroy dev Postgres data and restart with a fresh database
	$(COMPOSE) -f docker-compose.dev.yml down -v
	$(COMPOSE) -f docker-compose.dev.yml up -d postgres

migrate-up: ## Apply pending database migrations
	go run ./cmd/migrate up

migrate-status: ## Show migration status
	go run ./cmd/migrate status

seed-dev: migrate-up ## Seed local dev Postgres data (non-destructive)
	$(COMPOSE) -f docker-compose.dev.yml exec -T postgres \
		psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" < db/seed/dev.sql

seed-reset-dev: ## Reset + reseed local dev Postgres data for the current schema
	$(COMPOSE) -f docker-compose.dev.yml exec -T postgres \
		psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" < db/seed/reset.sql

ui-install: ## Install frontend dependencies
	cd frontend && npm ci

ui-build: ## Build frontend bundle and stage it for Go embed
	cd frontend && npm run build
	rm -rf $(UI_EMBED_DIR)
	mkdir -p $(UI_EMBED_DIR)
	cp -R $(UI_DIST_DIR)/. $(UI_EMBED_DIR)/

api-build: ui-build ## Build Go server binary with staged UI assets
	go build -o $(APP_NAME) ./cmd/server

lint-go: ## Run Go linters
	golangci-lint run

lint-frontend: ## Run frontend linters
	cd frontend && npm run lint

clean: ## Remove local build artifacts
	rm -rf $(UI_DIST_DIR) $(UI_EMBED_DIR)
	rm -rf ./cmd/server/web
	rm -f ./$(APP_NAME) ./server
