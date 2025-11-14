FROM golang:1.25-alpine AS builder

WORKDIR /app

# 1. зависимости
COPY go.mod go.sum ./
RUN go mod download

# 2. исходники
COPY . .

# 3. сборка только main-пакета (каталог .)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Рантайм-слой ----
FROM alpine:3.19

WORKDIR /app

# бинарник
COPY --from=builder /app/server .

# миграции, чтобы goose их видел по пути ./migrations
COPY migrations ./migrations

EXPOSE 8080

# дефолтные значения (можешь переопределять через docker run / compose)
ENV POSTGRES_USER=qa_user \
    POSTGRES_PASSWORD=qa_password \
    POSTGRES_DB=qa_db

CMD ["./server"]
