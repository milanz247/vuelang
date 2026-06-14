.PHONY: install dev dev-backend dev-frontend build run clean

# ── Install ───────────────────────────────────────────────────────────────────
install:
	@echo "→ Installing Air (Go hot-reload)..."
	go install github.com/air-verse/air@latest
	@echo "→ Installing frontend dependencies..."
	cd ui && npm install
	@echo "✓ Done. Run 'make dev' to start."

# ── Development ───────────────────────────────────────────────────────────────
dev:
	@echo ""
	@echo "  ┌──────────────────────────────────────────┐"
	@echo "  │            Vuelang  DEV                  │"
	@echo "  │   Open  →  http://localhost:8080         │"
	@echo "  │   .go   →  Air rebuilds  (<1s)           │"
	@echo "  │   .vue  →  Vite HMR (instant)           │"
	@echo "  └──────────────────────────────────────────┘"
	@$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	ENV=development air -c .air.toml

dev-frontend:
	cd ui && npm run dev

# ── Production ────────────────────────────────────────────────────────────────
build:
	@echo "→ Building Vue 3 frontend..."
	cd ui && npm run build
	@echo "→ Compiling Go binary with embedded frontend..."
	mkdir -p dist
	ENV=production go build -ldflags="-s -w" -o dist/vuelang .
	@echo "✓ Binary ready → dist/vuelang"

run:
	ENV=production ./dist/vuelang

# ── Cleanup ───────────────────────────────────────────────────────────────────
clean:
	rm -rf dist/ tmp/ build-errors.log
	rm -rf ui/dist/assets ui/dist/*.js ui/dist/*.css
