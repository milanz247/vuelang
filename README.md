# Vuelang V2

> **Enterprise-grade Go + Vue 3 full-stack framework**  
> Single binary. JWT auth. RBAC. Production secure.

```
  ┌────────────────────────────────────────────────────────┐
  │  Backend  →  Go 1.23 + Gin                            │
  │  Frontend →  Vue 3 + Vite + Tailwind + ShadcnVue      │
  │  Database →  MySQL 8.0                                 │
  │  Auth     →  JWT (access + refresh) + RBAC            │
  │  Deploy   →  Single binary or Docker                   │
  └────────────────────────────────────────────────────────┘
```

---

## Quick Start

```bash
# 1. Clone & install
git clone https://github.com/your-org/vuelang.git
cd vuelang
cp .env.example .env   # fill in DB credentials + JWT_SECRET

# 2. Install tools
make install           # installs Air, npm deps, CLI

# 3. Start dev server (hot reload for both Go and Vue)
make dev
#   → http://localhost:8080
```

## Production Build

```bash
make build             # builds Vue → embeds into single Go binary
make run               # starts production binary on :8080
```

## Docker

```bash
docker-compose up -d   # starts app + MySQL
```

---

## CLI

```bash
vuelang make:model       Product
vuelang make:controller  ProductController
vuelang make:middleware  AdminOnly
vuelang make:migration   create_products_table
vuelang make:seeder      ProductSeeder
vuelang version
```

---

## Adding a New Resource (5 Steps)

```bash
# 1. Scaffold
vuelang make:model       Product
vuelang make:controller  ProductController
vuelang make:migration   create_products_table

# 2. Edit database/migrations/YYYYMMDDHHMMSS_create_products_table.go
#    Add your CREATE TABLE statement

# 3. Register migration in database/migrations/runner.go

# 4. Create app/repositories/product_repository.go
#    Implement DB queries

# 5. Wire in internal/server/router.go
productCtrl := controllers.NewProductController(productSvc)
protected.GET("/products", productCtrl.Index)
```

---

## API Endpoints

### Public
```
POST /api/v1/auth/register         Register new account
POST /api/v1/auth/login            Get access + refresh tokens
POST /api/v1/auth/refresh          Rotate tokens
POST /api/v1/auth/forgot-password  Request reset email
POST /api/v1/auth/reset-password   Set new password

GET  /health                       Health check
```

### Protected (requires `Authorization: Bearer <access_token>`)
```
GET  /api/v1/auth/me               Current user + roles
POST /api/v1/auth/logout           Invalidate refresh token

# Admin only
GET    /api/v1/users               List users
GET    /api/v1/users/:id           Get user
POST   /api/v1/users               Create user
PUT    /api/v1/users/:id           Update user
DELETE /api/v1/users/:id           Delete user
```

### Response Format
```json
{
  "success": true,
  "message": "User created",
  "data": { ... },
  "errors": null,
  "meta": { "total": 100, "page": 1, "per_page": 15 }
}
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `ENV` | `development` | `development` or `production` |
| `JWT_SECRET` | *(required in prod)* | `openssl rand -base64 64` |
| `JWT_ACCESS_TTL_MINUTES` | `15` | Access token lifetime |
| `JWT_REFRESH_TTL_DAYS` | `7` | Refresh token lifetime |
| `DB_HOST` | `127.0.0.1` | MySQL host |
| `DB_NAME` | `vuelang` | Database name |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Comma-separated origins |
| `DB_SEED` | `false` | Set `true` to seed dev data |

---

## Project Structure

```
app/controllers/    HTTP handlers
app/middleware/     Auth, RBAC, rate limiting
app/models/         Data structs (no SQL)
app/repositories/   SQL queries (data access layer)
app/requests/       Validated request DTOs
app/services/       Business logic

bootstrap/app.go    Dependency injection wiring
config/app.go       Typed configuration

database/migrations/  Schema migration files
database/seeders/     Dev data seeders

internal/framework/   JWT, hash, response, rate limit
internal/platform/    DB connection, logger
internal/server/      Gin setup + route registration

cmd/vuelang/        CLI scaffolding tool
ui/                 Vue 3 frontend
```

---

## Security

See [SECURITY_CHECKLIST.md](SECURITY_CHECKLIST.md) before every production deployment.  
See [AUDIT_REPORT.md](AUDIT_REPORT.md) for the full security audit findings.

**Key security features:**
- JWT with 15-minute access tokens + 7-day rotating refresh tokens
- bcrypt cost 12 password hashing
- IP-based rate limiting (5 req/min on auth endpoints)
- CORS with explicit origin allowlist
- Full security header suite (CSP, X-Frame-Options, etc.)
- RBAC with DB-backed roles
- Audit log schema ready

---

## Seeded Accounts (development only)

When `DB_SEED=true`:

| Email | Password | Role |
|-------|----------|------|
| superadmin@vuelang.dev | password123 | super_admin |
| admin@vuelang.dev | password123 | admin |
| demo@vuelang.dev | password123 | user |

---

## License

MIT
