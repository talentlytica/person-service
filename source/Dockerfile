# Multi-stage build

# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files from app directory
COPY app/go.mod app/go.sum ./app/
COPY app/vendor ./app/vendor

# Download dependencies (skip if using vendor)
WORKDIR /build/app
# RUN go mod download - not needed when using vendor

# Copy source code
WORKDIR /build
COPY app/ ./app/

# Ensure go.sum is up to date and build the binary
WORKDIR /build/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -cover -o /build/person-service .

# Stage 2: Runtime
FROM alpine:3.18

# Copy ca-certificates from builder (avoids needing apk add which has network issues)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Create coverage directory
RUN mkdir -p /app/coverage && chown -R appuser:appuser /app/coverage

# Copy binary from builder
COPY --from=builder /build/person-service .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the application
CMD ["./person-service"]
