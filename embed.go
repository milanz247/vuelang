package main

import "embed"

// embeddedUI holds the compiled Vue 3 frontend (ui/dist).
// Populated by "make build" → "cd ui && npm run build".
// In dev mode the server proxies to Vite and never reads this.
//
//go:embed all:ui/dist
var embeddedUI embed.FS
