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

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/battlestation .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./battlestation"]