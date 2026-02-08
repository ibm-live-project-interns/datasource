# --- Stage 1: Builder ---
FROM golang:1.22-alpine AS builder

# ⚠️ CRITICAL: Install C compiler (gcc, musl-dev) for Kafka library
RUN apk add --no-cache git build-base

WORKDIR /app

# Download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -tags musl: Required for Alpine Linux compatibility with CGO
# CGO_ENABLED=1: Required for confluent-kafka-go
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o datasource_bin main.go

# --- Stage 2: Runner ---
FROM alpine:latest

WORKDIR /root/

# Install CA certs for network calls
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/datasource_bin .

# Command to run
CMD ["./datasource_bin"]