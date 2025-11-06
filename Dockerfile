# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy everything needed for building
COPY . .

# Download dependencies and build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o dt

# Runtime stage
FROM alpine:latest

# Create a non-root user
RUN adduser -D -h /home/appuser appuser

# Copy the binary and documentation
COPY --from=builder /app/dt /usr/local/bin/
COPY --from=builder /app/README.md /usr/share/doc/dt/
COPY --from=builder /app/LICENSE /usr/share/doc/dt/

# Set ownership and permissions
RUN chown -R appuser:appuser /usr/share/doc/dt && \
    chmod 755 /usr/local/bin/dt

# Switch to non-root user
USER appuser
WORKDIR /home/appuser

# Set default environment variables
ENV TERM=xterm-256color

ENTRYPOINT ["dt"]
CMD ["--help"]
