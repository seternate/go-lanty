.PHONY: help build run test clean vet lint docker-build docker-up docker-down db-up db-down swagger test-e2e test-e2e-all

# Variables
BINARY_NAME=lantyd
MAIN_PATH=./cmd/lantyd
DOCKER_COMPOSE=docker-compose
GO=go
ARGS?=--db "postgres://lanty:lanty@localhost:5432/lanty?sslmode=disable" --loglevel "debug"
APP_VERSION?=dev-build

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) -ldflags "-X main.AppVersion=$(APP_VERSION)" $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run:
	@echo "Running $(BINARY_NAME)..."
	$(GO) run $(MAIN_PATH) $(ARGS)

test:
	@echo "Running tests..."
	$(GO) test ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	$(GO) tool cover -func=coverage.out

test-e2e:
	@echo "Running e2e tests with venom..."
	@echo "Make sure backend is running: make docker-up-build"
	@$(DOCKER_COMPOSE) run --rm test run tests/games/*.venom.yml tests/health.venom.yml

test-e2e-all:
	@echo "Running all e2e tests with venom..."
	@echo "Make sure backend is running: make docker-up-build"
	@$(DOCKER_COMPOSE) run --rm test run tests/

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH)/main.go --outputTypes go,yaml ; \
		echo "Swagger documentation generated in ./docs"; \
	else \
		echo "swag not found. Install it with: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

## docker-build: Build the Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t go-lanty:latest .
	@echo "Docker build complete"

## docker-up: Start the full docker-compose stack
docker-up:
	@echo "Starting docker-compose stack..."
	$(DOCKER_COMPOSE) up -d
	@echo "Docker-compose stack started"

## docker-up-build: Start the full docker-compose stack and build images
docker-up-build:
	@echo "Building and starting docker-compose stack..."
	$(DOCKER_COMPOSE) up -d --build
	@echo "Docker-compose stack built and started"

## docker-down: Stop the docker-compose stack
docker-down:
	@echo "Stopping docker-compose stack..."
	$(DOCKER_COMPOSE) down
	@echo "Docker-compose stack stopped"

## db-up: Start database with migrations
db-up:
	@echo "Starting database with migrations..."
	$(DOCKER_COMPOSE) up -d database
	@echo "Waiting for database to be healthy..."
	@$(DOCKER_COMPOSE) up database-migration
	@echo "Database and migrations are ready"

## db-down: Stop database
db-down:
	@echo "Stopping database..."
	$(DOCKER_COMPOSE) stop database
	@echo "Database stopped"

## db-reset: Reset database (stop, remove volumes, start with migrations)
db-reset:
	@echo "Resetting database..."
	$(DOCKER_COMPOSE) down -v database database-migration
	$(DOCKER_COMPOSE) up -d database
	@echo "Waiting for database to be healthy..."
	@$(DOCKER_COMPOSE) up database-migration
	@echo "Database reset complete"

## dev: Start database and run the application locally
dev: db-up
	@echo "Starting development environment..."
	@echo "Database is ready. Starting application..."
	$(MAKE) run


