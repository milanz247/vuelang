# Vuelang V2 — Security Checklist

Use this before every production deployment.

## Authentication
- [ ] `JWT_SECRET` is at least 32 characters, randomly generated (`openssl rand -base64 64`)
- [ ] `JWT_SECRET` is NOT committed to source control
- [ ] Access tokens expire in ≤15 minutes
- [ ] Refresh tokens expire in ≤30 days
- [ ] Refresh tokens are stored in DB and invalidated on logout
- [ ] Password reset tokens expire in ≤1 hour
- [ ] bcrypt cost is ≥12
- [ ] Login endpoint does NOT reveal whether email exists (same response for wrong email vs wrong password)
- [ ] Forgot password endpoint does NOT reveal whether email exists

## HTTP Security Headers
- [ ] `Content-Security-Policy` is set
- [ ] `X-Frame-Options: DENY` is set
- [ ] `X-Content-Type-Options: nosniff` is set
- [ ] `Referrer-Policy: strict-origin-when-cross-origin` is set
- [ ] `Permissions-Policy` restricts camera/mic/geolocation
- [ ] HSTS is set at reverse proxy level (`Strict-Transport-Security: max-age=31536000; includeSubDomains`)
- [ ] API responses include `Cache-Control: no-store`

## CORS
- [ ] `CORS_ALLOWED_ORIGINS` does NOT contain `*` in production
- [ ] Allowed origins are the exact production domain(s)

## Rate Limiting
- [ ] Auth endpoints (login, register, forgot-password) are rate-limited (≤10 req/min per IP)
- [ ] General API is rate-limited (≤60 req/min per IP)
- [ ] Rate limit uses real client IP (check reverse proxy X-Forwarded-For)

## Input Validation
- [ ] All request bodies are validated via struct binding tags
- [ ] Email max length enforced (prevent oversized inputs)
- [ ] Password max length enforced (bcrypt input limit is 72 bytes)
- [ ] All SQL queries use parameterised statements (`?` placeholders)
- [ ] No user input is interpolated into SQL strings

## RBAC
- [ ] Every admin route has `RequireRole("admin", "super_admin")` middleware
- [ ] Default registered users get the `user` role (not `admin`)
- [ ] Role assignments require admin privileges

## Secrets Management
- [ ] `.env` is in `.gitignore`
- [ ] No secrets in `docker-compose.yml` (use `${VAR}` references)
- [ ] No secrets in Dockerfile
- [ ] DB password is NOT the default (`root`/`root`)
- [ ] DB user has minimal privileges (SELECT, INSERT, UPDATE, DELETE — no DROP, CREATE USER)

## Logging
- [ ] Production logs are JSON format
- [ ] Error logs do NOT include stack traces visible to clients
- [ ] Passwords and tokens are NEVER logged
- [ ] Audit logs capture write operations with user_id and IP

## Infrastructure
- [ ] HTTPS is enforced (TLS terminates at load balancer or Caddy)
- [ ] Database port (3306) is NOT exposed to the internet
- [ ] Docker containers run as non-root user
- [ ] `ENV=production` is set in the runtime environment
- [ ] Health check endpoint is accessible for load balancer probes
- [ ] Graceful shutdown timeout (15s) is greater than max request duration

## Database
- [ ] DB credentials are rotated from defaults
- [ ] Connection pool is tuned for expected load
- [ ] `parseTime=true` and `loc=UTC` are in DSN
- [ ] Foreign keys are enforced (InnoDB with FOREIGN KEY constraints)
- [ ] Backups are configured and tested

## Go Binary
- [ ] Built with `-ldflags="-s -w"` (strip debug info)
- [ ] `CGO_ENABLED=0` for static binary
- [ ] Binary is `scratch`-based Docker image (minimal attack surface)
