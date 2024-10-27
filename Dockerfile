# Используем минималистичный образ Go
FROM golang:1.23.2-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем goose для миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем все файлы в контейнер
COPY . .

# Компилируем Go-приложение
RUN go build -o migrate ./cmd/bot/main.go

# Стандартная команда для выполнения миграций
ENTRYPOINT ["./migrate"]
