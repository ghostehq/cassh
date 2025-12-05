# cassh-server Docker image
# Multi-stage build for minimal image size

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o /cassh-server ./cmd/cassh-server

# Runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' cassh
USER cassh

WORKDIR /app

# Copy binary (templates and static files are embedded at build time)
COPY --from=builder /cassh-server /app/cassh-server

# Environment variables (override via docker run -e or docker-compose)
# Required for production:
#   CASSH_SERVER_URL - Public URL (e.g., https://cassh.example.com)
#   CASSH_OIDC_CLIENT_ID - Entra app client ID
#   CASSH_OIDC_CLIENT_SECRET - Entra app client secret
#   CASSH_OIDC_TENANT - Entra tenant ID
#   CASSH_CA_PRIVATE_KEY - CA private key (paste full PEM content)
ENV CASSH_LISTEN_ADDR=:8080

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/cassh-server"]
