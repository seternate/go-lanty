.PHONY: help run test-unit test-integration test-e2e test-all test-e2e-coverage swagger clean db-up db-down db-reset

MAIN_PATH=./cmd/lantyd

.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	@echo "Running..."
	go run $(MAIN_PATH) --db "postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable" --loglevel "debug"

test-unit: ## Run unit tests only (excludes integration tests, no cache)
	@echo "Running unit tests..."
	go test -count=1 ./... -tags='!integration'

test-integration: ## Run integration tests only (requires DB, no cache)
	@echo "Running integration tests..."
	@echo "Make sure database is running (make db-up)"
	@INTEGRATION_PACKAGES=$$(find . -name '*_integration_test.go' -exec dirname {} \; | sort -u | sed 's|^\./|./|' | sed 's|^\.$$|./|' | tr '\n' ' '); \
	if [ -z "$$INTEGRATION_PACKAGES" ]; then \
		echo "No integration tests found"; \
		exit 0; \
	fi; \
	go test -count=1 -tags=integration $$INTEGRATION_PACKAGES

test-e2e: ## Run e2e tests with venom
	@echo "Running e2e tests with venom..."
	@echo "Make sure backend is running"
	@if ! command -v venom >/dev/null 2>&1; then \
		echo "venom not found. Install it with: go install github.com/ovh/venom/cmd/venom@latest"; \
		exit 1; \
	fi
	@venom run $$(find tests -name "*.venom.yml"); \
	rm -f venom.log

test-e2e-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out

test-all: ## Run unit, integration, and e2e tests sequentially (no cache for unit/integration)
	@echo "Running all tests..."
	$(MAKE) test-unit && $(MAKE) test-integration && $(MAKE) test-e2e

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH)/main.go --outputTypes go,yaml ; \
		echo "Swagger documentation generated in ./docs"; \
	else \
		echo "swag not found. Install it with: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi

clean: ## Remove artifacts
	rm -f coverage.out coverage.html
	@echo "Clean complete"

db-up: ## Start database with migrations
	@echo "Starting database with migrations..."
	docker compose up -d database
	@echo "Waiting for database to be healthy..."
	@docker compose up database-migration
	@echo "Database and migrations are ready"

db-down: ## Stop database
	@echo "Stopping database..."
	docker compose stop database
	@echo "Database stopped"

db-reset: ## Reset database (stop, remove volumes, start with migrations)
	@echo "Resetting database..."
	docker compose down -v database database-migration
	docker compose up -d database
	@echo "Waiting for database to be healthy..."
	@docker compose up database-migration
	@echo "Database reset complete"
