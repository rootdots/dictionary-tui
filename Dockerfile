# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy only the files needed for go mod download
COPY go.* ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o dt

# Runtime stage
FROM alpine:latest

# Create a non-root user
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/dt /usr/local/bin/

# Set default environment variables
ENV TERM=xterm-256color

ENTRYPOINT ["dt"]
CMD ["--help"]
