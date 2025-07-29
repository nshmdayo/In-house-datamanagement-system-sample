# Variables
APP_NAME := datamanagement-system
BINARY_NAME := main
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOLINT := golangci-lint

# Build the binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/server

# Build for production
.PHONY: build-prod
build-prod:
	@echo "Building $(APP_NAME) for production..."
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_NAME) ./cmd/server

# Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(BINARY_NAME)

# Run in development mode with hot reload
.PHONY: dev
dev:
	@echo "Starting development server..."
	air

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GOLINT) run

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: update-deps
update-deps:
	@echo "Updating dependencies..."
	$(GOMOD) get -u all
	$(GOMOD) tidy

# Docker commands
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

# Database commands
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	$(GOCMD) run ./cmd/migrate

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	$(GOCMD) run ./cmd/seed

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	docker-compose exec postgres psql -U postgres -d datamanagement_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Setup development environment
.PHONY: setup-dev
setup-dev:
	@echo "Setting up development environment..."
	cp .env.example .env
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Development environment setup complete!"
	@echo "Please edit .env file with your configuration"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Generate API documentation
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	swag init -g ./cmd/server/main.go -o ./docs

# Generate mocks
.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	mockgen -source=./internal/services/user.go -destination=./mocks/user_service_mock.go
	mockgen -source=./internal/services/audit.go -destination=./mocks/audit_service_mock.go

# Backup database
.PHONY: backup
backup:
	@echo "Creating database backup..."
	docker-compose exec postgres pg_dump -U postgres datamanagement_db > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Restore database
.PHONY: restore
restore:
	@echo "Restoring database..."
	@read -p "Enter backup file path: " backup_file; \
	docker-compose exec -T postgres psql -U postgres datamanagement_db < $$backup_file

# Health check
.PHONY: health
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Application is not running"

# View logs
.PHONY: logs
logs:
	@echo "Viewing application logs..."
	docker-compose logs -f backend

# Load test
.PHONY: load-test
load-test:
	@echo "Running load test..."
	k6 run tests/load/api_test.js

# Full CI pipeline
.PHONY: ci
ci: deps fmt lint security test build

# Full deployment pipeline
.PHONY: deploy
deploy: ci docker-build docker-compose-up

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  build-prod     - Build for production"
	@echo "  run            - Run the application"
	@echo "  dev            - Start development server with hot reload"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  bench          - Run benchmarks"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  security       - Run security scan"
	@echo "  deps           - Download dependencies"
	@echo "  update-deps    - Update dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start services with Docker Compose"
	@echo "  docker-compose-down - Stop services with Docker Compose"
	@echo "  db-migrate     - Run database migrations"
	@echo "  db-seed        - Seed database"
	@echo "  db-reset       - Reset database"
	@echo "  setup-dev      - Setup development environment"
	@echo "  install-tools  - Install development tools"
	@echo "  docs           - Generate API documentation"
	@echo "  mocks          - Generate mocks"
	@echo "  backup         - Backup database"
	@echo "  restore        - Restore database"
	@echo "  health         - Check application health"
	@echo "  logs           - View application logs"
	@echo "  load-test      - Run load test"
	@echo "  ci             - Run CI pipeline"
	@echo "  deploy         - Deploy application"
	@echo "  help           - Show this help message"
