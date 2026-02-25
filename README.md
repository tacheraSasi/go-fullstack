# GO-FullStack

Production-ready fullstack Go application built with **Gin + GORM + Templ + Alpine.js + Tailwind CSS**. Includes JWT auth, user dashboard, RBAC, customer/invoice modules, and Docker support.

## Features

**Backend**
- Gin HTTP API with clean repository → service → handler layering
- JWT authentication (login, register, forgot/reset password, logout)
- Role-based access control (users, roles, permissions)
- Customer and invoice CRUD modules
- SQLite / Postgres / MySQL support via GORM
- Structured logging with Logrus
- Health and readiness endpoints

**Frontend**
- Server-rendered pages with [Templ](https://templ.guide)
- [templUI](https://templui.io) component library (cards, forms, inputs, badges, icons, etc.)
- [Alpine.js](https://alpinejs.dev) for client-side interactivity
- [Tailwind CSS v4](https://tailwindcss.com) with dark mode
- HTMX included for progressive enhancement

**Auth Flow (Laravel Breeze-style)**
- `/auth/login` — login with email & password, stores JWT in localStorage
- `/auth/register` — registration with client-side validation
- `/auth/forgot-password` — request a password reset token
- `/auth/reset-password` — reset password with token

**User Dashboard**
- `/dashboard` — welcome page with user info, account status, roles, member-since date
- `/dashboard/settings` — edit profile (name, email) and change password
- Sidebar layout with auth guard (redirects to login if no token)

## Quick Start

### Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev) (recommended) or Make
- Node.js (for Tailwind CSS CLI)

### 1) Configure environment

```bash
cp .env.example .env
```

Adjust values in `.env` — especially `JWT_SECRET` for production.

### 2) Run in development

```bash
task dev
```

This starts Tailwind CSS in watch mode and Templ with hot reload + proxy. The app runs at `http://localhost:7331` (proxy) with the API on port `8090`.

Alternatively with Make:

```bash
make dev
```

### 3) Seed sample data (optional)

```bash
make seed
```

### 4) Build for production

```bash
make build
```

## Available Commands

**Taskfile (recommended):**

```bash
task dev          # Start dev server with hot reload + Tailwind watch
task templ        # Run templ generate with watch & proxy
task tailwind     # Watch Tailwind CSS changes
```

**Makefile:**

```bash
make help           # Show all available commands
make dev            # Development mode with hot reload
make build          # Build the binary
make run            # Build and run
make seed           # Seed database
make test           # Run tests
make lint           # Run linter
make docker-up      # Start with docker-compose
make docker-down    # Stop docker-compose
```

## Project Structure

```
cmd/
  api/            → Application entrypoint
  seed/           → Database seeder
internals/
  config/         → Environment configuration
  dtos/           → Request/response DTOs with validation
  handlers/       → HTTP handlers (controllers)
  middlewares/    → Auth, CORS, logging, admin middleware
  models/         → GORM models (User, Role, Permission, Customer, Invoice, etc.)
  repositories/   → Data access layer
  services/       → Business logic layer
  utils/          → Response helpers
components/       → templUI components (accordion, badge, button, card, form, icon, input, etc.)
ui/
  layouts/        → Shared layouts (BaseLayout, DashboardLayout)
  pages/          → Page templates (home, login, register, dashboard, settings, etc.)
pkg/
  database/       → DB connection and migration
  jwt/            → JWT generation and validation
  logger/         → Structured logger setup
  styles/         → Terminal styling
assets/
  css/            → Tailwind input.css and generated output.css
  js/             → Component JS (templUI minified scripts)
```

## API Route Overview

Base path: `/api/v1`

| Group | Routes | Auth |
|---|---|---|
| Public | `POST /login`, `POST /register`, `POST /forgot-password`, `POST /reset-password` | None |
| Protected | `POST /logout`, `GET/PUT /users/:id`, `PUT /users/:id/password`, `GET /users/:id/roles` | JWT |
| Protected | `GET/POST /customers`, `GET/PUT/DELETE /customers/:id` | JWT |
| Protected | `GET/POST /invoices`, `GET/PUT/DELETE /invoices/:id` | JWT |
| Admin | `GET /admin/users`, `DELETE /admin/users/:id`, `POST/DELETE /admin/users/:id/roles/:roleId` | JWT + Admin |
| Admin | CRUD `/admin/roles/*`, `/admin/permissions/*` | JWT + Admin |

**Web Pages:**

| Route | Description |
|---|---|
| `/` | Landing page |
| `/auth/login` | Login |
| `/auth/register` | Register |
| `/auth/forgot-password` | Forgot password |
| `/auth/reset-password` | Reset password |
| `/dashboard` | User dashboard (auth required) |
| `/dashboard/settings` | Profile & password settings |
| `/health` | Liveness check |
| `/health/ready` | Readiness check |

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080` | HTTP server port |
| `GIN_MODE` | `release` | Gin mode (`debug`, `release`, `test`) |
| `JWT_SECRET` | `secret` | JWT signing secret (**change in production**) |
| `JWT_EXPIRES_IN` | `24` | Token expiry in hours |
| `CORS_ALLOWED_ORIGINS` | `*` | Comma-separated origins or `*` |
| `LOG_FILE_PATH` | `logs/app.log` | Log file output |
| `DB_TYPE` | `sqlite` | `sqlite`, `postgres`, or `mysql` |
| `DB_PATH` | `core.db` | SQLite database file path |
| `DB_HOST` | `localhost` | DB host (Postgres/MySQL) |
| `DB_PORT` | `5432` | DB port (Postgres/MySQL) |
| `DB_USER` | `user` | DB user (Postgres/MySQL) |
| `DB_PASSWORD` | `password` | DB password (Postgres/MySQL) |
| `DB_NAME` | `dbname` | DB name (Postgres/MySQL) |

## Docker

```bash
make docker-up
```

See [README.Docker.md](README.Docker.md) for details.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | Gin |
| ORM | GORM |
| Templating | Templ |
| Components | templUI |
| Interactivity | Alpine.js |
| Styling | Tailwind CSS v4 |
| Auth | JWT (golang-jwt/jwt/v5) |
| Database | SQLite (default), Postgres, MySQL |
| Logging | Logrus |

## Notes

- Auto-migrations run on startup — no manual SQL needed.
- Admin routes require both JWT authentication and the admin role.
- The user dashboard is client-side protected via Alpine.js auth guards — unauthenticated users are redirected to `/auth/login`.
- After login, users are redirected to `/dashboard`.
