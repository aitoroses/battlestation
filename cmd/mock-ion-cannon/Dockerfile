FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-ion-cannon ./cmd/mock-ion-cannon

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/mock-ion-cannon .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./mock-ion-cannon"]