# HackathonHub

Участники находят хакатоны, регистрируются, собирают команду и сдают решение до дедлайна. Организаторы создают хакатоны, настраивают треки и кейсы, смотрят заявки и сабмиты.

Стек: **Go** (API) + **React** (SPA) + **PostgreSQL**. Спека API лежит в [`openapi.yaml`](openapi.yaml).

**Развернутое решение:** [web](https://web-production-196fd.up.railway.app) · [API](https://api-production-3cf0.up.railway.app/api/v1/health)

---

## Фичи

- Регистрация / вход, JWT.
- Каталог опубликованных хакатонов, карточка с треками и кейсами
- Участник: регистрация на хакатон, команда (создать / вступить / роли), сдача submission (ссылки на репо, демо, презентацию)
- Организатор: CRUD хакатонов, треков, кейсов, публикация, списки регистраций и сабмитов
- При старте API поднимаются миграции (goose) и seed с тестовыми данными

Полный перечень эндпоинтов - в OpenAPI.

---

## Структура репозитория

```
.
├── backend/                 # Go API
│   ├── cmd/
│   │   ├── app/             # точка входа HTTP-сервера
│   │   └── seed/            # сид
│   ├── internal/
│   │   ├── delivery/http/   # хендлеры (chi)
│   │   ├── usecase/         # бизнес-логика
│   │   ├── adapter/repository/postgres_repo/
│   │   ├── domain/          # entities, dto, порты репозиториев
│   │   ├── infrastructure/  # postgres, jwt, seed, миграции
│   │   └── server/          # сборка сервера и роутов
│   ├── pkg/                 # middleware, логгер, http-хелперы
│   ├── Dockerfile
│   └── railway.toml
├── frontend/                # React + Vite + TypeScript
│   ├── src/
│   │   ├── api/             # клиент к REST
│   │   ├── pages/           # экраны
│   │   ├── components/
│   │   └── context/         # AuthContext
│   ├── Dockerfile           # nginx со статикой
│   └── railway.toml
├── docker/                  # Dockerfile'ы для деплоя из корня (monorepo)
├── openapi.yaml
├── docker-compose.yml
├── Makefile
└── .env.example
```

Бэкенд разложен по слоям: **handler -> usecase -> repository -> postgres**. Фронт ходит в API через `fetch`, токен подставляется в заголовок `Authorization: Bearer`.

---

## Как это связано

```
   Клиент
        │
        ▼
   Go API
        │
        ▼
   PostgreSQL
```

1. Пользователь логинится - API отдаёт JWT, фронт кладёт его в `localStorage`.
2. Защищённые запросы идут с заголовком `Authorization`. Роль `organizer` открывает `/organizer/*` на API и отдельные страницы на фронте.
3. При старте API сам накатывает SQL-миграции из `backend/internal/infrastructure/database/postgres/migrations/`.
4. Если `SEED_DISABLE` не `1`, создаются демо-пользователи и опубликованный хакатон.

---

## Требования

| Инструмент | Версия (ориентир)                      |
| ---------- | -------------------------------------- |
| Go         | 1.23+                                  |
| Node.js    | 20+                                    |
| Docker     | для Postgres локально                  |
| Make       | опционально, но проще всего через него |

На Windows без Make можно выполнять те же команды вручную (см. ниже). (ну или wsl\MSYS2 накатить тоже поможет, в общем любую unix подобную среду :D )

---

## Быстрый старт

### macOS / Linux

```bash
git clone https://github.com/virtualreus/PiRKSPCourseWork.git
cd PiRKSPCourseWork

cp .env.example .env
cp frontend/.env.example frontend/.env

make dev
```

Откроется:

- фронт: http://localhost:5173
- API: http://localhost:8080/api/v1/health

Остановить API и Vite - `Ctrl+C`. Postgres в Docker после этого ещё крутится; выключить: `make down`.

### Windows (PowerShell)

Нужны Go, Node, Docker Desktop.

```powershell
git clone https://github.com/virtualreus/PiRKSPCourseWork.git
cd PiRKSPCourseWork

Copy-Item .env.example .env
Copy-Item frontend\.env.example frontend\.env

docker compose up -d --wait

cd backend
go mod download
go run ./cmd/app
```

Во втором терминале:

```powershell
cd frontend
npm install
npm run dev
```

`make dev` на Windows можно через [Git Bash](https://git-scm.com/) или WSL - тогда всё как на Linux.

---

## Тестовые аккаунты (seed)

Создаются при первом запуске API, если в `.env` не стоит `SEED_DISABLE=1`.

| Роль        | Email            | Пароль  |
| ----------- | ---------------- | ------- |
| Организатор | `admin@admin.ru` | `admin` |
| Участник    | `user@user.ru`   | `user`  |

Организатор: раздел «Организатор» в шапке. Участник: каталог -> хакатон -> участие / команда / сдача работы.

Пересоздать только сид без пересоздания БД: `make seed` (нужен запущенный Postgres).

---

## Переменные окружения

**Корень (`.env`)** - бэкенд и docker-compose:

| Переменная     | Назначение                    |
| -------------- | ----------------------------- |
| `PORT`         | порт API (по умолчанию 8080)  |
| `PG_DSN`       | строка подключения к Postgres |
| `JWT_SECRET`   | секрет для подписи токенов    |
| `CORS_ORIGIN`  | origin фронта для CORS        |
| `SEED_DISABLE` | `1` - не сыпать демо-данные   |

**`frontend/.env`:**

| Переменная     | Назначение                                               |
| -------------- | -------------------------------------------------------- |
| `VITE_API_URL` | базовый URL API, например `http://localhost:8080/api/v1` |

На Railway у API обычно `DATABASE_URL=${{Postgres.DATABASE_URL}}`, у фронта при сборке прокидывается `VITE_API_URL`.

---

## Makefile help

```bash
make setup       # go mod + npm install + .env
make db          # только Postgres в Docker
make dev         # postgres + API + frontend
make backend     # только API
make frontend    # только Vite
make seed        # демо-данные
make down        # остановить postgres
```

---

## API

Базовый префикс: `/api/v1`.

| Группа           | Примеры                                                                    |
| ---------------- | -------------------------------------------------------------------------- |
| **Auth**         | `POST /auth/register`, `POST /auth/login`                                  |
| **Публичное**    | `GET /hackathons`, `GET /hackathons/{id}`, `GET /hackathons/{id}/teams`    |
| **Пользователь** | `GET /users/me`, `GET /users/me/dashboard`, `PATCH /users/me`              |
| **Участие**      | `POST /hackathons/{id}/register`, команды, `PUT /teams/{id}/submission`    |
| **Организатор**  | `/organizer/hackathons`, треки, кейсы, publish, registrations, submissions |

Авторизация: `Authorization: Bearer <token>`.

healthcheck'ер `GET /health` -> `{"status":"ok","database":"ok"}`.

Детали, схемы тел запросов и коды ошибок - в [`openapi.yaml`](openapi.yaml).

Пример:

```bash
curl http://localhost:8080/api/v1/health

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@user.ru","password":"user"}'
```

---

## Деплой (Railway)

В проекте три сервиса: **Postgres**, **api**, **web**. Monorepo: для GitHub-деплоя у `api` и `web` лучше указать Root Directory (`backend` / `frontend`) или использовать Dockerfile из `docker/`.

```bash
npm install
make railway-login
make railway-deploy
```

Отдельно: `make railway-api`, `make railway-web`. URL продакшена подставь в `Makefile` (`RAILWAY_API_URL`, `RAILWAY_WEB_URL`) или в переменные сервисов в панели Railway.

---

## Разработка

**Миграции** - SQL в `backend/internal/infrastructure/database/postgres/migrations/hackathon/`, накатываются при старте API.

**Линт фронта:** `cd frontend && npm run lint`

**Сборка фронта:** `cd frontend && npm run build` -> `frontend/dist`

**Сборка API в Docker:**

```bash
cd backend && docker build -t hackathon-api .
```

---

## Автор

Тищенко Никита Сергеевич, nt3008@yandex.ru
