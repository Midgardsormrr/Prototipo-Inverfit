# ==============================================================================
# Inverfit API - Multi-stage Production Dockerfile
# ==============================================================================

# Stage 1: Build binary
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache ca-certificates tzdata git

# Download Go module dependencies (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Compile static binary with target OS and architecture
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
    -trimpath \
    -ldflags="-w -s" \
    -o /bin/inverfit-api ./cmd/api

# Stage 2: Minimal runtime
FROM alpine:3.20

# Install runtime dependencies: ca-certificates, tzdata, and curl for healthchecks
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

# Set timezone default
ENV TZ=America/Santiago

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/inverfit-api /app/inverfit-api

# Assign permissions to non-privileged user
RUN chown -R appuser:appgroup /app

USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["/app/inverfit-api"]
