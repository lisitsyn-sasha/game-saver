# Используем официальный образ Go
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN go build -o game-saver ./cmd/main.go

# Финальный контейнер
FROM alpine:latest
WORKDIR /root/

# Устанавливаем зависимости
RUN apk --no-cache add ca-certificates

# Копируем бинарник из builder-контейнера
COPY --from=builder /app/game-saver .

# Устанавливаем переменную окружения для конфига
ENV CONFIG_PATH=/config/local.yaml

# Копируем конфиг в контейнер
COPY config/local.yaml /config/local.yaml

# Открываем порт
EXPOSE 8082

# Запускаем приложение
CMD ["./game-saver"]