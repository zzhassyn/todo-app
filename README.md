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

Реализованы четыре фичи:

```
internal/features/
├── users/    # профили пользователей (CRUD без регистрации)
├── auth/     # регистрация, вход, текущая сессия
├── tasks/    # задачи (To-Do), всегда привязаны к автору
└── folders/  # папки/списки для группировки задач
```

### Слои внутри фичи

- **`service`** — бизнес-логика. Объявляет интерфейс `*Repository`, который реализует слой `repository`, и (при необходимости) узкие интерфейсы для зависимостей от других фич — например `tasks_service.UsersChecker`, `tasks_service.FoldersChecker` и `auth_service.UsersRegistry`, реализуемые соответственно `users_service.UsersService` и `folders_service.FoldersService`, но фичи `tasks`/`auth` не импортируют пакеты `users`/`folders` напрямую.
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
через `version`. Задачи поддерживают приоритет (`low`/`medium`/`high`),
необязательный срок выполнения (`due_at`), теги (многие-ко-многим,
таблицы `tags`/`task_tags`, общий словарь тегов на всё приложение) и
принадлежность папке (`folder_id`, см. фичу Folders ниже — задача без
папки считается лежащей в общем "буфере"). Удаление — мягкое: обычный
`DELETE`/`archive` помечают задачу `archived_at`, ничего физически не
удаляя; безвозвратное удаление — отдельный шаг, доступный только из
архива (см. `DELETE .../permanent` ниже).

| Метод | Путь | Описание |
|---|---|---|
| POST | `/api/v1/tasks` | Создать задачу (автор — текущий пользователь) |
| GET | `/api/v1/tasks?limit=&offset=&completed=&archived=&priority=&tag=&folder_id=` | Список своих задач с фильтрами |
| GET | `/api/v1/tasks/{id}` | Получить свою задачу (включая архивную) |
| PATCH | `/api/v1/tasks/{id}` | Частично обновить `title`/`description`/`priority`/`due_at`/`tags`/`folder_id` |
| DELETE | `/api/v1/tasks/{id}` | Архивировать задачу (синоним `POST .../archive`) |
| POST | `/api/v1/tasks/{id}/complete` | Отметить выполненной (`completed_at = now()`) |
| POST | `/api/v1/tasks/{id}/uncomplete` | Снять отметку о выполнении |
| POST | `/api/v1/tasks/{id}/archive` | Архивировать (мягкое удаление) |
| POST | `/api/v1/tasks/{id}/unarchive` | Восстановить из архива |
| DELETE | `/api/v1/tasks/{id}/permanent` | Безвозвратно удалить — **только если задача уже архивирована** (иначе `409`) |
| GET | `/api/v1/tags` | Список всех тегов (для автокомплита) |

`archived` в `GET /tasks` по умолчанию (или при `archived=false`) показывает
только неархивные задачи; `archived=true` — только архивные. Оба набора
одновременно не возвращаются за один запрос.

`folder_id` в `GET /tasks` принимает либо конкретный UUID папки (только
задачи в ней), либо специальное значение `none` (только задачи без папки —
"буфер"); если параметр не передан, фильтрация по папке не применяется
вовсе. Чужой (несуществующий или принадлежащий другому пользователю)
`folder_id`, переданный при создании/патче задачи или как фильтр списка,
даёт `400 Bad Request`, а не `403` — папки других пользователей нельзя
отличить по ответу от просто не существующих.

Поле `tags` в `PATCH` следует тому же соглашению, что и везде у `Nullable`
полей: отсутствие поля или `null` — теги не трогать; `[]` — снять все теги;
`["a", "b"]` — заменить набор тегов на указанный. Несуществующие теги
создаются автоматически (find-or-create по имени). Поле `folder_id` в
`PATCH` следует тому же соглашению: отсутствие поля — папку не трогать;
явный `null` — убрать из папки (вернуть в буфер); UUID — переместить в
эту папку.

### Пример: зарегистрироваться и сохранить cookie

```bash
curl -c cookies.txt -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Ada Lovelace", "email": "ada@example.com", "password": "hunter22"}'
```

### Пример: создать задачу с приоритетом, сроком, тегами и папкой

```bash
curl -b cookies.txt -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Купить молоко",
    "priority": "high",
    "due_at": "2026-07-01T12:00:00Z",
    "tags": ["errands", "urgent"],
    "folder_id": "5882de4e-580d-4cfe-b7c5-9c87279a1729"
  }'
```

