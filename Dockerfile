# Use Go 1.24.2 for building
FROM golang:1.24.2 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o news-api ./cmd/main

# Use a lightweight base image for production
FROM alpine:latest

# Set environment variable for release mode
ENV GIN_MODE=release

# Working directory inside the container
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/news-api .

# Copy the .env file if needed at runtime
COPY .env .

# Expose the port your app uses
EXPOSE 8080

# Run the compiled binary
CMD ["./news-api"]
