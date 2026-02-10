COMPOSE := docker compose
APP_NAME := pulselounge

.PHONY: help up up-build down reset logs ps restart ui-install ui-build api-run api-build run clean

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start stack in background
	$(COMPOSE) up -d

up-build: ## Build and start stack in background
	$(COMPOSE) up --build -d

down: ## Stop stack
	$(COMPOSE) down

reset: ## Stop stack and remove volumes (fresh Postgres)
	$(COMPOSE) down -v

restart: ## Restart app and postgres containers
	$(COMPOSE) restart app postgres

logs: ## Tail app and postgres logs
	$(COMPOSE) logs -f app postgres

ps: ## Show compose services status
	$(COMPOSE) ps

ui-install: ## Install frontend dependencies
	cd frontend && npm install

ui-build: ## Build frontend bundle embedded by Go
	cd frontend && npm run build

api-build: ## Build Go server binary
	go build -o $(APP_NAME) ./cmd/server

api-run: ## Run Go server (expects UI already built)
	set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/server

run: ui-build api-run ## Build UI and run API locally

clean: ## Remove local build artifacts
	rm -f ./$(APP_NAME) ./server
