# Используем минималистичный образ Go для сборки приложения
FROM golang:1.23.2-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы модуля и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта в рабочую директорию контейнера
COPY . .

# Компилируем Go-приложение в бинарный файл main
RUN go build -o main ./cmd/bot/main.go

# Устанавливаем команду по умолчанию для запуска контейнера
ENTRYPOINT ["./main"]
