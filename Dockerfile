# Runtime stage
FROM alpine:3.18

# Install runtime dependencies and create non-root user
RUN apk add --no-cache ca-certificates=20241121-r1 tzdata=2025b-r0 && \
    addgroup -g 1001 -S shelly && \
    adduser -u 1001 -S shelly -G shelly

# Set working directory
WORKDIR /app

# Copy pre-built binary from GoReleaser
COPY shelly-exporter .

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
    CMD ["sh", "-c", "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"]

# Run the application
ENTRYPOINT ["./shelly-exporter"]
CMD ["--config=/etc/shelly-exporter/config.yaml"]
