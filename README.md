<div align="center">

<img src="https://raw.githubusercontent.com/milanz247/vuelang/main/ui/src/assets/vuelang-logo.svg" alt="Vuelang Logo" width="120" height="120" />

# Vuelang Framework

**Enterprise-grade, full-stack MVC web framework**  
Bridging the performance of Go with the reactivity of Vue 3

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)
[![Author](https://img.shields.io/badge/Author-Milan%20Madusanka-blue?style=flat-square)](https://github.com/milanz247)

> **Author**: Milan Madusanka, Associate TechOps Engineer

</div>

---

## What is Vuelang?

Vuelang is an enterprise-grade, full-stack MVC web framework that bridges the raw performance of a **Go backend** with the reactivity of a **Vue 3 frontend**. Inspired by proven MVC architectural patterns, it provides a highly structured, scalable environment for rapid application development — compiled down to a **single executable binary**.

```
  ┌────────────────────────────────────────────────────────┐
  │  Backend  →  Go 1.25 + Gin                            │
  │  Frontend →  Vue 3 + Vite + Tailwind + ShadcnVue      │
  │  Database →  MySQL 8.0                                 │
  │  Auth     →  JWT (access + refresh tokens) + RBAC     │
  │  Deploy   →  Single binary  or  Docker                 │
  └────────────────────────────────────────────────────────┘
```

---

## Table of Contents

- [System Architecture](#system-architecture)
- [Core Features](#core-features)
- [Directory Structure](#directory-structure)
- [Requirements](#requirements)
- [Installation](#installation)
- [Development Workflow](#development-workflow)
- [Make Commands Reference](#make-commands-reference)
- [CLI Scaffolding Tool](#cli-scaffolding-tool)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [Adding New Resources](#adding-new-resources)
- [Database](#database)
- [Docker](#docker)
- [Security](#security)
- [Seeded Accounts](#seeded-accounts-development-only)
- [Roadmap](#roadmap)

---

## System Architecture

Vuelang separates concerns into distinct layers. The Vue 3 frontend operates as a Single Page Application (SPA) communicating asynchronously with the Go API via JSON over HTTP.

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend (SPA)                      │
│          Vue 3  ·  Vite HMR  ·  Tailwind CSS           │
│                  Shadcn Vue Components                  │
└────────────────────────┬────────────────────────────────┘
                         │  JSON over HTTP
┌────────────────────────▼────────────────────────────────┐
│                  Go Backend (Gin)                        │
│                                                         │
│   API Router → HTTP Middleware → Controllers            │
│                                       ↕                 │
│                                   Models / Repos        │
└────────────────────────┬────────────────────────────────┘
                         │  SQL Queries
                ┌────────▼────────┐
                │  MySQL Database  │
                └─────────────────┘
```

---

## Core Features

| Feature | Description |
|---------|-------------|
| 🏗️ **MVC Architecture** | Clean separation across `controllers`, `models`, `middleware`, `services`, and `repositories` |
| 🔐 **JWT Authentication** | Access tokens (15 min) + rotating refresh tokens (7 days) |
| 👮 **RBAC** | Database-backed Role-Based Access Control with middleware enforcement |
| 🔥 **Hot Reloading** | Air for Go backend (`<1s rebuild`) + Vite HMR for Vue (instant) |
| 📦 **Single Binary Deploy** | Vue frontend is embedded directly into the compiled Go binary |
| 🛡️ **Security Headers** | CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy out of the box |
| 🚦 **Rate Limiting** | IP-based token bucket (5 req/min on auth, 60 req/min general) |
| 🐳 **Docker Ready** | Multi-stage Dockerfile + `docker-compose.yml` with MySQL |
| 🧰 **CLI Scaffolding** | `vuelang make:model`, `make:controller`, `make:migration` and more |
| 🧪 **Testing** | `go test ./...` with race detection and coverage |
| 📋 **Audit Logging** | Schema-ready audit log table capturing user_id, action, and IP |
| 🌐 **CORS** | Configurable, origin-allowlisted CORS with no wildcard in production |

---

## Directory Structure

```
vuelang/
├── app/
│   ├── controllers/      # HTTP request handlers, business logic entry points
│   ├── middleware/        # Auth (JWT), RBAC, Rate Limiting interceptors
│   ├── models/            # Data structs (no SQL — pure Go types)
│   ├── repositories/      # SQL queries and database access layer
│   ├── requests/          # Validated request DTOs (struct binding)
│   └── services/          # Business logic orchestration layer
│
├── bootstrap/
│   └── app.go            # Dependency injection wiring
│
├── config/
│   └── app.go            # Typed, validated configuration from env
│
├── database/
│   ├── migrations/        # Schema migration files (numbered, ordered)
│   └── seeders/           # Development data seeders
│
├── internal/
│   ├── framework/         # JWT, hashing, response helpers, rate limiter
│   ├── platform/          # DB connection pool, structured logger
│   └── server/            # Gin engine setup + route registration
│
├── cmd/
│   └── vuelang/           # CLI scaffolding tool source
│
├── ui/                    # Vue 3 frontend (Vite + Tailwind + ShadcnVue)
│   ├── src/
│   │   ├── api/           # Axios client + auth API calls
│   │   ├── router/        # Vue Router configuration
│   │   ├── stores/        # Pinia stores (auth with refresh interceptor)
│   │   └── views/         # Page components (Login, Register, Dashboard…)
│   └── package.json
│
├── main.go               # Application entry point
├── embed.go              # go:embed directive for Vue dist
├── Makefile              # All automation commands
├── Dockerfile            # Multi-stage production Docker build
├── docker-compose.yml    # App + MySQL stack
└── .env.example          # Environment variable template
```

---

## Requirements

| Tool | Minimum Version |
|------|-----------------|
| Go | 1.25.0 |
| Node.js | 18.x LTS or higher |
| npm | 9.x or higher |
| MySQL | 8.0 |

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/milanz247/vuelang.git
cd vuelang
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Open `.env` and fill in your database credentials and generate a JWT secret:

```bash
# Generate a secure JWT secret
openssl rand -base64 64
```

### 3. Install All Dependencies

This single command installs the Air hot-reload tool, all npm packages, and the Vuelang CLI:

```bash
make install
```

### 4. Run Database Migrations

```bash
make migrate
```

### 5. (Optional) Seed Development Data

```bash
make seed
```

---

## Development Workflow

Start the full development environment with concurrent hot-reloading for both frontend and backend:

```bash
make dev
```

```
  ┌──────────────────────────────────────────────────┐
  │              Vuelang V2  DEV                     │
  │   App  →  http://localhost:8080                  │
  │   .go  →  Air rebuilds  (<1s)                    │
  │   .vue →  Vite HMR  (instant)                   │
  └──────────────────────────────────────────────────┘
```

- Go API and Vue frontend are both served at **`http://localhost:8080`**
- Edit any `.go` file → Air detects the change and rebuilds in under 1 second
- Edit any `.vue` file → Vite HMR updates the browser instantly with no full reload

---

## Make Commands Reference

### Development

| Command | Description |
|---------|-------------|
| `make install` | Install Air, npm deps, and the Vuelang CLI |
| `make dev` | Start full dev server (Go + Vue concurrently, hot-reload) |
| `make dev-backend` | Start Go backend only (Air) |
| `make dev-frontend` | Start Vue/Vite dev server only |

### Build & Run

| Command | Description |
|---------|-------------|
| `make build` | Build optimized Vue frontend, embed into Go binary → `dist/vuelang` |
| `make run` | Run the production binary (`dist/vuelang`) |
| `make clean` | Remove `dist/`, `tmp/`, coverage files, and built frontend assets |

### Database

| Command | Description |
|---------|-------------|
| `make migrate` | Run all pending database migrations |
| `make seed` | Seed development data (roles, demo users) |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all tests with race detection and coverage (`go test ./... -v -race -cover`) |
| `make test-cover` | Generate HTML coverage report (`coverage.out`) |

### CLI Shortcuts via Make

| Command | Description |
|---------|-------------|
| `make make-model NAME=Product` | Scaffold a new model |
| `make make-controller NAME=ProductController` | Scaffold a new controller |
| `make make-middleware NAME=AdminOnly` | Scaffold a new middleware |
| `make make-migration NAME=create_products_table` | Scaffold a new migration |

### Docker

| Command | Description |
|---------|-------------|
| `make docker-build` | Build the production Docker image |
| `make docker-up` | Start app + MySQL containers (`docker-compose up -d`) |
| `make docker-down` | Stop all containers |
| `make docker-logs` | Tail logs from the app container |

---

## CLI Scaffolding Tool

After `make install`, the `vuelang` CLI is available globally:

```bash
# Scaffold individual pieces
vuelang make:model       Product
vuelang make:controller  ProductController
vuelang make:middleware  AdminOnly
vuelang make:migration   create_products_table
vuelang make:seeder      ProductSeeder

# Check version
vuelang version
```

---

## API Endpoints

### Public Routes

```
POST  /api/v1/auth/register          Register a new user account
POST  /api/v1/auth/login             Authenticate — returns access + refresh tokens
POST  /api/v1/auth/refresh           Rotate and reissue tokens
POST  /api/v1/auth/forgot-password   Send password reset email
POST  /api/v1/auth/reset-password    Set a new password via reset token

GET   /health                        Health check (used by load balancer / Docker)
```

### Protected Routes

All protected routes require the `Authorization: Bearer <access_token>` header.

```
GET   /api/v1/auth/me                Get current user profile + roles
POST  /api/v1/auth/logout            Invalidate the current refresh token
```

### Admin-Only Routes

```
GET    /api/v1/users                 List all users (paginated)
GET    /api/v1/users/:id             Get a single user by ID
POST   /api/v1/users                 Create a new user
PUT    /api/v1/users/:id             Update an existing user
DELETE /api/v1/users/:id             Delete a user
```

### Response Envelope

All API responses follow this standard JSON envelope:

```json
{
  "success": true,
  "message": "User created",
  "data": { "id": 1, "email": "user@example.com" },
  "errors": null,
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 15
  }
}
```

---

## Environment Variables

Copy `.env.example` to `.env` and configure the following:

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `PORT` | `8080` | No | HTTP listen port |
| `ENV` | `development` | No | `development` or `production` |
| `JWT_SECRET` | — | **Yes (prod)** | Min 32 chars — generate with `openssl rand -base64 64` |
| `JWT_ACCESS_TTL_MINUTES` | `15` | No | Access token lifetime in minutes |
| `JWT_REFRESH_TTL_DAYS` | `7` | No | Refresh token lifetime in days |
| `DB_HOST` | `127.0.0.1` | No | MySQL host |
| `DB_PORT` | `3306` | No | MySQL port |
| `DB_USER` | `root` | No | MySQL username |
| `DB_PASSWORD` | — | **Yes** | MySQL password |
| `DB_NAME` | `vuelang` | No | MySQL database name |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | No | Comma-separated allowed origins |
| `DB_SEED` | `false` | No | Set `true` to automatically seed dev data on startup |

---

## Adding New Resources

Follow these 5 steps to add a fully wired resource (e.g., a `Product`):

### Step 1 — Scaffold Files

```bash
vuelang make:model       Product
vuelang make:controller  ProductController
vuelang make:migration   create_products_table
```

### Step 2 — Define the Schema

Edit the generated migration file in `database/migrations/`:

```go
// database/migrations/YYYYMMDDHHMMSS_create_products_table.go
func Up(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE products (
            id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            name       VARCHAR(255)    NOT NULL,
            price      DECIMAL(10, 2)  NOT NULL DEFAULT 0.00,
            created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
    `)
    return err
}
```

### Step 3 — Register the Migration

Add your migration to `database/migrations/runner.go`:

```go
migrations.Register("YYYYMMDDHHMMSS_create_products_table", Up, Down)
```

### Step 4 — Implement the Repository

Create `app/repositories/product_repository.go` with your SQL queries:

```go
func (r *ProductRepository) FindAll() ([]models.Product, error) {
    rows, err := r.db.Query("SELECT id, name, price FROM products ORDER BY id DESC")
    // ...
}
```

### Step 5 — Wire Routes

Register the controller in `internal/server/router.go`:

```go
productCtrl := controllers.NewProductController(productSvc)

protected.GET("/products",        productCtrl.Index)
protected.GET("/products/:id",    productCtrl.Show)
protected.POST("/products",       productCtrl.Store)
protected.PUT("/products/:id",    productCtrl.Update)
protected.DELETE("/products/:id", productCtrl.Destroy)
```

---

## Database

### Running Migrations

```bash
make migrate
```

### Seeding Development Data

```bash
make seed
```

Alternatively, set `DB_SEED=true` in `.env` to seed automatically on startup.

### Migration File Convention

Migration files in `database/migrations/` are numbered and run in order:

```
0001_create_users_table.go
0002_create_roles_table.go
0003_create_role_user_table.go
0004_create_password_resets_table.go
0005_create_refresh_tokens_table.go
0006_create_audit_logs_table.go
```

---

## Docker

### Start the Full Stack

Launches both the app and a MySQL 8.0 container:

```bash
docker-compose up -d
```

### Required Docker Environment Variables

```bash
export JWT_SECRET=$(openssl rand -base64 64)
export DB_PASSWORD=your_secure_password
docker-compose up -d
```

### Individual Docker Commands

```bash
make docker-build   # Build production image
make docker-up      # Start all containers in background
make docker-down    # Stop and remove containers
make docker-logs    # Tail live logs from the app container
```

The production Docker image uses a **multi-stage build** with a `scratch` base image for the smallest possible attack surface.

---

## Security

Vuelang ships with a comprehensive security posture out of the box:

| Area | Implementation |
|------|---------------|
| **Password Hashing** | bcrypt with cost factor 12 |
| **JWT Tokens** | 15-minute access tokens + 7-day rotating refresh tokens stored in DB |
| **Rate Limiting** | IP-based token bucket — 5 req/min on auth endpoints, 60 req/min general |
| **Security Headers** | CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy |
| **CORS** | Explicit origin allowlist — no wildcard in production |
| **RBAC** | Database-backed roles enforced at the middleware layer |
| **Input Validation** | All request bodies validated via struct binding tags |
| **SQL Safety** | All queries use parameterized statements (`?` placeholders) — no interpolation |
| **Audit Logging** | Schema-ready audit_logs table capturing user_id, action, and IP |
| **Email Privacy** | Login and forgot-password endpoints return identical responses for unknown emails |

Before every production deployment, work through [SECURITY_CHECKLIST.md](SECURITY_CHECKLIST.md).

---

## Seeded Accounts (Development Only)

When `DB_SEED=true` or after running `make seed`:

| Email | Password | Role |
|-------|----------|------|
| `superadmin@vuelang.dev` | `password123` | `super_admin` |
| `admin@vuelang.dev` | `password123` | `admin` |
| `demo@vuelang.dev` | `password123` | `user` |

> ⚠️ **Never use these credentials in production.** Seed data is for local development only.

---

## Roadmap

### V2 → V3 (Next Major Release)

**Core Framework**
- Email system — SMTP mailer with templates and queue integration
- `sqlc` integration — type-safe SQL → Go code generation
- Event system — `event.Dispatch(UserRegistered{})` with in-process listeners
- Redis-backed job queue with retry, delay, and priority
- File storage abstraction (local / S3 / R2)
- WebSocket support for real-time events
- Redis caching layer with TTL helpers
- Task scheduler — cron-like `schedule.Daily(fn)` with distributed lock

**Developer Experience**
- `vuelang new <project>` — scaffold a complete new project
- `vuelang make:resource Product` — generate model + migration + controller + service + repo at once
- `vuelang migrate:rollback` — reversible migrations
- `vuelang migrate:status` — show which migrations have run
- OpenAPI / Swagger auto-generation from route definitions

**Security Additions**
- MFA (TOTP — Google Authenticator compatible)
- Email verification gate on registration
- Account lockout after N failed login attempts
- API key support (`X-API-Key` header)

**Frontend**
- Form validation with vee-validate + zod schema mirroring
- Toast notification system
- Admin panel for user management
- Dark mode with persisted preference

See [ROADMAP.md](ROADMAP.md) for the full plan including V3 → V4 features.

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

<div align="center">

Built with ❤️ by [Milan Madusanka](https://github.com/milanz247)

</div>