### Пример: список незавершённых задач с тегом `work`

```bash
curl -b cookies.txt "http://localhost:8080/api/v1/tasks?completed=false&tag=work&limit=20&offset=0"
```

### Пример: список задач без папки (буфер)

```bash
curl -b cookies.txt "http://localhost:8080/api/v1/tasks?folder_id=none"
```

### Пример: архивировать, затем безвозвратно удалить задачу

```bash
curl -b cookies.txt -X POST http://localhost:8080/api/v1/tasks/1/archive
# повторный DELETE /tasks/1 (без /archive) до этого шага вернул бы 409
curl -b cookies.txt -X DELETE http://localhost:8080/api/v1/tasks/1/permanent
```

## Фича: Folders

Папки/списки для группировки задач — лёгкая фича над `tasks`: создание,
список своих папок, удаление. У папки нет мягкого удаления, в отличие от
задач: `DELETE /folders/{id}` — это **настоящее, безвозвратное удаление**,
и благодаря `ON DELETE CASCADE` на `tasks.folder_id` оно **физически
удаляет все задачи внутри папки, включая уже архивированные**, минуя
обычную защиту "сначала архивируй, потом удаляй навсегда" у задач. Это
осознанный выбор (а не побочный эффект), сделанный по явному запросу;
если в будущем понадобится более бережное поведение, альтернатива —
`ON DELETE SET NULL` (задачи становятся неприписанными к папке вместо
удаления).

| Метод | Путь | Описание |
|---|---|---|
| POST | `/api/v1/folders` | Создать папку (привязана к текущему пользователю) |
| GET | `/api/v1/folders` | Список своих папок |
| DELETE | `/api/v1/folders/{id}` | Удалить папку **и каскадно все задачи внутри неё** |

`id` папки — UUID, генерируемый на стороне Go (`google/uuid`, уже
использовавшийся в проекте для request-id) перед вставкой в БД, а не
автоинкрементный `SERIAL`, как у остальных сущностей — это единственная
сущность в проекте с UUID-идентификатором.

### Пример: создать папку и задачу в ней

```bash
FOLDER_ID=$(curl -s -b cookies.txt -X POST http://localhost:8080/api/v1/folders \
  -H "Content-Type: application/json" -d '{"title": "Работа"}' | jq -r .id)

curl -b cookies.txt -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d "{\"title\": \"Подготовить отчёт\", \"folder_id\": \"$FOLDER_ID\"}"
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
make migrate-up        # применить миграции (000001_init, 000002_auth, 000003_task_extras, 000004_folders)
make todoapp-run       # собрать и запустить приложение
```

`AUTH_JWT_SECRET` обязателен и не имеет значения по умолчанию — сгенерируйте
случайную строку, например `openssl rand -base64 48`.

Миграции создают схему `todoapp` с таблицами `users` (включая `email` и
`password_hash` после `000002_auth`), `tasks` (включая `priority`, `due_at`,
`archived_at` после `000003_task_extras`, и `folder_id` после
`000004_folders`), `tags`/`task_tags`, а также `folders` (после
`000004_folders`).

## Запуск фронтенда

См. [`frontend/README.md`](./frontend/README.md). Короткая версия:

```bash
cd frontend
npm install
npm run dev
```

## CI

`.github/workflows/ci.yml` запускается на каждый push/PR в `main` и состоит
из двух независимых джобов:

- **Backend** — `gofmt -l` (фейлит, если что-то не отформатировано), `go
  vet`, `go build`, `go test` (в проекте пока нет `_test.go` файлов, так что
  это сейчас сводится к проверке, что `go test` не падает на сборке) и
  [`golangci-lint`](https://golangci-lint.run/) (конфиг в `.golangci.yml`;
  отключена только проверка имён пакетов `ST1003`, поскольку в проекте
  намеренно используется snake_case в именах пакетов вида `core_errors`,
  `users_service`, чтобы по импорту сразу было видно, какому слою/фиче он
  принадлежит).
- **Frontend** — `npm ci`, ESLint, `vite build`.

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
- `google/uuid` — генерация UUID для папок (`folders.id`/`tasks.folder_id`); до фичи Folders использовался только для request-id в middleware
- собственная минимальная реализация JWT (HS256) на стандартной библиотеке — без внешних JWT-зависимостей
- `golang-migrate` — миграции БД (через Docker)

Фронтенд:
- React 19 + Vite
- без UI-фреймворков — design tokens на чистых CSS-переменных
