FROM golang:1.25 AS builder

WORKDIR /app

# сначала зависимости
COPY go.mod go.sum ./
RUN go mod download

# ставим goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# теперь исходники
COPY . .

# собираем бинарник
RUN go build -o subscription main.go

# финальный образ (чистый)
FROM debian:bookworm-slim

WORKDIR /app

# утилиты для миграций
RUN apt-get update && apt-get install -y postgresql-client ca-certificates && rm -rf /var/lib/apt/lists/*

# копируем бинарник и миграции
COPY --from=builder /app/subscription .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY entrypoint.sh .

RUN chmod +x /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]


