# Todo App

A full-featured To-Do application built with a **Go** backend (following Clean Architecture and Feature-Driven principles) and a modern **React 19 + Vite** Single Page Application (SPA) frontend.

```text
.
├── backend/       # Go backend, entrypoints, and migrations
├── frontend/      # React SPA client
├── docker-compose.yaml
└── Makefile
```

## Features

- **Authentication**: Secure JWT-based authentication using HTTP-only cookies with refresh token rotation.
- **Task Management**: Create, read, update, and soft-delete tasks. Includes priority levels, due dates, and custom tags.
- **Folder Organization**: Group tasks into folders. Cascading deletion is supported.
- **Search & Filtering**: Full-text search (tsvector + GIN indexes) with scope support, along with robust metadata filtering.
- **Recurring Tasks**: Create tasks on a Cron schedule using a robust background worker.
- **Bulk Actions**: Multi-select support allowing you to batch complete, archive, or move multiple tasks simultaneously.
- **Drag and Drop**: Reorder tasks and move them between lists seamlessly using native HTML5 drag and drop.
- **Keyboard Shortcuts**: Power-user friendly with global hotkeys:
  - `j` / `k` for quick task navigation
  - `e` for instant archiving
  - `cmd+k` / `ctrl+k` for jumping to the search bar
- **Optimistic UI**: Snappy interactions using React 19's `useOptimistic` and `startTransition`, updating local state instantly without waiting for network responses.

## Backend Architecture

The Go backend follows a Clean Architecture approach and is divided into two top-level blocks:

1. **`backend/internal/core`**: The reusable core, independent of business logic. It contains domain primitives (`domain`), unified errors (`errors`), logging (`logger`), authentication context (`auth`), PostgreSQL connection pooling (`repository/postgres/pool`), and HTTP infrastructure (server, router, middleware).
2. **`backend/internal/features`**: Independent vertical slices of features. Each feature contains its own `repository` -> `service` -> `transport` layers.

### Implemented Features
- **`users`**: User profiles and management.
- **`auth`**: Registration, login, JWT issuance, and refresh token rotation.
- **`tasks`**: Task lifecycle, soft-deletion, bulk actions, and recurring task generation.
- **`folders`**: Grouping functionality with cascading hard-deletion for contained tasks.

## Tech Stack

**Backend**
- **Go 1.25** with standard `net/http` (`http.ServeMux`)
- **pgx/v5** for PostgreSQL connection pooling
- **go.uber.org/zap** for structured logging
- **go-playground/validator** for HTTP request validation
- **golang.org/x/crypto/bcrypt** for password hashing
- **github.com/robfig/cron/v3** for recurring task scheduling
- **golang-migrate** for database migrations

**Frontend**
- **React 19** + **Vite**
- Custom design tokens with pure CSS (no bloated UI frameworks)
- Modern Hooks (`useOptimistic`, `useActionState`)

**Infrastructure**
- **Docker & Docker Compose** (Full stack: Postgres, Backend, Frontend Nginx)
- **Testcontainers** for repository-level integration tests
- **GitHub Actions** CI/CD pipeline (linting, tests, builds)

## Running Locally

The easiest way to get the full application running is by using Docker Compose.

### 1. Environment Setup
Copy the example environment variables:
```bash
cp .env.example .env
```
Fill in the necessary values inside `.env` (e.g., `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, and `AUTH_JWT_SECRET`).
> **Note:** Generate a secure random string for `AUTH_JWT_SECRET` (e.g., `openssl rand -base64 48`).

### 2. Start the Full Stack (Docker)
You can build and start the entire stack (PostgreSQL, Go Backend, React Frontend) with one command:
```bash
docker compose up --build -d
```
The frontend will be available at `http://localhost:80`, and the backend API will be proxied automatically.

### Alternative: Local Development

If you prefer to run the Go application and Vite dev server directly on your host machine:

**1. Database and Migrations:**
```bash
make env-up
make migrate-up
```

**2. Run the Backend:**
```bash
make todoapp-run      
```

**3. Run the Frontend:**
```bash
cd frontend
npm install
npm run dev          
```

## Troubleshooting

- **"Failed to connect to server" (Frontend)**: Ensure your backend is running and the `VITE_API_BASE_URL` in `frontend/.env` correctly points to the Go server (default: `http://localhost:8080/api/v1`).
- **Logged out immediately / Tasks won't load**: Check `HTTP_ALLOWED_ORIGIN` on the backend. It must **exactly** match the origin URL (including the port) of your frontend server, otherwise the browser will reject the authentication cookie due to CORS policies.
- **Database Connection Errors**: Verify your `TODO_APP_PG_*` variables inside `.env` map exactly to the Postgres instance credentials.
