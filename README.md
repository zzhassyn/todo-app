# todo-app

Полноценное To-Do приложение: бэкенд на Go (Clean Architecture,
Feature-Driven) и SPA-клиент на React + Vite.

```
.
├── cmd/, internal/, migrations/   ← Go-бэкенд (этот README)
└── frontend/                      ← React-клиент, см. frontend/README.md
```

## Архитектура бэкенда

Проект разделён на два верхнеуровневых блока:

- **`internal/core`** — переиспользуемое ядро, не зависящее от бизнес-логики конкретных фич: домен-примитивы (`domain`), общие ошибки (`errors`), логгер (`logger`), JWT и контекст аутентификации (`auth`), абстракция над пулом соединений Postgres (`repository/postgres/pool`) и HTTP-инфраструктура (`transport/http`: сервер, роутер, middleware, request/response, утилиты).
- **`internal/features`** — независимые друг от друга вертикальные срезы (фичи). Каждая фича содержит свои слои `repository` → `service` → `transport`, организованные по правилу зависимостей Clean Architecture: внешние слои зависят от внутренних через интерфейсы, объявленные во внутреннем слое.

Реализованы три фичи:

```
internal/features/
├── users/   # профили пользователей (CRUD без регистрации)
├── auth/    # регистрация, вход, текущая сессия
└── tasks/   # задачи (To-Do), всегда привязаны к автору
```

### Слои внутри фичи

- **`service`** — бизнес-логика. Объявляет интерфейс `*Repository`, который реализует слой `repository`, и (при необходимости) узкие интерфейсы для зависимостей от других фич — например `tasks_service.UsersChecker` и `auth_service.UsersRegistry`, оба реализуются `users_service.UsersService`, но фичи `tasks` и `auth` не импортируют пакет `users` напрямую.
- **`repository/postgres`** — реализация доступа к данным на pgx/Postgres. Возвращает и принимает только `domain`-сущности.
- **`transport/http`** — HTTP-хендлеры, DTO и маршруты (`Routes()`), которые регистрируются в `core_http_server.APIVersionRouter`. Маршруты, требующие аутентификации, получают `core_http_middleware.Auth(...)` через поле `Route.Middleware`.
- **`domain`** (в `core`, общий для всех фич) — сущности, инварианты и бизнес-правила валидации (`Validate()`, `ApplyPatch()`, `Complete()`/`Uncomplete()` и т.д.). Это самый внутренний, ничего не знающий о фреймворках слой.

## Аутентификация

JWT (HS256, минимальная реализация на стандартной библиотеке — `internal/core/auth`,
без внешних зависимостей) выдаётся в **httpOnly cookie** после `/auth/register`
или `/auth/login`. Middleware `core_http_middleware.Auth` проверяет cookie на
каждом защищённом маршруте и кладёт claims (`user_id`, `email`) в контекст
запроса; хендлеры читают их через `core_auth.FromContext`.

- Все `tasks`-маршруты и большинство `users`-маршрутов требуют аутентификации.
- `tasks` всегда скоупится на текущего пользователя: `author_user_id` берётся
  из токена, а не из тела запроса или query-параметра, и сервис возвращает
  `404` (не `403`) при попытке обратиться к чужой задаче — это не даёт
  угадывать чужие ID задач.
- CORS настроен на конкретный origin (`HTTP_ALLOWED_ORIGIN`) с
  `Access-Control-Allow-Credentials: true`, что обязательно для работы с
  cookie из браузера.

| Метод | Путь | Защищён | Описание |
|---|---|---|---|
| POST | `/api/v1/auth/register` | нет | Регистрация, выдаёт cookie |
| POST | `/api/v1/auth/login` | нет | Вход, выдаёт cookie |
| POST | `/api/v1/auth/logout` | нет | Очищает cookie |
| GET | `/api/v1/auth/me` | да | Текущий пользователь |

## Фича: Users

| Метод | Путь | Описание |
|---|---|---|
| GET | `/api/v1/users?limit=&offset=` | Список пользователей |
| GET | `/api/v1/users/{id}` | Получить пользователя |
| PATCH | `/api/v1/users/{id}` | Частично обновить пользователя |
| DELETE | `/api/v1/users/{id}` | Удалить пользователя |

Создание пользователя выполняется только через `POST /auth/register` —
отдельного публичного `POST /users` нет, чтобы не было двух путей создания
учётной записи с разной обработкой пароля.

