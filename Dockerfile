# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache ca-certificates=20250911-r0 git=2.49.1-r0 tzdata=2025b-r0

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev} -X main.commit=${COMMIT:-unknown} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o shelly-exporter \
    ./cmd/shelly-exporter

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates=20241121-r1 tzdata=2025b-r0

# Create non-root user
RUN addgroup -g 1001 -S shelly && \
    adduser -u 1001 -S shelly -G shelly

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/shelly-exporter .

# Copy configuration template
COPY examples/config.yaml /etc/shelly-exporter/config.yaml

# Change ownership to non-root user
RUN chown -R shelly:shelly /app /etc/shelly-exporter

# Switch to non-root user
USER shelly

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["./shelly-exporter"]
CMD ["--config=/etc/shelly-exporter/config.yaml"]
