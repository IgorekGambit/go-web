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

## Отладка (Delve)

Образ собран с [Delve](https://github.com/go-delve/delve). Порт **40000** — для подключения удалённого отладчика (IDE / `dlv connect`).

**Обычный запуск** (по умолчанию): контейнер стартует с `./app`.

**Режим отладки** — переопредели команду и пробрось порт:

```bash
docker compose run --service-ports app dlv exec ./app --headless --listen=:40000 --api-version=2
```

Либо в `docker-compose.yml` добавь сервис/профиль с `command` и `ports: - "40000:40000"`.

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
