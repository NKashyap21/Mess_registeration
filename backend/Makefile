.PHONY: build run test clean docker-up docker-down dev migrate

# Variables
APP_NAME := mess-registration
DOCKER_COMPOSE := docker-compose
GO_FILES := $(shell find . -name "*.go" -type f)

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) .

# Run the application locally (requires local PostgreSQL)
run:
	@echo "Running $(APP_NAME)..."
	@go run .

# Run with hot reload using air
dev:
	@echo "Starting development server with hot reload..."
	@air

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Docker commands
docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

docker-logs:
	@echo "Showing Docker logs..."
	@$(DOCKER_COMPOSE) logs -f

docker-build:
	@echo "Building Docker image..."
	@$(DOCKER_COMPOSE) build

# Database commands
db-reset:
	@echo "Resetting database..."
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d db
	@sleep 5
	@$(DOCKER_COMPOSE) up -d app

# Development workflow
setup: docker-up
	@echo "Waiting for database to be ready..."
	@sleep 10
	@echo "Setup complete! Application should be running on http://localhost:8080"

# Quick development start
quick-dev: docker-down docker-up
	@echo "Quick development setup complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Generate documentation
docs:
	@echo "The API documentation is available in APIdocs.md"

# Show help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  dev         - Run with hot reload (requires air)"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"
	@echo "  docker-logs - Show Docker logs"
	@echo "  db-reset    - Reset database"
	@echo "  setup       - Initial setup with Docker"
	@echo "  quick-dev   - Quick development setup"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  help        - Show this help"
