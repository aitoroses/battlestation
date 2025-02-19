# Build settings
BINARY_NAME=battlestation
GO=go

# Build the application
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/battlestation

# Run all tests
test:
	$(GO) test -v ./...

# Run integration tests
test-integration:
	$(GO) test -v ./tests

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

# Build docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run docker compose
docker-up:
	docker-compose up --build -d

# Stop docker compose
docker-down:
	docker-compose down

# Run test cases script
test-cases:
	./tests.sh

# Format code
fmt:
	$(GO) fmt ./...

# Run linter
lint:
	$(GO) vet ./...

.PHONY: build test test-integration run clean docker-build docker-up docker-down test-cases fmt lint