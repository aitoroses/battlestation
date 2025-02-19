# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o battlestation ./cmd/battlestation

# Development stage
FROM golang:1.22-alpine AS dev

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Source code will be mounted as volume

# Production stage
FROM alpine:latest AS prod

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/battlestation .

# Expose metrics port
EXPOSE 8080

# Run the application
CMD ["./battlestation"]