# Server Dockerfile
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

# Enable Go dependency caching
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags)" \
    -o /hush ./cmd/server

# Migration builder
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /migrate ./cmd/migrate/main.go

# Server runtime
FROM alpine:3.19 as run-server

RUN apk add --no-cache \
    ca-certificates \
    tzdata

COPY --from=builder /hush /usr/local/bin/hush

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup -h /app && \
    chown -R appuser:appgroup /app

USER appuser
WORKDIR /app

EXPOSE 8080
ENV PORT=8080 \
    HOSTNAME="0.0.0.0" \
    GIN_MODE=release

CMD ["hush"]

# Migration runtime
FROM alpine:3.19 as run-migrate

RUN apk add --no-cache \
    postgresql-client \
    ca-certificates

COPY --from=builder /migrate /usr/local/bin/migrate
COPY --from=builder /app/db/migrations /migrations

USER nobody:nobody
WORKDIR /migrations

ENTRYPOINT ["migrate"]

