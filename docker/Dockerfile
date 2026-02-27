# ========= Build stage =========
FROM golang:1.22-alpine AS builder
WORKDIR /gogogo-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate/main.go

# ========= Runtime stage =========
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /gogogo-api/app .
COPY --from=builder /gogogo-api/migrate .
COPY --from=builder /gogogo-api/internal/orm/postgres/migrations ./internal/orm/postgres/migrations

EXPOSE 8080
CMD ["./app"]