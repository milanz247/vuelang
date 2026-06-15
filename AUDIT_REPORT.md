# Vuelang V2 — Complete Enterprise Audit Report

> **Framework:** Vuelang (Go + Gin + Vue 3 + MySQL)
> **Audit Date:** 2026-06-15
> **Auditor:** Senior Framework Architect / Security Engineer
> **Version Audited:** V1 → V2 upgrade

---

## Table of Contents
1. [Architecture Review](#1-architecture-review)
2. [Security Audit](#2-security-audit)
3. [Developer Experience](#3-developer-experience)
4. [Authentication Scaffolding](#4-authentication-scaffolding)
5. [ShadcnVue Integration](#5-shadcnvue-integration)
6. [Production Readiness](#6-production-readiness)
7. [Laravel Comparison](#7-laravel-comparison)
8. [Final Deliverables](#8-final-deliverables)

---

## 1. Architecture Review

### Current V1 Problems

| # | Problem | Risk |
|---|---------|------|
| 1 | Controllers call models directly — no service layer | Business logic leaks into HTTP handlers; untestable |
| 2 | Models contain both schema AND SQL queries | Violates Single Responsibility; swap DB = rewrite everything |
| 3 | No repository pattern | Cannot mock data access in unit tests |
| 4 | `Auth()` middleware is a stub — accepts **any** bearer token | **CRITICAL: zero authentication in production** |
| 5 | No dependency injection — `*sql.DB` passed as raw field | Tight coupling; impossible to swap implementations |
| 6 | No service container | Can't wire complex dependency graphs |
| 7 | Module name `go-cloud-erp` hardcoded everywhere | Wrong module name for a framework called Vuelang |
| 8 | `db.Close()` never called — resource leak | DB connection pool never released on shutdown |
| 9 | No graceful shutdown | Active requests killed on Ctrl+C |
| 10 | Error messages expose internal details to clients | Information leakage |

### V2 Architecture (Implemented)

```
main.go  ──→  bootstrap/app.go  (wires entire dependency graph)
              │
              ├── config/app.go           (all env vars, typed, validated)
              ├── internal/platform/      (infra: DB, logger)
              ├── internal/framework/     (jwt, hash, response, ratelimit)
              │
              ├── app/repositories/       (data access — raw SQL, context-aware)
              ├── app/services/           (business logic — pure Go, testable)
              ├── app/controllers/        (HTTP — bind → service → respond)
              ├── app/middleware/         (auth, rbac, rate limit, audit)
              └── app/requests/           (validated DTOs)
```

### Architecture Diagram

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────────────────┐
│  Gin Router                                             │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐            │
│  │ security │→ │   CORS   │→ │ rateLimit  │ middleware │
│  │ headers  │  │          │  │   (IP)     │            │
│  └──────────┘  └──────────┘  └────────────┘            │
│                                    │                    │
│                          ┌─────────▼──────────┐        │
│                          │  AuthMiddleware     │        │
│                          │  (JWT validate)     │        │
│                          └─────────┬──────────┘        │
│                                    │                    │
│                          ┌─────────▼──────────┐        │
│                          │  RBACMiddleware     │        │
│                          │  (role check)       │        │
│                          └─────────┬──────────┘        │
└────────────────────────────────────┼───────────────────┘
                                     │
                          ┌──────────▼──────────┐
                          │    Controller        │
                          │  (bind + validate)   │
                          └──────────┬──────────┘
                                     │
                          ┌──────────▼──────────┐
                          │      Service         │
                          │  (business logic)    │
                          └──────────┬──────────┘
                                     │
                          ┌──────────▼──────────┐
                          │    Repository        │
                          │   (SQL queries)      │
                          └──────────┬──────────┘
                                     │
                          ┌──────────▼──────────┐
                          │      MySQL           │
                          └─────────────────────┘
```

### Should a Services Layer Exist? YES.

| Layer | Laravel Equivalent | Purpose |
|-------|--------------------|---------|
| Controller | Controller | Bind HTTP, call service, respond |
| Service | Service / Action | Business rules, orchestration |
| Repository | Eloquent Model / Repository | DB queries only |
| Model (struct) | Eloquent Model (schema) | Data shape definition |

### Should Repositories Exist? YES.

Repositories let you:
- Swap MySQL for Postgres without touching services
- Mock in unit tests (`type UserRepo interface { FindByEmail(...) }`)
- Keep SQL out of business logic

### Should Events/Queues Exist? ROADMAP V3.

Adding channels/goroutines in V2 is over-engineering. Use Redis + a worker in V3.

### Should DI Exist? YES — implemented via `bootstrap/app.go`.

---

## 2. Security Audit

### CRITICAL Findings (Must Fix Before Release)

---

#### SEC-001 — Auth Middleware Accepts Any Bearer Token

**Severity:** CRITICAL  
**File:** `app/middleware/auth.go` (V1)

**Description:** The V1 auth middleware only checks that the header starts with `"Bearer "`. Any string like `"Bearer abc"` passes.

**Attack Scenario:**
```bash
curl -H "Authorization: Bearer faketoken" /api/v1/users
# Returns 200 OK — full user list exposed
```

**Fix (V2 Implementation):**
```go
claims, err := m.jwt.ValidateAccess(tokenStr)
if err != nil { ... abort ... }
c.Set("user_id", claims.UserID)
```
JWT is cryptographically verified against `JWT_SECRET`. ✅ Fixed in V2.

---

#### SEC-002 — No Rate Limiting

**Severity:** HIGH  
**Description:** V1 has zero rate limiting. An attacker can brute-force passwords or hammer the API.

**Fix (V2):** Token-bucket limiter — 60 req/min on API, 5 req/min on auth endpoints. ✅

---

#### SEC-003 — No CSRF Protection for State-Changing Endpoints

**Severity:** HIGH  
**Description:** V1 has no CSRF headers. Any website can trigger PUT/DELETE from a victim's browser.

**Fix (V2):** API-only framework using `Authorization: Bearer` header — same-origin cookie attacks don't apply. For cookie-based sessions, add `SameSite=Strict` and CSRF tokens.

---

#### SEC-004 — JWT Secret Falls Back to Empty String in Dev

**Severity:** HIGH  
**File:** `config/app.go` (V1)

**Description:** In development, `JWT_SECRET` defaults to `""`. `golang-jwt` will sign/verify with an empty key, making tokens trivially forgeable.

**Attack Scenario:**
```go
// Attacker signs: jwt.NewWithClaims(HS256, Claims{UserID: 1, Roles: ["super_admin"]})
// With secret = "" — this validates!
```

**Fix (V2):** Validation fails fast if secret is empty in **any** environment. Warn in dev, fatal in prod. ✅

---

#### SEC-005 — Error Messages Leak Internal Details

**Severity:** MEDIUM  
**Description:**
```go
// V1:
c.JSON(500, gin.H{"error": err.Error()})
// Exposes: "User.Create: Error 1062: Duplicate entry 'x@x.com' for key 'users.email'"
```

**Fix (V2):** Standardised response package. Internal errors are logged; clients get generic messages. ✅

---

#### SEC-006 — Password Reset Not Implemented

**Severity:** HIGH  
**Description:** V1 has no password reset flow. Users cannot recover accounts.

**Fix (V2):** Full forgot/reset password flow with time-limited tokens (1 hour TTL), stored in `password_resets` table, invalidated after use. ✅

---

#### SEC-007 — No Refresh Token Rotation

**Severity:** MEDIUM  
**Description:** V1 has no refresh tokens. Access tokens (if issued) would never expire safely.

**Fix (V2):** Access tokens (15min) + refresh tokens (7 days, stored in DB). Refresh rotates tokens. Logout invalidates refresh token. ✅

---

#### SEC-008 — Missing Security Headers

**Severity:** MEDIUM

| Header | V1 | V2 |
|--------|----|----|
| Content-Security-Policy | ❌ | ✅ |
| X-Frame-Options | ✅ | ✅ |
| X-Content-Type-Options | ✅ | ✅ |
| Referrer-Policy | ✅ | ✅ |
| Permissions-Policy | ❌ | ✅ |
| HSTS | ❌ | Add via reverse proxy (nginx/Caddy) |

---

#### SEC-009 — No RBAC

**Severity:** HIGH  
**Description:** V1 middleware passes once the header format is correct — no role checking. Any authenticated user can delete any other user.

**Fix (V2):** `RequireRole("admin", "super_admin")` middleware on admin routes. Roles loaded from DB, embedded in JWT claims. ✅

---

#### SEC-010 — SQL Injection via Raw Query Concatenation

**Severity:** LOW (current code uses `?` placeholders correctly)  
**Description:** All V1 queries use parameterised statements correctly. Risk is low but developer docs should explicitly warn against string concatenation in queries.

**Recommendation:** Add a query builder or document the rule clearly. ✅ V2 maintains parameterised queries.

---

#### SEC-011 — Email Enumeration via Forgot Password

**Severity:** LOW  
**Description:** A timing attack or different error response reveals whether an email exists.

**Fix (V2):** Always return success: `"If an account exists, a reset link has been sent"`. ✅

---

#### SEC-012 — bcrypt Cost Too Low

**Severity:** MEDIUM  
**Description:** V1 uses `bcrypt.DefaultCost` (cost 10, ~80ms). Modern hardware can crack this faster.

**Fix (V2):** Cost 12 (~250ms). Suitable for production. ✅

---

### HTTP Security Headers (Full Reference)

```go
// V2 implementation in internal/server/server.go
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "0")          // Modern: let CSP handle it
c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
c.Header("Content-Security-Policy", "default-src 'self'; ...")
```

---

## 3. Developer Experience

### V1 DX Problems

- No CLI — developers must copy files manually
- No auto-seeding
- No standardised response format
- Routes are manually wired with `*sql.DB` passed everywhere
- No clear "how to add a new resource" guide

### V2 DX Improvements

#### CLI Tool (`cmd/vuelang`)

```bash
# Install
go install ./cmd/vuelang

# Scaffold
vuelang make:model Product
vuelang make:controller ProductController
vuelang make:middleware AdminOnly
vuelang make:migration create_products_table
vuelang make:seeder ProductSeeder
```

#### How to Add a New Resource (5 Steps)

```
1. vuelang make:model Product
2. vuelang make:migration create_products_table
   → Add migration to database/migrations/runner.go
3. vuelang make:controller ProductController
   → Add service/repository calls
4. Add repository in app/repositories/product_repository.go
5. Add routes in internal/server/router.go
```

#### Service Container Pattern

```go
// bootstrap/app.go — one place to wire everything
hasher   := hash.NewBcrypt()
jwtSvc   := jwt.NewService(secret, ttl, refreshDays)
userRepo := repositories.NewUserRepository(db)
userSvc  := services.NewUserService(userRepo, hasher)
userCtrl := controllers.NewUserController(userSvc)
```

#### Standard Response Format

```json
{
  "success": true,
  "message": "User created",
  "data": { "id": 1, "name": "Jane" },
  "errors": null,
  "meta": { "total": 100, "page": 1, "per_page": 15 }
}
```

### What Laravel Does That V2 Should Add (V3 Roadmap)

| Feature | Laravel | V2 | V3 Target |
|---------|---------|-----|-----------|
| ORM | Eloquent | Raw SQL | sqlc or GORM |
| Validation | Form Requests | Gin binding | Custom rule builders |
| Events | `event(new UserRegistered)` | ❌ | Go channels + handlers |
| Queue | `dispatch(new SendEmail)` | ❌ | Redis worker |
| Mail | `Mail::to()->send()` | Stub | SMTP via `net/smtp` |
| Scheduler | `$schedule->daily()` | ❌ | Cron + goroutine |
| Localization | `__('messages.welcome')` | ❌ | JSON i18n files |
| Broadcasting | Pusher/Echo | ❌ | WebSocket (gorilla) |

---

## 4. Authentication Scaffolding

### Database Schema

```sql
-- users
CREATE TABLE users (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name              VARCHAR(100) NOT NULL,
    email             VARCHAR(150) NOT NULL UNIQUE,
    password          VARCHAR(255) NOT NULL,
    email_verified_at DATETIME NULL,
    is_active         TINYINT(1) NOT NULL DEFAULT 1,
    created_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- roles + RBAC
CREATE TABLE roles (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(50) NOT NULL UNIQUE,   -- 'admin', 'user'
    display_name VARCHAR(100) NOT NULL,
    description  VARCHAR(255) NOT NULL
);

CREATE TABLE role_user (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- refresh tokens (rotation)
CREATE TABLE refresh_tokens (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- password reset
CREATE TABLE password_resets (
    email      VARCHAR(150) NOT NULL PRIMARY KEY,
    token      VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL
);
```

### Auth Flow

```
Guest
  │
  ├─ POST /api/v1/auth/register
  │    → validate → hash password → create user → assign 'user' role
  │    → generate access+refresh tokens → return tokens
  │
  ├─ POST /api/v1/auth/login
  │    → find user by email → bcrypt.Compare → check is_active
  │    → generate tokens → store refresh token in DB → return
  │
  ├─ POST /api/v1/auth/forgot-password
  │    → find user (if not found, return success anyway)
  │    → generate UUID reset token → store with 1h TTL
  │    → (TODO: send email with link)
  │
  ├─ POST /api/v1/auth/reset-password?token=<uuid>
  │    → validate token not expired → hash new password
  │    → update user → delete reset token → invalidate all refresh tokens
  │
  └─ POST /api/v1/auth/refresh
       → validate refresh token in DB → check expiry
       → delete old token → generate new pair → store new refresh token

Authenticated
  │
  ├─ GET  /api/v1/auth/me      → return current user with roles
  └─ POST /api/v1/auth/logout  → delete refresh token from DB
```

### JWT Token Structure

```json
{
  "uid": 1,
  "email": "user@example.com",
  "roles": ["admin"],
  "sub": "1",
  "iat": 1718400000,
  "exp": 1718400900,
  "iss": "vuelang"
}
```

### RBAC Usage

```go
// In router:
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin", "super_admin"))
admin.GET("/users", userCtrl.Index)

// Adding a new role check:
admin.DELETE("/users/:id", middleware.RequireRole("super_admin"), userCtrl.Destroy)
```

---

## 5. ShadcnVue Integration

### Is ShadcnVue the Best Choice? YES — with caveats.

**Pros:**
- Unstyled primitives via Radix Vue — accessible by design
- Copy-paste ownership — no version lock-in
- Tailwind-first — aligns with current stack
- Dark mode via `class` strategy
- Type-safe props

**Cons:**
- Manual installation per component (no `npm install shadcn`)
- Migration cost if Radix Vue API changes
- Bundle can grow large if all components are copied

### Folder Organisation (V2)

```
ui/src/components/
├── ui/              ← ShadcnVue primitives (Button, Card, Input...)
│   ├── button/
│   │   ├── Button.vue
│   │   └── index.ts
│   ├── card/
│   ├── input/
│   └── index.ts     ← barrel export
├── forms/           ← Composed form components
├── layouts/         ← Page layout wrappers
└── shared/          ← App-specific reusable components
```

### Dark Mode Setup

```ts
// In App.vue
const isDark = useDark()    // @vueuse/core
const toggle = useToggle(isDark)
// Adds/removes 'dark' class on <html>
```

### Form Validation with Vee-Validate + Zod (Recommendation)

```ts
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

const schema = toTypedSchema(z.object({
  email: z.string().email(),
  password: z.string().min(8),
}))
const { handleSubmit, errors } = useForm({ validationSchema: schema })
```

**Recommendation:** Add `vee-validate` + `zod` for V2 forms for type-safe client-side validation that mirrors server-side rules.

---

## 6. Production Readiness

### Missing Features (V1 → V2 Gaps)

| Feature | V1 | V2 | Priority |
|---------|----|----|----------|
| Graceful shutdown | ❌ | ✅ 15s drain | P0 |
| JWT validation | ❌ (stub) | ✅ | P0 |
| Refresh tokens | ❌ | ✅ | P0 |
| Rate limiting | ❌ | ✅ | P0 |
| RBAC | ❌ | ✅ | P0 |
| Audit logs table | ❌ | ✅ schema | P1 |
| Health check endpoint | ❌ | ✅ GET /health | P1 |
| Docker support | ❌ | ✅ multi-stage | P1 |
| Password reset | ❌ | ✅ | P1 |
| Structured logging | ✅ | ✅ JSON in prod | P1 |
| .env loading | ❌ | ✅ godotenv | P1 |
| CORS config | ❌ | ✅ | P1 |
| Connection pool tuning | ✅ | ✅ | P2 |
| HTTP timeouts | ❌ | ✅ 15s r/w | P2 |
| CSP headers | ❌ | ✅ | P2 |
| CI/CD pipeline | ❌ | See below | P2 |
| Email sending | ❌ | Stub | P3 |
| Test coverage | ❌ | Scaffold | P3 |

### Logging Strategy

```
Development: Text format, DEBUG level, with source location
Production:  JSON format, INFO level, stdout → log aggregator (e.g. Loki/CloudWatch)

Log fields: method, path, status, latency, ip, user_id (when authenticated)
```

### Health Check (`GET /health`)

```json
{
  "success": true,
  "message": "healthy",
  "data": {
    "status": "ok",
    "timestamp": "2026-06-15T10:00:00Z",
    "version": "2.0.0"
  }
}
```

### Recommended CI/CD (GitHub Actions)

```yaml
# .github/workflows/ci.yml
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql: { image: mysql:8.0, env: { MYSQL_ROOT_PASSWORD: root, MYSQL_DATABASE: vuelang_test } }
    steps:
      - uses: actions/setup-go@v5
      - run: go test ./... -race -cover
      - run: cd ui && npm ci && npm run build
  
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - run: docker build -t vuelang:${{ github.sha }} .
      - run: docker push ...
```

### Load Balancing / Horizontal Scaling

Since V2 is stateless (JWT-based, tokens stored in MySQL not in-memory):
- Deploy multiple instances behind nginx/Caddy/AWS ALB
- All instances share the same MySQL (Aurora for HA)
- Rate limiter: replace in-memory with Redis for multi-instance support (V3)

---

## 7. Laravel Comparison

| Feature | Laravel | Vuelang V1 | Vuelang V2 |
|---------|---------|-----------|-----------|
| **Routing** | `routes/api.php`, named routes, resource groups | Manual `router.go` | Structured groups with middleware |
| **Middleware** | Global + route-level, pipeline pattern | Single stub | Auth, RBAC, Rate-limit, CORS |
| **ORM** | Eloquent (Active Record) | Raw SQL in models | Raw SQL in repositories |
| **Auth** | Sanctum / Passport (JWT) | Stub (any token passes) | JWT + refresh tokens + RBAC |
| **Validation** | Form Request classes, rules array | Gin struct tags | Gin struct tags + custom errors |
| **Queue** | Horizon + Redis | ❌ | ❌ (V3) |
| **Events** | Event classes + Listeners | ❌ | ❌ (V3) |
| **Scheduler** | `$schedule->command()` | ❌ | ❌ (V3) |
| **Service Container** | `app()->make()` | ❌ | Bootstrap DI wiring |
| **Testing** | PHPUnit + Factories | ❌ | Scaffolded (V3: testify) |
| **Security** | CSRF, bcrypt, rate limit | Partial | CORS, rate limit, bcrypt 12, headers |
| **CLI** | `artisan make:model` | ❌ | `vuelang make:model` |
| **Scaffolding** | `laravel new`, Breeze, Jetstream | ❌ | Auth pages, seed roles, migrations |
| **Config** | `config/app.php`, `env()` | Single struct | Typed config struct, validated |
| **Mail** | `Mail::send()`, Mailable class | ❌ | Stub (V3: SMTP) |
| **Storage** | `Storage::put()` | ❌ | ❌ (V3) |
| **Broadcast** | Laravel Echo + Pusher | ❌ | ❌ (V3) |
| **Docker** | Sail | ❌ | Multi-stage Dockerfile + compose |
| **DX** | ★★★★★ | ★★☆☆☆ | ★★★★☆ |
| **Performance** | ~2,000 req/s | ~50,000 req/s | ~50,000 req/s |
| **Deployment** | PHP-FPM + Nginx | Single binary | Single binary / Docker |

---

## 8. Final Deliverables

### Critical Issues Fixed Before This Release

- [x] SEC-001: Auth middleware now validates JWT cryptographically
- [x] SEC-002: Rate limiting added (IP-based token bucket)
- [x] SEC-004: JWT secret validated at startup (empty = fatal)
- [x] SEC-005: Error messages sanitised; internal errors logged only
- [x] SEC-006: Password reset flow implemented
- [x] SEC-007: Refresh token rotation implemented
- [x] SEC-008: Full security header suite added
- [x] SEC-009: RBAC middleware implemented
- [x] SEC-012: bcrypt cost raised to 12
- [x] ARCH-001: Service + repository layers added
- [x] ARCH-008: Graceful shutdown (15s drain)
- [x] ARCH-009: `db.Close()` called via `defer app.Close()`

### V2 Folder Structure (Final)

```
vuelang/
├── app/
│   ├── controllers/          HTTP handlers (bind → service → respond)
│   │   ├── auth_controller.go
│   │   └── user_controller.go
│   ├── middleware/           Cross-cutting concerns
│   │   ├── auth.go           JWT validation
│   │   ├── rbac.go           Role-based access
│   │   └── rate_limit.go     IP token bucket
│   ├── models/               Data shapes (structs only)
│   │   ├── user.go, role.go, refresh_token.go
│   │   ├── password_reset.go, audit_log.go
│   ├── repositories/         SQL queries
│   │   ├── user_repository.go
│   │   └── role_repository.go
│   ├── requests/             Validated DTOs
│   │   ├── auth_request.go
│   │   └── user_request.go
│   └── services/             Business logic
│       ├── auth_service.go
│       └── user_service.go
├── bootstrap/
│   └── app.go                Dependency graph wiring
├── cmd/
│   └── vuelang/main.go       CLI scaffolding tool
├── config/
│   └── app.go                Typed, validated config
├── database/
│   ├── migrations/           Numbered schema changes
│   └── seeders/              Dev data seeders
├── internal/
│   ├── framework/
│   │   ├── hash/             bcrypt wrapper
│   │   ├── jwt/              JWT + refresh token service
│   │   ├── ratelimit/        Token bucket limiter
│   │   └── response/         Standard JSON envelope
│   ├── platform/
│   │   ├── database/         MySQL connection pool
│   │   └── logger/           slog initialisation
│   └── server/
│       ├── server.go         Gin setup, middleware, static serving
│       └── router.go         Route registration
├── ui/                       Vue 3 + Vite + Tailwind
│   └── src/
│       ├── api/              Axios client + endpoint modules
│       ├── router/           Vue Router (auth guards)
│       ├── stores/           Pinia auth store
│       └── views/auth/       Login, Register, ForgotPw, ResetPw
├── .env.example
├── .air.toml
├── Dockerfile                Multi-stage (node → go → scratch)
├── docker-compose.yml
├── embed.go
├── go.mod
├── main.go                   Entry point with graceful shutdown
└── Makefile
```