## Фича: Tasks

CRUD задач + явные операции завершения/возобновления. Optimistic locking
через `version`.

| Метод | Путь | Описание |
|---|---|---|
| POST | `/api/v1/tasks` | Создать задачу (автор — текущий пользователь) |
| GET | `/api/v1/tasks?limit=&offset=&completed=` | Список своих задач, фильтр по статусу |
| GET | `/api/v1/tasks/{id}` | Получить свою задачу |
| PATCH | `/api/v1/tasks/{id}` | Частично обновить `title`/`description` |
| DELETE | `/api/v1/tasks/{id}` | Удалить задачу |
| POST | `/api/v1/tasks/{id}/complete` | Отметить выполненной (`completed_at = now()`) |
| POST | `/api/v1/tasks/{id}/uncomplete` | Снять отметку о выполнении |

### Пример: зарегистрироваться и сохранить cookie

```bash
curl -c cookies.txt -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Ada Lovelace", "email": "ada@example.com", "password": "hunter22"}'
```

### Пример: создать задачу (используя сохранённую cookie)

```bash
curl -b cookies.txt -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко", "description": "2 литра, можно растительное"}'
```

### Пример: список незавершённых задач

```bash
curl -b cookies.txt "http://localhost:8080/api/v1/tasks?completed=false&limit=20&offset=0"
```

## Обработка ошибок

Все ошибки бизнес-логики и репозитория оборачиваются в `core_errors`
(`ErrNotFound`, `ErrInvalidArgument`, `ErrConflict`, `ErrUnauthorized`) и
транслируются HTTP-слоем в соответствующие коды:

| Domain error | HTTP статус |
|---|---|
| `ErrInvalidArgument` | 400 |
| `ErrUnauthorized` | 401 |
| `ErrNotFound` | 404 |
| `ErrConflict` (в т.ч. конфликт версии при PATCH, дубликат email) | 409 |
| прочее / panic | 500 |

## Запуск бэкенда

```bash
cp .env.example .env   # заполнить POSTGRES_USER/PASSWORD/DB и AUTH_JWT_SECRET
make env-up            # поднять Postgres в Docker
make migrate-up        # применить миграции (000001_init, 000002_auth)
make todoapp-run       # собрать и запустить приложение
```

`AUTH_JWT_SECRET` обязателен и не имеет значения по умолчанию — сгенерируйте
случайную строку, например `openssl rand -base64 48`.

Миграции создают схему `todoapp` с таблицами `users` (включая `email` и
`password_hash` после `000002_auth`) и `tasks`.

## Запуск фронтенда

См. [`frontend/README.md`](./frontend/README.md). Короткая версия:

```bash
cd frontend
npm install
npm run dev
```

## Диагностика частых проблем

- **Фронтенд пишет "не получилось подключиться к серверу"** — почти всегда
  значит, что Go-бэкенд не запущен или недоступен по адресу, заданному в
  `frontend/.env` (`VITE_API_BASE_URL`, по умолчанию
  `http://localhost:8080/api/v1`). Проверьте: `curl http://localhost:8080/api/v1/auth/me`
  должен вернуть `401` (а не `connection refused`).
- **После входа сразу разлогинивает / задачи не подгружаются** — проверьте
  `HTTP_ALLOWED_ORIGIN` на бэкенде: он должен **точно** совпадать с origin,
  с которого открыт фронтенд (порт включительно), иначе браузер не примет
  cookie с токеном даже при успешном ответе сервера.
- **CORS-ошибка в консоли браузера** — то же самое: `HTTP_ALLOWED_ORIGIN`
  не совпадает с реальным origin фронтенда, либо фронтенд ходит на бэкенд
  без `credentials: "include"` (в этом проекте уже включено по умолчанию).

## Стек

Бэкенд:
- Go 1.25, стандартный `net/http` (`http.ServeMux` с path-параметрами)
- `pgx/v5` — драйвер и пул соединений Postgres
- `go.uber.org/zap` — структурированное логирование
- `go-playground/validator` — валидация входящих HTTP-запросов
- `golang.org/x/crypto/bcrypt` — хеширование паролей
- собственная минимальная реализация JWT (HS256) на стандартной библиотеке — без внешних JWT-зависимостей
- `golang-migrate` — миграции БД (через Docker)

Фронтенд:
- React 19 + Vite
- без UI-фреймворков — design tokens на чистых CSS-переменных
