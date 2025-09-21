# Go-Solid Makefile
# Provides convenient commands for running the application

.PHONY: help server cron migration build clean test docker-up docker-down

# Default target
help:
	@echo "Go-Solid Application Commands:"
	@echo ""
	@echo "Development Commands:"
	@echo "  make server     - Run the HTTP server"
	@echo "  make cron       - Run the cron service"
	@echo "  make migration  - Run database migrations"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build      - Build both server and cron binaries"
	@echo "  make clean      - Clean build artifacts"
	@echo ""
	@echo "Database Commands:"
	@echo "  make docker-up  - Start MySQL and Redis using Docker"
	@echo "  make docker-down- Stop MySQL and Redis containers"
	@echo ""
	@echo "Testing Commands:"
	@echo "  make test       - Run tests"
	@echo ""

# Run the HTTP server
server:
	@echo "Starting HTTP server..."
	go run cmd/server/main.go

# Run the cron service
cron:
	@echo "Starting cron service..."
	go run cmd/cron/main.go

# Run database migrations
migration:
	@echo "Running database migrations..."
	go run cmd/migration/main.go

# Build all binaries
build:
	@echo "Building server binary..."
	go build -o bin/server cmd/server/main.go
	@echo "Building cron binary..."
	go build -o bin/cron cmd/cron/main.go
	@echo "Building migration binary..."
	go build -o bin/migration cmd/migration/main.go
	@echo "Build complete! Binaries available in bin/ directory"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Start Docker services
docker-up:
	@echo "Starting MySQL and Redis containers..."
	docker-compose up -d
	@echo "Database services started!"
	@echo "  MySQL: localhost:3306"
	@echo "  Redis: localhost:6379"

# Stop Docker services
docker-down:
	@echo "Stopping MySQL and Redis containers..."
	docker-compose down
	@echo "Database services stopped!"

# Development setup (build and run server)
dev: build
	@echo "Starting development environment..."
	./bin/server

# Production setup (build optimized binaries)
prod-build:
	@echo "Building optimized binaries for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/cron cmd/cron/main.go
	@echo "Production binaries built!"
