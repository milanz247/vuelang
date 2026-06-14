<div align="center">

# Vuelang

**A full-stack Go + Vue 3 framework with a Laravel-style structure.**

Single binary in production. Hot-reload in development. Zero runtime dependencies.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Vue](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.5+-3178C6?logo=typescript&logoColor=white)](https://typescriptlang.org)
[![Tailwind](https://img.shields.io/badge/Tailwind-3.4+-06B6D4?logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

---

## What is Vuelang?

Vuelang is a full-stack starter framework combining a **Go (Gin)** backend with a **Vue 3** frontend. The structure mirrors Laravel — dedicated folders for controllers, services, repositories, models, requests, middleware, routes, and database migrations.

In **development**, both sides hot-reload through a single port. In **production**, the entire app compiles into one Go binary with the frontend embedded inside it.

---

## Quick start

```bash
git clone https://github.com/yourusername/vuelang.git
cd vuelang
make install
cp .env.example .env
make dev
```

Open **http://localhost:8080**.

---

## Two commands to know

```bash
make dev      # hot-reload development on http://localhost:8080
make build    # production single binary → dist/vuelang
```

---

## Project structure

```
vuelang/
│
├── main.go                              # Single entry point (dev + prod, env-driven)
├── embed.go                             # //go:embed all:ui/dist
├── Makefile
├── .air.toml                            # Air hot-reload config
├── .env.example
│
├── database/                            # ← Like Laravel's database/
│   ├── migrations/
│   │   ├── runner.go                    # Runs all migrations in order
│   │   ├── 0001_create_users_table.go
│   │   ├── 0002_create_stores_table.go
│   │   └── ...                          # Add new table files here
│   └── seeders/
│       └── user_seeder.go
│
├── internal/
│   ├── app/                             # ← Your application code
│   │   ├── models/                      # Domain structs (User, Product, …)
│   │   │   ├── user.go
│   │   │   └── greeting.go
│   │   │
│   │   ├── repositories/               # Data access layer (interface + implementation)
│   │   │   ├── user_repository.go       # Interface + MySQL implementation
│   │   │   └── greeting_repository.go  # Interface + memory implementation
│   │   │
│   │   ├── services/                   # Business logic layer
│   │   │   ├── user_service.go
│   │   │   └── greeting_service.go
│   │   │
│   │   └── http/
│   │       ├── controllers/            # HTTP handlers (thin — delegate to services)
│   │       │   ├── user_controller.go
│   │       │   └── greeting_controller.go
│   │       ├── requests/               # Input validation structs
│   │       │   └── user_request.go
│   │       └── middleware/             # Gin middleware
│   │           ├── auth.go             # JWT auth (placeholder — add your JWT library)
│   │           └── cors.go
│   │
│   ├── routes/
│   │   └── api.go                      # ← All API routes in one file (like routes/api.php)
│   │
│   ├── config/
│   │   └── config.go                   # Reads env vars; validates secrets in production
│   │
│   ├── platform/
│   │   ├── database/
│   │   │   └── mysql.go                # Connection pool
│   │   └── logger/
│   │       └── logger.go               # Structured slog (text in dev, JSON in prod)
│   │
│   └── server/
│       ├── server.go                   # Server struct, middleware stack, Start/StartDev
│       └── router.go                   # Mounts /api/v1 and delegates to routes.Register
│
└── ui/                                 # Vue 3 frontend
    ├── vite.config.ts
    ├── src/
    │   ├── App.vue
    │   ├── assets/index.css            # Tailwind + CSS variables
    │   └── components/ui/              # shadcn-vue components
    └── dist/                           # Built by make build (gitignored)
```

---

## How the layers connect

The dependency chain is always **one direction**. Nothing in an inner layer knows about an outer one.

```
HTTP Request
     │
     ▼
Controller          (internal/app/http/controllers/)
  • Parse & validate request (using Requests structs)
  • Call service method
  • Write JSON response
  • Knows HTTP — nothing else
     │
     ▼
Service             (internal/app/services/)
  • Business logic, validation rules, password hashing
  • Orchestrates one or more repositories
  • No HTTP, no SQL — pure Go
     │
     ▼
Repository          (internal/app/repositories/)
  • Implements the interface defined in the same file
  • Talks to MySQL (or memory, or Redis)
  • Returns models — never raw SQL rows
     │
     ▼
Model               (internal/app/models/)
  • Plain Go struct
  • JSON tags for serialisation
  • No methods, no logic
```

Dependency injection happens in **one place**: `internal/routes/api.go`. That file creates every repository, service, and controller and wires them together.

---

## Adding a new resource (e.g. Product)

### 1. Migration — `database/migrations/0003_create_products_table.go`

```go
package migrations

import "database/sql"

func CreateProductsTable(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            name       VARCHAR(150)   NOT NULL,
            price      DECIMAL(10,2)  NOT NULL DEFAULT 0.00,
            stock      INT            NOT NULL DEFAULT 0,
            is_active  TINYINT(1)    NOT NULL DEFAULT 1,
            created_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
    `)
    return err
}
```

Register it in `database/migrations/runner.go`:
```go
{"0003_create_products_table", CreateProductsTable},
```

### 2. Model — `internal/app/models/product.go`

```go
package models

type Product struct {
    ID       uint    `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Stock    int     `json:"stock"`
    IsActive bool    `json:"is_active"`
}
```

### 3. Repository — `internal/app/repositories/product_repository.go`

```go
package repositories

import (
    "context"
    "database/sql"
    "go-cloud-erp/internal/app/models"
)

type ProductRepository interface {
    FindAll(ctx context.Context) ([]models.Product, error)
    FindByID(ctx context.Context, id uint) (*models.Product, error)
    Create(ctx context.Context, p *models.Product) (*models.Product, error)
    Delete(ctx context.Context, id uint) error
}

type mysqlProductRepository struct{ db *sql.DB }

func NewMySQLProductRepository(db *sql.DB) ProductRepository {
    return &mysqlProductRepository{db: db}
}

func (r *mysqlProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT id, name, price, stock, is_active FROM products`)
    // ... scan rows
}
```

### 4. Request — `internal/app/http/requests/product_request.go`

```go
package requests

