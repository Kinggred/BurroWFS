FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env.local ./
# Ensure static build for Alpine compatibility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate-tool ./cmd/migrate_up
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-tool ./cmd/app

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=builder /app/migrate-tool .
COPY --from=builder /app/app-tool .
COPY --from=builder /app/.env.local .

# Debug: list files and permissions
RUN ls -l /app

ENTRYPOINT ["/app/app-tool"]
