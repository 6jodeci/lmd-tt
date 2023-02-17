# Используем образ Golang для сборки бинарника
FROM golang:1.17 as builder

# Копируем исходный код в рабочую директорию контейнера
COPY . /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Используем образ alpine в качестве окончательного образа
FROM alpine:latest

# Копируем бинарник в контейнер из образа builder
COPY --from=builder /app/app /app

# Устанавливаем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["/app"]
