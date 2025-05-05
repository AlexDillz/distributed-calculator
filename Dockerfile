# Базовый образ с Go
FROM golang:1.23 AS builder

# Установка protobuf-компилятора
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

# Рабочая директория
WORKDIR /app

# Зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходников
COPY . .

# Генерация gRPC кода
RUN protoc --go_out=. --go-grpc_out=. internal/proto/tasks.proto

# Сборка бинарников
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go
RUN CGO_ENABLED=0 go build -o /agent cmd/agent/main.go

# Финальный образ
FROM gcr.io/distroless/static-debian12

COPY --from=builder /server /server
COPY --from=builder /agent /agent

ENTRYPOINT ["/server"]
