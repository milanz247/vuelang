package main

import "embed"

// embeddedUI holds the compiled Vue 3 frontend (ui/dist).
// It is populated by "make build" which runs "cd ui && npm run build" first.
//
// In development (make dev) this variable is never read — the server
// proxies to Vite instead — so ui/dist does not need to exist.
//
//go:embed all:ui/dist
var embeddedUI embed.FS
