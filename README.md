# Subscriptions API

REST API сервис для управления подписками. Написан на **Go (Golang)** с использованием:
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- [GORM](https://gorm.io/) — ORM для работы с PostgreSQL
- [Goose](https://github.com/pressly/goose) — миграции
- [Swagger](https://github.com/swaggo/gin-swagger) — автодокументация API
- Docker + docker-compose — контейнеризация

---

## Запуск проекта

### 1. Клонировать репозиторий
```bash
git clone https://github.com/kintoa/subscriptions.git
cd subscriptions
```

### 2. Поднять сервис
Создайте и заполните файл .env

Выполните команду:
```bash
docker-compose up --build
```

API будет доступен по адресу:
```bash
http://localhost:8081
```

Документация Swagger:
```bash
http://localhost:8081/swagger/index.html
```


### Структура проекта
``` bash
.
├── db/                # Инициализация БД
├── migrations/        # Goose миграции
├── dto/               # DTO для сериализации/десериализации
├── middleware/        # Middleware 
├── models/            # GORM-модели
├── routes/            # Роуты
├── main.go            # Точка входа
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Методы API

- `GET /subscriptions` — получить список подписок  
- `GET /subscriptions/{id}` — получить подписку по ID  
- `POST /subscriptions` — создать новую подписку  
- `PATCH /subscriptions/{id}` — обновить подписку (частично)  
- `DELETE /subscriptions/{id}` — удалить подписку  
- `GET /subscriptions/summary` — получить суммарную стоимость подписок с фильтрацией  

## Примеры запросов и ответов

### Request
```bash
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Response
```json
{
  "id": 1,
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null
}
```

## Логирование

Для логов написан logging middleware, который логирует запросы и ответы методов

### Пример:
```json
{
  "latency": "26.2647ms",
  "level": "info",
  "method": "POST",
  "msg": "http request",
  "path": "/subscriptions",
  "request": {
    "price": 300,
    "service_name": "yandex",
    "start_date": "06-2025",
    "user_id": "1"
  },
  "response": {
    "id": 12,
    "price": 300,
    "service_name": "yandex",
    "start_date": "06-2025",
    "user_id": "1"
  },
  "status": 201,
  "time": "2025-08-27 23:45:01"
}
```
