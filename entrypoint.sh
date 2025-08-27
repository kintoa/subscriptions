#!/bin/sh
set -e

DB_URL="host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=$DB_SSLMODE"

echo "🔗 Используем строку подключения: $DB_URL"

echo "🚀 Выполнение миграций..."
goose -dir /app/migrations postgres "$DB_URL" up

echo "✅ Запуск приложения..."
./subscription
