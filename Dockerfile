FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache libc6-compat
COPY --from=builder /app/main .

CMD ["/main"]