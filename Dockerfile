# Базовый образ с Go
FROM golang:1.23 AS builder

# Установка зависимостей
RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Рабочая директория
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Генерация proto файлов
RUN protoc --go_out=. --go-grpc_out=. internal/proto/tasks.proto

# Сборка сервера и агента
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go
RUN CGO_ENABLED=0 go build -o /agent cmd/agent/main.go

# Финальный образ (минималистичный)
FROM gcr.io/distroless/static-debian12

# Копирование бинарников
COPY --from=builder /server /server
COPY --from=builder /agent /agent

# Запуск по умолчанию (можно переопределить через docker-compose)
CMD ["/server"]