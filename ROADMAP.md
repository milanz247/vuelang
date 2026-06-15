# Vuelang Roadmap

## V1 → V2 (Current Release)

### What's Done ✅
- Real JWT authentication (access + refresh tokens)
- Role-Based Access Control (RBAC)
- Service + Repository architecture layers
- Rate limiting (IP-based token bucket)
- Full security header suite (CSP, X-Frame, etc.)
- CORS with configurable origins
- Password reset flow
- Graceful shutdown (15-second drain)
- Multi-stage Docker build
- CLI scaffolding tool (`vuelang make:model` etc.)
- Standardised JSON response envelope
- bcrypt cost 12 hashing
- Audit log schema
- Typed, validated configuration
- Vue 3 auth pages (Login, Register, ForgotPw, ResetPw, Dashboard)
- Pinia auth store with token refresh interceptor
- `.env` loading via godotenv

---

## V2 → V3 (Next Major Release)

### Core Framework

| Feature | Description | Priority |
|---------|-------------|----------|
| **Email system** | SMTP mailer with templates, queue integration | P0 |
| **sqlc integration** | Type-safe SQL → Go code generation, replaces raw SQL | P0 |
| **Event system** | `event.Dispatch(UserRegistered{})` with in-process listeners | P1 |
| **Queue system** | Redis-backed job queue with retry, delay, priority | P1 |
| **File storage** | `storage.Put(file)` abstraction (local / S3 / R2) | P1 |
| **WebSocket** | Real-time events via gorilla/websocket | P2 |
| **Caching** | Redis cache with TTL (`cache.Remember(key, ttl, fn)`) | P2 |
| **Task scheduler** | Cron-like `schedule.Daily(fn)` with distributed lock | P2 |
| **Localization** | `i18n.T("auth.login_failed")` with JSON translation files | P3 |

### Developer Experience

| Feature | Description |
|---------|-------------|
| `vuelang new <project>` | Scaffold complete new project |
| `vuelang make:resource Product` | Generate model + migration + controller + service + repo at once |
| `vuelang migrate:rollback` | Reversible migrations with `down()` function |
| `vuelang migrate:status` | Show which migrations have run |
| `vuelang serve --port 9000` | Start dev server from CLI |
| Test factories | `factory.Create[User](db)` for seeding tests |
| Hot config reload | Reload env without restart |
| OpenAPI generation | Auto-generate Swagger docs from route definitions |

### Architecture

| Change | Rationale |
|--------|-----------|
| Multi-instance rate limiting (Redis) | In-memory limiter only works with 1 instance |
| Distributed refresh token blacklist | MySQL works; Redis is faster at scale |
| Pagination helpers | `paginator.Paginate(db, query, page, perPage)` |
| Database transactions in service layer | `db.WithTx(ctx, fn)` wrapper |
| Custom error types | Structured errors carry HTTP status code |

### Security

| Feature | Description |
|---------|-------------|
| MFA (TOTP) | Two-factor via TOTP (Google Authenticator compatible) |
| Email verification | Send email on register, gate access until verified |
| Session blacklisting | Invalidate all tokens for a user on password change |
| API keys | `X-API-Key` header support for machine-to-machine auth |
| Account lockout | Lock after N failed login attempts |
| IP allowlisting | Middleware to restrict by IP range |

### Frontend

| Feature | Description |
|---------|-------------|
| Form validation | vee-validate + zod schema mirroring backend rules |
| Toast notifications | Global notification system |
| Admin panel | User management UI |
| Dark mode toggle | Persist preference in localStorage |
| Component library | Full ShadcnVue component set |
| TypeScript strict mode | `strict: true` in tsconfig |

---

## V3 → V4 (Future)

- Plugin system (`vuelang.RegisterPlugin(...)`)
- Multi-tenant support
- GraphQL gateway
- gRPC service integration
- Horizontal auto-scaling with state synchronised via Redis/NATS
- Admin UI generator from struct definitions
- Real-time collaboration features

---

## Critical Issues — Must Fix Before Public V2 Release

| # | Issue | Status |
|---|-------|--------|
| 1 | Add `go.sum` by running `go mod tidy` after cloning | Action needed |
| 2 | Wire audit log writes into controller middleware | Stub in schema |
| 3 | Implement email sending in `ForgotPassword` service | Stub with TODO |
| 4 | Add `HSTS` header at reverse proxy layer (nginx/Caddy) | Infrastructure |
| 5 | Replace in-memory rate limiter with Redis for multi-instance | V3 priority |
| 6 | Add unit tests for auth service and repositories | Before public release |
| 7 | Add `email_verified_at` gate to protected routes | V2 enhancement |
| 8 | CLI: add `vuelang migrate:rollback` | V2 enhancement |
