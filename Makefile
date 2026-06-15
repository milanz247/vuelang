.PHONY: install dev dev-backend dev-frontend build run clean migrate seed test

APP_NAME=vuelang
DIST=dist/$(APP_NAME)

# ── Install ───────────────────────────────────────────────────────────────────
install:
	@echo "→ Installing Air (Go hot-reload)…"
	go install github.com/air-verse/air@latest
	@echo "→ Installing frontend dependencies…"
	cd ui && npm install
	@echo "→ Creating dummy ui/dist placeholder for Go embed…"
	mkdir -p ui/dist && echo "placeholder" > ui/dist/index.html
	@echo "→ Installing CLI…"
	go install ./cmd/vuelang
	@echo "✓ Done. Run 'make dev' to start."

# ── Development ───────────────────────────────────────────────────────────────
dev:
	@echo ""
	@echo "  ┌──────────────────────────────────────────────────┐"
	@echo "  │              Vuelang V2  DEV                     │"
	@echo "  │   App  →  http://localhost:9090                  │"
	@echo "  │   .go  →  Air rebuilds  (<1s)                    │"
	@echo "  │   .vue →  Vite HMR  (instant)                   │"
	@echo "  └──────────────────────────────────────────────────┘"
	@$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	ENV=development air -c .air.toml

dev-frontend:
	cd ui && npm run dev

# ── Production build ──────────────────────────────────────────────────────────
build:
	@echo "→ Building Vue 3 frontend…"
	cd ui && npm run build
	@echo "→ Compiling Go binary with embedded frontend…"
	mkdir -p dist
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$$(git describe --tags --always 2>/dev/null || echo dev)" \
		-o $(DIST) .
	@echo "✓ Binary ready → $(DIST)"

run:
	ENV=production ./$(DIST)

# ── Database ──────────────────────────────────────────────────────────────────
migrate:
	go run . --migrate-only

seed:
	go run . --seed

# ── Testing ───────────────────────────────────────────────────────────────────
test:
	go test ./... -v -race -cover

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ── CLI shortcuts ─────────────────────────────────────────────────────────────
make-model:
	vuelang make:model $(NAME)

make-controller:
	vuelang make:controller $(NAME)

make-middleware:
	vuelang make:middleware $(NAME)

make-migration:
	vuelang make:migration $(NAME)

# ── Cleanup ───────────────────────────────────────────────────────────────────
clean:
	rm -rf dist/ tmp/ build-errors.log coverage.out
	rm -rf ui/dist/assets ui/dist/*.js ui/dist/*.css ui/dist/*.html

# ── Docker ────────────────────────────────────────────────────────────────────
docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
