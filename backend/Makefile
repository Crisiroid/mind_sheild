.PHONY: build run test clean deps migrate

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server ./cmd/server

# Run the application
run:
	@echo "Starting application..."
	go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run migrations (handled by auto-migrate on startup)
migrate:
	@echo "Migrations are run automatically on startup"
	go run ./cmd/server

# Help command
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make migrate  - Run migrations and start server"
