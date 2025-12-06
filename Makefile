.PHONY: help docker-up docker-down docker-logs docker-restart migrate-up migrate-down migrate-status run-local build test

# Variables
BINARY_NAME=tranquara_api
MIGRATIONS_PATH=./migrations
COMPOSE_FILE=docker-compose.yml

# Docker database connection (inside Docker network)
DB_URL_DOCKER=postgres://postgres:Nhatdien123@db:5432/tranquara_core?sslmode=disable

# Local database connection (for development without Docker)
DB_URL_LOCAL=postgresql://postgres:Nhatdien123@localhost:5432/tranquara_core?sslmode=disable

help: ## Show this help message
	@echo '🚀 Tranquara Core Service - Makefile Commands'
	@echo ''
	@echo 'Docker Commands (Primary Development):'
	@awk 'BEGIN {FS = ":.*?## "} /^docker-[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Migration Commands:'
	@awk 'BEGIN {FS = ":.*?## "} /^migrate-[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Local Development Commands:'
	@awk 'BEGIN {FS = ":.*?## "} /^(run-local|build|test):.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================
# Docker Commands (Primary Development)
# ============================================

docker-up: ## Start all services (core_service, db, keycloak, migrate)
	@echo "🚀 Starting all services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 Core Service: http://localhost:4000"
	@echo "🔐 Keycloak: http://localhost:4200 (admin/admin)"
	@echo "💾 Adminer: http://localhost:8080"

docker-down: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker-compose down

docker-logs: ## Show logs from all services
	docker-compose logs -f

docker-logs-core: ## Show logs from core service only
	docker-compose logs -f core_service

docker-logs-db: ## Show logs from database only
	docker-compose logs -f db

docker-restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	docker-compose restart

docker-rebuild: ## Rebuild and restart services
	@echo "🔨 Rebuilding services..."
	docker-compose up -d --build

docker-clean: ## Stop services and remove volumes (⚠️  deletes all data)
	@echo "⚠️  WARNING: This will delete all database data!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds..."
	@sleep 5
	docker-compose down -v

# ============================================
# Migration Commands
# ============================================

migrate-up: ## Run all database migrations (via Docker migrate service)
	@echo "🔄 Running migrations via Docker..."
	docker-compose up migrate
	@echo "✅ Migrations complete!"

migrate-status: ## Check current migration version
	@echo "📊 Checking migration status..."
	docker-compose run --rm migrate -path /migrations -database "$(DB_URL_DOCKER)" version

migrate-down: ## Rollback the last migration (⚠️  use with caution)
	@echo "⚠️  Rolling back last migration..."
	docker-compose run --rm migrate -path /migrations -database "$(DB_URL_DOCKER)" down 1

migrate-force: ## Force migration to specific version (usage: make migrate-force VERSION=1)
	@echo "⚠️  Forcing migration to version $(VERSION)..."
	docker-compose run --rm migrate -path /migrations -database "$(DB_URL_DOCKER)" force $(VERSION)

migrate-create: ## Create new migration file (usage: make migrate-create NAME=add_column_to_users)
	@echo "📝 Creating new migration: $(NAME)..."
	docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(NAME)

# ============================================
# Local Development Commands (Without Docker)
# ============================================

run-local: ## Run the application locally (requires local Go installation)
	@echo "🏃 Running application locally..."
	go run ./cmd/api

build: ## Build the application binary
	@echo "🔨 Building application..."
	go build -o bin/$(BINARY_NAME) ./cmd/api

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

migrate-local-up: ## Run migrations on local database (requires golang-migrate installed)
	@echo "🔄 Running migrations on local database..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" up

migrate-local-down: ## Rollback migrations on local database
	@echo "⚠️  Rolling back local migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" down 1

# ============================================
# Utility Commands
# ============================================

db-shell: ## Open PostgreSQL shell in Docker
	@echo "💾 Opening database shell..."
	docker-compose exec db psql -U postgres -d tranquara_core

db-reset: ## Reset database (drop and recreate - ⚠️  deletes all data)
	@echo "⚠️  WARNING: This will delete all database data!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds..."
	@sleep 5
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS tranquara_core;"
	docker-compose exec db psql -U postgres -c "CREATE DATABASE tranquara_core;"
	@echo "✅ Database reset complete. Run 'make migrate-up' to apply migrations."

.DEFAULT_GOAL := help
