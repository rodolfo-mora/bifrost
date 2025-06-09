# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bifrost ./cmd/bifrost

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' bifrost

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bifrost /app/bifrost

# Copy configuration files
COPY --from=builder /app/config.yaml /app/config.yaml

# Set ownership
RUN chown -R bifrost:bifrost /app

# Switch to non-root user
USER bifrost

# Expose ports
EXPOSE 4318 4317

# Run the application
ENTRYPOINT ["/app/bifrost"] 