type CreateProductRequest struct {
    Name  string  `json:"name"  binding:"required,min=2"`
    Price float64 `json:"price" binding:"required,min=0"`
    Stock int     `json:"stock" binding:"min=0"`
}
```

### 5. Service — `internal/app/services/product_service.go`

```go
package services

type ProductService struct {
    repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

func (s *ProductService) List(ctx context.Context) ([]models.Product, error) {
    return s.repo.FindAll(ctx)
}
```

### 6. Controller — `internal/app/http/controllers/product_controller.go`

```go
package controllers

type ProductController struct{ service *services.ProductService }

func NewProductController(s *services.ProductService) *ProductController {
    return &ProductController{service: s}
}

func (ctrl *ProductController) Index(c *gin.Context) {
    products, err := ctrl.service.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": products})
}
```

### 7. Wire it — `internal/routes/api.go`

```go
// inside Register(), in the protected group:
productRepo       := repositories.NewMySQLProductRepository(db)
productService    := services.NewProductService(productRepo)
productController := controllers.NewProductController(productService)

protected.GET("/products",        productController.Index)
protected.GET("/products/:id",    productController.Show)
protected.POST("/products",       productController.Store)
protected.PUT("/products/:id",    productController.Update)
protected.DELETE("/products/:id", productController.Destroy)
```

That's the complete pattern. Every resource follows the same seven steps.

---

## Migrations

Migrations live in `database/migrations/`. Each file contains one function. Add new tables by creating a new numbered file and registering it in `runner.go`.

```bash
# Migrations run automatically on startup
make dev
make run
```

Each migration uses `CREATE TABLE IF NOT EXISTS` so it is safe to re-run every time.

---

## API routes

All routes are defined in **one file**: `internal/routes/api.go`.

```
GET  /api/v1/greeting          public
GET  /api/v1/users             protected (Bearer token)
GET  /api/v1/users/:id         protected
POST /api/v1/users             protected
PUT  /api/v1/users/:id         protected
DELETE /api/v1/users/:id       protected
```

---

## Frontend (Vue 3)

API calls always use `/api/v1` — same URL in dev and production.

```ts
// ui/src/services/api.ts
const BASE = '/api/v1'

export async function getProducts() {
    const res = await fetch(`${BASE}/products`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return res.json()
}
```

Add shadcn-vue components:
```bash
cd ui
npx shadcn-vue@latest add button
npx shadcn-vue@latest add dialog
npx shadcn-vue@latest add table
npx shadcn-vue@latest add input
```

---

## Configuration

```env
# .env
PORT=8080
ENV=development           # "production" → serves embedded ui/dist

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=vuelang_db

# Generate with: openssl rand -base64 48
JWT_SECRET=change-me-before-deploying
```

In production the server refuses to start if `JWT_SECRET` is empty or still set to the default value.

---

## Security built in

| Header | Value |
|---|---|
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `X-XSS-Protection` | `1; mode=block` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Cache-Control` (API routes) | `no-store` |
| CORS | Configurable per-environment |
| JWT validation | Middleware placeholder in `middleware/auth.go` |
| Secret validation | Server refuses to start with default JWT secret in production |

---

## Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.26, Gin |
| Database | MySQL (optional) |
| Frontend | Vue 3, TypeScript, Vite 5 |
| UI components | shadcn-vue, Tailwind CSS 3 |
| Go hot-reload | Air |
| Frontend hot-reload | Vite HMR |
| Password hashing | bcrypt |

---

## License

MIT — see [LICENSE](LICENSE).
