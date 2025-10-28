# Dockerfile
FROM golang:1.24 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Build specific package from main (here: ./cmd/app)
RUN go build -ldflags="-s -w" -o /app/main ./cmd/app

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/main /app/main

USER nonroot
ENV PORT=8080
EXPOSE 8080

CMD ["/app/main"]
