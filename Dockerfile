FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY ingestor/shared/ ./ingestor/shared/

# Copy datasource files
COPY datasource/go.mod datasource/go.sum ./datasource/
WORKDIR /app/datasource
RUN go mod download

COPY datasource/ .
RUN go build -o datasource .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/datasource/datasource .
COPY --from=builder /app/datasource/.env .

CMD ["./datasource"]
