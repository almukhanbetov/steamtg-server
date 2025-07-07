# Используем официальный минимальный Go образ
FROM golang:1.24

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходники
COPY . .

# Собираем приложение
RUN go build -o main .

# Открываем порт
EXPOSE 8989

# Команда запуска
CMD ["./main"]
