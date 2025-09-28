# Задачник на Golang
Архитектура: Domain -> Repository -> Service -> Http-handler

## Быстрый старт
1. Установить Golang 1.22+ и MySQL 8+
2. Создайте БД и таблицы, выполнив `sql/init.sql`
3. Создайте `.env` и отредактируйте доступ к БД.
4. Запустите сервер:
```bash
go run ./cmd/server
```
5. Проверьте: `GET http://localhost:8080/health -> 200 OK`

## Пример .env файла:
```.env
APP_NAME=ToDoApp
APP_ENV=dev
APP_PORT=8080

MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DB=todolist
MYSQL_PARAMS=charset=utf8mb4&parseTime=true&loc=Local

JWT_SECRET=secret_jwt
JWT_TTL_MINUTES=120
COOKIE_DOMAIN=localhost
```

## Структура проекта
* internal/models - модели
* internal/repository - интерфейс и реализации для работы с БД
* internal/service - бизнес-логика
* internal/transport/http - http-хендлеры и маршруты