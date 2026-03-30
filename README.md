# go-web

Минимальное веб-приложение на Go: MVC, роутинг, PostgreSQL, Docker.

## Стек

- **Go** — stdlib + chi (роутинг)
- **PostgreSQL** — БД
- **Docker Compose** — окружение

## Структура

```
cmd/app/          — точка входа
internal/
  db/             — фасад БД
  handlers/       — контроллеры (MVC: C)
  models/         — модели (MVC: M)
  router/         — маршруты
resources/views/  — шаблоны (MVC: V): layouts/, pages/
```

## Запуск

```bash
docker compose up --build
```

Приложение: http://localhost:8081

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_DSN` | Connection string PostgreSQL | `postgres://user:pass@db:5432/app` |
| `PORT` | Порт HTTP-сервера | `8080` |

## Локальная разработка

```bash
go run ./cmd/app
```

Требуется запущенный PostgreSQL (или `docker compose up db`).
