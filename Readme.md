# Distributed Calculator

Распределённая система для асинхронного вычисления арифметических выражений в многопользовательском режиме с персистентностью в SQLite и взаимодействием по gRPC.

---

## Описание проекта

Система состоит из двух компонентов:

1. **Orchestrator (Server)**  
   - REST API для пользователей (`/api/v1/...`)  
   - gRPC-клиент для агентов  
   - Хранит пользователей и результаты выражений в SQLite  
   - Маршрутит запросы, разбивая выражение на подзадачи (по один RPC-запрос на каждую операцию)

2. **Agent**  
   - gRPC-сервер, получает задачи от Orchestrator, считает их локально (моно-поточными goroutine)  
   - Возвращает результат через тот же поток  
   - Может масштабироваться добавлением новых контейнеров/инстансов  

**Основные возможности**  
- Регистрация и логин пользователей (JWT)  
- Запрос на вычисление выражения: клиент получает результат, когда сервер завершил все подзадачи  
- Сохранение истории вычислений в БД  
- REST-ENDPOINTS:  
  - `POST /api/v1/register`  
  - `POST /api/v1/login`  
  - `POST /api/v1/calculate`  
  - `GET  /api/v1/expressions`  
  - `GET  /api/v1/expressions/{id}`  
- Взаимодействие Orchestrator ↔ Agent по gRPC  

---

## Архитектура

+-----------+ HTTP +-----------------+ gRPC +---------+
| Клиент | <----------> | Orchestrator | <---------> | Agent |
| (curl, | REST | (REST + gRPC | Tasks | (Compute|
| Postman) | | client) | | Worker)|
+-----------+ +-----------------+ +---------+
| |
| +-------------------------------+ |
+-----> | SQLite (expressions) | <------------+
+-------------------------------+

- Клиент отправляет JWT-авторизованный REST-запрос на Orchestrator.  
- Orchestrator транслирует полное выражение в поток gRPC-запросов агенту.  
- Агент считает выражение (каждую операцию) и потоково возвращает результаты.  
- Orchestrator агрегирует ответы, сохраняет в БД и возвращает пользователю.

---

## Что нужно

- Go 1.23 и выше  
- Docker & Docker-Compose (для контейнерного запуска)  
- `protoc` + плагины `protoc-gen-go` и `protoc-gen-go-grpc` (для генерации proto, если меняете `.proto`)  

---

## Клонирование репозитория

```bash
git clone https://github.com/AlexDillz/distributed-calculator.git
cd distributed-calculator

Локальный запуск без Docker
Сборка proto (если меняли .proto):

make proto

Запуск Orchestrator:

DATABASE_PATH=./calc.db GRPC_PORT=:50051 HTTP_PORT=:8080 go run cmd/server/main.go

Запуск Agent (в другом терминале):

GRPC_PORT=:50051 COMPUTING_POWER=4 go run cmd/agent/main.go
COMPUTING_POWER — число параллельных воркеров (горутин).

Запуск через Docker & Docker-Compose
Сборка образа:

make build

Поднять сервисы:

make up

Или напрямую:

docker-compose up --build -d

Просмотр логов:

make logs

Остановка и очистка:

make down

Environment Variables

Переменная	      Описание	                            По умолчанию

DATABASE_PATH	   Путь к файлу SQLite DB	             calc.db
GRPC_PORT	      Порт для gRPC сервера агента	       :50051
HTTP_PORT	      Порт для HTTP сервера Orchestrator	 :8080
COMPUTING_POWER	Количество горутин в Agent	          1
JWT_SECRET	      Секрет для подписи JWT	             secret_key

Примеры запросов (curl)
1. Регистрация

curl -i -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"login":"user1","password":"pass123"}'

Ответ:
HTTP/1.1 200 OK
OK

2. Логин

curl -i -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user1","password":"pass123"}'

Ответ:
HTTP/1.1 200 OK
{
  "token": "<ваш_JWT_токен>"
}

3. Вычисление выражения

curl -i -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"expression":"2+2*2"}'

Ответ:
HTTP/1.1 200 OK
{"result":6}

4. Список всех выражений

curl -i -X GET http://localhost:8080/api/v1/expressions \
  -H "Authorization: Bearer <TOKEN>"

Ответ:
HTTP/1.1 200 OK
{
  "expressions":[
    {"id":1,"status":"done","result":6},
    {"id":2,"status":"error","result":0}
  ]
}

5. Детали конкретного выражения

curl -i -X GET http://localhost:8080/api/v1/expressions/1 \
  -H "Authorization: Bearer <TOKEN>"

Ответ:

HTTP/1.1 200 OK
{
  "expression":{"id":1,"status":"done","result":6}
}

Примеры в Postman

Создайте коллекцию Distributed Calculator.
Добавьте Request Register:

Method: POST
URL: {{baseUrl}}/api/v1/register
Body → raw JSON { "login":"user1","password":"pass123" }

Добавьте Request Login:

Method: POST
URL: {{baseUrl}}/api/v1/login
Body → raw JSON { "login":"user1","password":"pass123" }

Tests:

pm.environment.set("jwt", pm.response.json().token);


Добавьте Request Calculate:

Method: POST
URL: {{baseUrl}}/api/v1/calculate

Headers:
Authorization: Bearer {{jwt}}
Content-Type: application/json

Body: { "expression":"2+2*2" }

Request ListExpressions:

Method: GET
URL: {{baseUrl}}/api/v1/expressions
Header: Authorization: Bearer {{jwt}}

Request GetExpressionById:

Method: GET
URL: {{baseUrl}}/api/v1/expressions/1
Header: Authorization: Bearer {{jwt}}

В переменных окружения (Environments) задайте:

baseUrl = http://localhost:8080

Тестирование
Unit-тесты:

go test ./pkg/calculation
go test ./internal/storage
go test ./internal/agent
go test ./internal/server

Integration-тесты:
go test ./tests/integration
