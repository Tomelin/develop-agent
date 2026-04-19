.PHONY: help dev build test test-integration lint migrate docs clean

# Variables
GO := go
GOLANGCI_LINT := golangci-lint
DOCKER_COMPOSE := docker-compose

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -h "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/\\$$//' | sed -e 's/##//'

## dev: Start the development environment using Docker Compose
dev:
	$(DOCKER_COMPOSE) up -d

## build: Compile the binary for production
build:
	cd src/backend && $(GO) build -o bin/api cmd/api/main.go

## test: Run unit tests with coverage
test:
	cd src/backend && $(GO) test -v -cover ./...

## test-integration: Run integration tests with Docker environment
test-integration:
	cd src/backend && $(GO) test -tags=integration -v ./...

## lint: Run golangci-lint
lint:
	cd src/backend && $(GOLANGCI_LINT) run

## migrate: Run database migrations
migrate:
	@echo "Migration target not yet implemented"

## docs: Generate Swagger documentation
docs:
	@echo "Docs generation not yet implemented"

## clean: Remove build artifacts
clean:
	rm -rf src/backend/bin
	$(DOCKER_COMPOSE) down -v
