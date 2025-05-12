FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o socket-app .

# Create a lightweight production image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/socket-app .

# Copy any necessary resources (including PDF files)
COPY resources/ ./resources/

# Expose UDP port (assuming the application uses UDP)
# Note: This is for documentation purposes. You can expose additional ports at runtime.
EXPOSE 8080/udp

# Set the entrypoint
ENTRYPOINT ["/app/socket-app"]

# Default command (can be overridden)
CMD ["serve"]
