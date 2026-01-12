.PHONY: help run test test-coverage test-e2e swagger clean db-up db-down db-reset

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

test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.out

test-e2e: ## Run e2e tests with venom
	@echo "Running e2e tests with venom..."
	@echo "Make sure backend is running"
	@if ! command -v venom >/dev/null 2>&1; then \
		echo "venom not found. Install it with: go install github.com/ovh/venom/cmd/venom@latest"; \
		exit 1; \
	fi
	@mkdir -p test-results
	@venom run --output-dir test-results $$(find tests -name "*.venom.yml")

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
	rm -rf test-results/
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
