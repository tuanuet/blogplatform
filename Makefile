.PHONY: all build run test clean lint swagger deps docker-up docker-down migrate-up migrate-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api
MAIN_PATH=./cmd/api

# Swagger
SWAG=swag

# Docker
DOCKER_COMPOSE=docker compose

# Default target
all: lint test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GORUN) $(MAIN_PATH)

# Run with hot reload (requires air)
dev:
	air

# Run tests
test:
	$(GOTEST) -v -cover ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Generate Swagger documentation
swagger:
	$(SWAG) init -g cmd/api/main.go -o docs

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install development tools
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Database migrations (using golang-migrate)
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/boilerplate?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/boilerplate?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with hot reload (air)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make lint           - Run linter"
	@echo "  make swagger        - Generate Swagger docs"
	@echo "  make deps           - Download dependencies"
	@echo "  make tools          - Install dev tools"
	@echo "  make docker-up      - Start Docker services"
	@echo "  make docker-down    - Stop Docker services"
	@echo "  make migrate-up     - Run migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo "  make clean          - Clean build artifacts"
