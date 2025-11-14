# QA API Service

Небольшой HTTP-сервис вопросов/ответов на Go c хранением в PostgreSQL и
миграциями через goose.

## Логика работы

-   При старте сервис:
    -   читает переменные окружения `POSTGRES_HOST`, `POSTGRES_PORT`,
        `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` (есть
        дефолты: `qa_user` / `qa_password` / `qa_db`, хост `localhost`,
        порт `5432`);
    -   собирает DSN и подключается к PostgreSQL;
    -   запускает миграции из каталога `./migrations` (`goose up`);
    -   если передан флаг `-reset-db`, сначала делает `goose down-to 0`,
        затем `goose up`;
    -   поднимает HTTP-сервер на порту `:8080`.

### Основные эндпоинты

**Вопросы**

-   `GET /questions/` --- список всех вопросов.
-   `POST /questions/`
    -   тело: `{ "text": "текст вопроса" }`
    -   ответ: созданный вопрос.
-   `GET /questions/{id}/`
    -   ответ: `{ "question": {...}, "answers": [...] }` --- вопрос и
        все его ответы.
-   `DELETE /questions/{id}/`
    -   удаляет вопрос; ответы удаляются каскадно логикой SQL (ON DELETE
        CASCADE).

**Ответы**

-   `POST /questions/{id}/answers/`
    -   тело: `{ "user_id": "строка", "text": "текст ответа" }`
    -   ответ: созданный ответ.
-   `GET /answers/{id}/` --- получить ответ по id.
-   `POST /answers/`
    -   тело:
        `{ "question_id": 1, "user_id": "строка", "text": "текст ответа" }`
-   `DELETE /answers/{id}/` --- удалить ответ.

Все ответы отдаются в JSON.

## Запуск в Docker

### Подготовка `.env`

    POSTGRES_USER=qa_user
    POSTGRES_PASSWORD=qa_password
    POSTGRES_DB=qa_db
    POSTGRESS_HOST=localhost
    POSTGRESS_PORT=5432

### Сборка и запуск

    docker build -t qa-api-server:latest .
    docker compose build
    docker compose up -d

API доступно на `http://localhost:8080`.